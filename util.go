package sequencer

import "fmt"

func Initialize(stype string, opt Options, conf interface{}) (Sequencer, error) {
	stype = ADAPTER_TYPE_REDIS
	var seq Sequencer
	switch stype {
	case ADAPTER_TYPE_REDIS:
		config, ok := conf.(RedisConfig)
		if !ok {
			return seq, fmt.Errorf("Expected RedisConfig, Got %T", conf)
		}
		redisS := InitializeRedisSequencer(opt, config)
		seq = redisS

	}
	return seq, nil
}
