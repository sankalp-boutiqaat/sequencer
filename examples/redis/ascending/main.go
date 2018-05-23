package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sanksons/sequencer"
)

func main() {

	//
	// Sequencer Options.
	//
	var options sequencer.Options = sequencer.Options{
		//Define Key and Bucket.
		//Bucket is important it makes sure all data relevent to Sequencer
		// is located in single bucket.
		Key: sequencer.Key{
			Name:   "ascseq",
			Bucket: "1",
		},
		//Start counter from
		Start: 1,
		//Max limit upto which to increament counter. -2 for no limit
		Limit: -2,
	}

	//
	// Redis adapter configuration.
	//
	var conf sequencer.RedisConfig = sequencer.RedisConfig{
		Addrs: []string{"172.17.0.2:30001", "172.17.0.2:30003"},
	}

	//
	// Init the sequencer.
	//
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
		time.Sleep(100 * time.Millisecond)
	}
}
