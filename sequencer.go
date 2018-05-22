package sequencer

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

const ADAPTER_TYPE_REDIS = "redis"

type Options struct {
	Name    string
	Start   int
	Limit   int
	Rolling bool
	Reverse bool
}

var _ Sequencer = (*RedisSequencer)(nil)

type Sequencer interface {
}

type options struct {
	name      string
	start     int
	limit     int
	isRolling bool
	isReverse bool
}

func InitializeRedisSequencer(opt Options, conf RedisConfig) *RedisSequencer {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisSimple := RedisSequencer{
		client: client,
		options: options{
			name:      opt.Name,
			start:     opt.Start,
			limit:     opt.Limit,
			isRolling: opt.Rolling,
			isReverse: opt.Reverse,
		},
	}
	return &redisSimple
}

type RedisConfig struct {
	Addr     string
	Password string
	PoolSize int
}

var ErrLimitReached error = errors.New("Limit Reached")

//handle case when key does not exists and start is specified.
type RedisSequencer struct {
	options
	client *redis.Client
}

func (this *RedisSequencer) Init() error {
	_, err := this.initKey()
	if err != nil {
		return err
	}
	return nil
}

func (this *RedisSequencer) Next() (int, error) {
	if !this.isRolling && !this.isReverse {
		return this.incr()
	}

	if !this.isRolling && this.isReverse {
		return this.decr()
	}

	if this.isRolling && !this.isReverse {
		val, err := this.incr()
		if err == nil {
			return val, nil
		}
		if err != nil && err != ErrLimitReached {
			return 0, err
		}
		//limit reached reset.
		err = this.reset()
		if err == nil {
			return this.incr()
		}
		return 0, err
	}

	if this.isRolling && this.isReverse {
		val, err := this.decr()
		if err == nil {
			return val, nil
		}
		if err != nil && err != ErrLimitReached {
			return 0, err
		}
		//limit reached reset.
		err = this.reset()
		if err == nil {
			return this.decr()
		}
		return 0, err
	}

	return 0, fmt.Errorf("Invalid Option. How did u landed here?")

}

func (this *RedisSequencer) initKey() (bool, error) {
	//create key only if it not exists.
	ret := this.client.SetNX(this.name, this.start, 0)
	if ret.Err() != nil {
		return false, ret.Err()
	}
	return ret.Val(), nil
}

func (this *RedisSequencer) decr() (int, error) {
	ret := this.client.Decr(this.name)
	if ret.Err() != nil {
		return 0, ret.Err()
	}
	res := int(ret.Val())
	if res <= this.limit {
		return 0, ErrLimitReached
	}
	return int(res), nil
}

func (this *RedisSequencer) incr() (int, error) {
	ret := this.client.Incr(this.name)
	if ret.Err() != nil {
		return 0, ret.Err()
	}
	res := int(ret.Val())
	if res >= this.limit {
		return 0, ErrLimitReached
	}
	return res, nil
}

//@todo: need to work on atomicity here.
func (this *RedisSequencer) reset() error {
	txpipe := this.client.TxPipeline()
	ret := txpipe.Get(this.name)
	valI64, _ := ret.Int64()
	val := int(valI64)
	if !this.isReverse && val >= this.limit {
		val = this.start
		txpipe.Set(this.name, this.start, 0)
	}
	if this.isReverse && val <= this.limit {
		val = this.start
		txpipe.Set(this.name, this.start, 0)
	}
	cmdr, err := txpipe.Exec()
	fmt.Printf("commander: %v\n", cmdr)
	if err != nil {
		return err
	}
	return nil
}

// Exposed methods lies here

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
