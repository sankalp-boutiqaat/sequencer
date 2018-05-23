package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sanksons/sequencer"
)

func main() {

	var options sequencer.Options = sequencer.Options{

		Key: sequencer.Key{
			Name:   "seq2",
			Bucket: "1",
		},

		Start:   10,
		Limit:   0,
		Rolling: true,
		Reverse: true,
	}

	var conf sequencer.RedisConfig = sequencer.RedisConfig{
		Addrs: []string{"172.17.0.2:30001", "172.17.0.2:30003"},
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
