package sequencer

//
// Options to be passed on to Sequencer.
//
type Options struct {
	Key struct {
		Name   string
		Bucket string
	}
	Start   int
	Limit   int
	Rolling bool
	Reverse bool
}
