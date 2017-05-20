package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

var (
	confPath = flag.String("conf", "conf.json", "path to config file")
)

func main() {

	flag.Parse()

	conf, err := NewConfig(*confPath)
	if err != nil {
		log.Fatal(err)
	}

	t := NewThing(conf)

	log.Printf("thing initialized: %v", t)
}

type Config struct {
	Foo string
	Baz int
}

// NewConfig reads a json config file.
func NewConfig(path string) (Config, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	// Declare a config. Initial values are the "zero values".
	var conf Config
	// Pass a pointer to conf to Unmarshal. Read data into it.
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return Config{}, err
	}
	// conf has data
	return conf, nil
}

// Thing is our thing that needs configuration.
type Thing struct {
	Msg string
}

// NewThing is our Thing constructor.
func NewThing(conf Config) *Thing {

	// Declare our Thing
	var t Thing
	// Set our field
	t.Msg = fmt.Sprintf("Foo is %s and Baz is %v", conf.Foo, conf.Baz)
	// Return a pointer to our Thing
	return &t
}
