package inkwell

import "math/rand"

// Ladder is what somebody says the n-th time the player does the same
// fruitless thing. A game that answers a repeated failure with one line
// sounds like a machine; a Ladder answers the first attempt with the
// first rung, the second with the second, and keeps a character's
// patience visibly running out.
//
// Past the last rung a Ladder picks one at random, which keeps a stock
// refusal from turning into a stuck record — unless it names a Last,
// which stands from then on. A character who has given up has given up.
type Ladder struct {
	// Rungs are the lines, in the order the player earns them.
	Rungs []Text
	// Last stands once the rungs run out. Empty picks a rung at random.
	Last Text
}

// At is the line for the n-th attempt, counting the first as 1. A Ladder
// with no rungs says nothing at all.
func (l Ladder) At(n int) Text {
	switch {
	case len(l.Rungs) == 0:
		return l.Last
	case n < 1:
		return l.Rungs[0]
	case n <= len(l.Rungs):
		return l.Rungs[n-1]
	case !l.Last.Empty():
		return l.Last
	default:
		return l.Rungs[rand.Intn(len(l.Rungs))]
	}
}
