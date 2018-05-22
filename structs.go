package sequencer

type Options struct {
	Name    string
	Start   int
	Limit   int
	Rolling bool
	Reverse bool
}
