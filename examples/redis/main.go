package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sanksons/sequencer"
)

func main() {

	var options sequencer.Options = sequencer.Options{
		Name:    "sequencer1",
		Start:   1,
		Limit:   10,
		Rolling: true,
		Reverse: false,
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
