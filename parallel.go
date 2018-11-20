package Li

import (
	"runtime"
)

type WorkFunc func(interface{}) interface{}

type Parallel struct {
	Func func(interface{}) interface{}
	In   chan interface{}
	Out  chan interface{}
	N    int
}

var FINISH = new(interface{})

func NewParallel(maxWork int, fn WorkFunc) *Parallel {
	par := &Parallel{
		Func: fn,
		In:   make(chan interface{}, maxWork),
		Out:  make(chan interface{}, maxWork),
		N:    runtime.GOMAXPROCS(0),
	}
	for i := 0; i < par.N; i++ {
		go func() {
			for {
				input := <-par.In
				if input == FINISH {
					break
				}
				output := fn(input)
				par.Out <- output

			}
		}()
	}

	return par
}

func (par *Parallel) Add1(input interface{}) {
	par.In <- input
}

func (par *Parallel) Wait1() {
	<-par.Out
}

func (par *Parallel) Finish() {
	for i := 0; i < par.N; i++ {
		par.In <- FINISH
	}
}
