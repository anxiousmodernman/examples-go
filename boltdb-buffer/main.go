package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

const bktName = "messages"

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGHUP, syscall.SIGQUIT)

	ch := make(chan Msg)

	src := NewSource(50, ch)
	src.Start()

	maxBatch := 1000

	bt, err := NewBufferedThing("example.db", ch, maxBatch)
	if err != nil {
		fmt.Println("something is very wrong", err.Error())
		os.Exit(1)
	}
	bt.Start()

	// keep our program alive until one one of the signals is caught. Note that
	// Ctrl+C will raise a SIGINT.
	for {
		select {
		case sig := <-quit:
			fmt.Println("got signal:", sig)
			os.Exit(0)
		}
	}
}

// Msg is our message type.
type Msg struct {
	Value string
}

// Source is our type that emits messages. This simulates getting messages from
// a kafka queue.
type Source struct {
	// Sleep is a time in milliseconds to sleep before sending a message.
	Sleep int
	// After our sleep, we send on this channel.
	SendChan chan Msg
}

// NewSource constructs a Source from a channel and a sleep interval.
func NewSource(sleep int, ch chan Msg) *Source {
	var src Source
	src.Sleep = sleep
	src.SendChan = ch
	return &src
}

// Start starts our source. We spin off a goroutine. Note that we don't have
// a way to stop the Source without closing the program.
func (s *Source) Start() {
	// spin off a goroutine to send on our channel so Start return immediately.
	go func() {
		interval := time.Duration(s.Sleep) * time.Millisecond
		for {
			time.Sleep(interval)
			msg := Msg{Value: "blah"}
			s.SendChan <- msg
		}
	}()
}

// BufferedThing receives messages from the Source, buffers them to boltdb, and
// periodically flushes them.
type BufferedThing struct {
	// DB field is a pointer to an opened/instantiated boltdb file.
	DB       *bolt.DB
	RecvChan chan Msg
	MaxBatch int
}

func NewBufferedThing(path string, ch chan Msg, maxBatch int) (*BufferedThing, error) {

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	// creat our bucket if it doesn't exist
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	tx.CreateBucketIfNotExists([]byte(bktName))
	tx.Commit()

	bt := BufferedThing{
		DB:       db,
		RecvChan: ch,
		MaxBatch: maxBatch,
	}

	return &bt, nil
}

func (bt *BufferedThing) Start() {
	go bt.receiveLoop()
	go bt.batchProcess()
}

func (bt *BufferedThing) receiveLoop() {
	for {
		select {
		case msg := <-bt.RecvChan:

			err := bt.DB.Update(func(tx *bolt.Tx) error {
				bkt := tx.Bucket([]byte(bktName))
				id := uuid.Must(uuid.NewV4())
				return bkt.Put(id.Bytes(), []byte(msg.Value))
			})
			if err != nil {
				fmt.Println("updated failed; you should probably retry 3x then crash")
			}
		}
	}
}

// Key is an alias to []byte. A []Key is prettier than [][]byte when reading out
// of boltdb.
type Key []byte

// Value is an alias to []byte. A []Value is prettier than [][]byte when reading
// out of boltdb.
type Value []byte

// batchProcess reads, processes, and deletes key-value pairs from our messages
// bucket. This happens on an interval inside an infinite for loop.
func (bt *BufferedThing) batchProcess() {
	for {

		// sleep for an interval
		time.Sleep(time.Duration(5) * time.Second)

		// declare some slices outside of our DB query, so that we can use
		// these values in subsequent queries.
		var keys []Key
		var values []Value

		// Read up to MaxBatch keys and values from boltdb
		err := bt.DB.View(func(tx *bolt.Tx) error {
			bkt := tx.Bucket([]byte(bktName))
			cur := bkt.Cursor()

			count := 0

			// Scan from the first-sorted key in our bucket
			for k, v := cur.First(); k != nil; k, v = cur.Next() {
				if k == nil {
					fmt.Println("initial key was nil")
					return nil
				}

				// add to our slices
				keys = append(keys, Key(k))
				values = append(values, Value(v))

				count++
				if count >= bt.MaxBatch {
					fmt.Println("max reached")
					// we've already got N keys/values in our slices, so we
					// simply return from our query here.
					return nil
				}
			}
			// We return here if we had less than MaxBatch keys in our bucket.
			// E.g., we scanned all of them.
			return nil
		})
		if err != nil {
			fmt.Println("error:", err.Error())
			panic("couldn't read from bolt? we die now")
		}

		// Make sure the world makes sense. We probably don't need to do this.
		if len(keys) != len(values) {
			panic("expected 1:1 key to value")
		}

		// Make sure we have at least 1 k-v pair. We definitely NEED to do this.
		if len(keys) < 1 {
			fmt.Println("no keys/values; skipping processing")
			continue
		}

		// Now we do our fake processing. Here we'd send to some API or database.
		process(values)

		// We processed our stuff! Remove the just-processed keys/values from
		// boltdb. We can rely on the UUIDs in the keys slice to delete only
		// the successfully processed keys.
		err = bt.DB.Update(func(tx *bolt.Tx) error {
			bkt := tx.Bucket([]byte(bktName))
			for _, k := range keys {
				if err := bkt.Delete([]byte(k)); err != nil {
					return fmt.Errorf("delete failed: %v", err)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("error:", err)
			panic("could not delete processed keys")
		}

		// Monitor how many messages are left in our queue. If we're keeping
		// pace with incoming messages, we will keep this number hovering around
		// zero.
		err = bt.DB.View(func(tx *bolt.Tx) error {
			bkt := tx.Bucket([]byte(bktName))
			stats := bkt.Stats()
			fmt.Printf("keys in boltdb processing queue: %v\n", stats.KeyN)
			return nil
		})
	}
}

func process(values []Value) {
	fmt.Printf("successfully processed %v values\n", len(values))
}
