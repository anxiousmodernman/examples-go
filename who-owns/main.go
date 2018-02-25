package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"syscall"
)

func main() {
	// in shell we do this
	// ls -l | awk '{print $3, $4 }'
	fi, err := os.Stat("foo.txt")
	if err != nil {
		log.Fatal(err)
	}

	switch fi.Sys().(type) {
	case *syscall.Stat_t:
		statt, ok := fi.Sys().(*syscall.Stat_t)
		if !ok {
			log.Fatal("wrong type")
		}
		// these are uint32
		fmt.Println(statt.Gid)
		fmt.Println(statt.Uid)
		uid, gid := fmt.Sprintf("%v", statt.Uid), fmt.Sprintf("%v", statt.Gid)
		fmt.Println(uid, gid)

		u, err := user.LookupId(uid)
		if err != nil {
			log.Fatal(err)
		}
		g, err := user.LookupGroupId(gid)
		if err != nil {
			log.Fatal(err)
		}

		mode := fi.Mode().Perm()

		if (0600 ^ mode) == 0 {

			fmt.Println("how does binary work")
		}

		if samePerms(0600, uint32(mode)) {

			fmt.Println("how does binary work pt. 2")
		}

		fmt.Printf("what is my mode: %o\n", uint32(mode&os.ModeSetgid))
		fmt.Printf("what is my mode: %o\n", uint32(mode))
		fmt.Printf("what is my M0d3: %o\n", uint32(0600))

		fmt.Printf("%s %s\n", u.Username, g.Name)
	}
}

func samePerms(yours int, mode uint32) bool {
	if (uint32(yours) ^ mode) == 0 {
		return true
	}
	return false
}
