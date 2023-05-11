package main

import (
	"log"
	"os"

	"github.com/ADTEULAA2023/tacochain/cmd"
)

func main() {
	err := os.Mkdir("./tmp", os.ModePerm)
	if err != nil {
		log.Panicln(err)
	}

	cmd.Execute()
}
