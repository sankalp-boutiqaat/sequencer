package sequencer

//Redis Adapter.
var _ Sequencer = (*RedisSequencer)(nil)

//
// Sequencer interface
//
type Sequencer interface {
	//Takes care of initialization of sequncer.
	Init() error

	//Gives the next sequence.
	Next() (int, error)

	//Reset the sequencer to original state.
	Reset() error

	//Completely destroy all information related to sequencer.
	Destroy() error
}

//
// These are the common options are used by sequencer adpaters.
//
type options struct {
	name      string
	bucket    string
	start     int
	limit     int
	isRolling bool
	isReverse bool
}
