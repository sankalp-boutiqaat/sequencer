package sequencer

//
// Options to be passed on to Sequencer.
//
type Options struct {
	Key     Key
	Start   int
	Limit   int
	Rolling bool
	Reverse bool
}

type Key struct {
	Name   string
	Bucket string
}
