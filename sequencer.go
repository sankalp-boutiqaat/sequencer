package sequencer

const ADAPTER_TYPE_REDIS = "redis"

var _ Sequencer = (*RedisSequencer)(nil)

type Sequencer interface {
	Next() (int, error)
	Reset() error
	Destroy() error
}

type options struct {
	name      string
	start     int
	limit     int
	isRolling bool
	isReverse bool
}
