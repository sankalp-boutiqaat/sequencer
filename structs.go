package sequencer

//
// Options to be passed on to Sequencer.
//
type Options struct {
	Name    string
	Start   int
	Limit   int
	Rolling bool
	Reverse bool
}
