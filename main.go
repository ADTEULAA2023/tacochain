package main

import (
	"errors"
	"log"
	"os"

	"github.com/ADTEULAA2023/tacochain/cmd"
)

func main() {
	path := "./tmp"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	cmd.Execute()
}
