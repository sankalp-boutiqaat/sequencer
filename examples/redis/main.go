package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sanksons/sequencer"
)

func main() {

	var options sequencer.Options = sequencer.Options{
		Name:    "sequencer6",
		Start:   10,
		Limit:   0,
		Rolling: true,
		Reverse: true,
	}

	var conf sequencer.RedisConfig = sequencer.RedisConfig{
		Addr: "localhost:6379",
	}

	sequenceG, err := sequencer.Initialize(sequencer.ADAPTER_TYPE_REDIS, options, conf)

	if err != nil {
		log.Fatal(err)
	}
	for {
		d, err := sequenceG.Next()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(d)
		time.Sleep(1 * time.Second)
	}
}
