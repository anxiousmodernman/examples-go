package main

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
)

func main() {
	// We just want to print a json representation, so we don't need to
	// give a specific type here.
	var anything map[string]interface{}
	toml.DecodeFile("inventory.toml", &anything)
	data, _ := json.MarshalIndent(&anything, "", "   ")
	fmt.Println(string(data))
}
