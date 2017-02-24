package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	i := 10
	for i > 0 {
		t := time.Now()
		daysAgo := random(10000)
		fmt.Println("daysago", daysAgo)
		t2 := t.AddDate(0, 0, -daysAgo).Add(time.Duration(int64(-daysAgo)) * time.Second)
		fmt.Println(t2)
		i--
	}
}

func random(max int64) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(int(max)) // careful casting
}
