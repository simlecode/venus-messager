package utils

import (
	"sync"
)

type Par struct {
	par_ctlor chan interface{}
	index     int64
	group     sync.WaitGroup
	mutx      sync.Mutex
	errs      *MultiError
}

func NewPar(maxParCount int) *Par {
	return &Par{
		par_ctlor: make(chan interface{}, maxParCount),
		group:     sync.WaitGroup{},
		errs:      new(MultiError),
	}
}

func (this *Par) Go(f func() error) {
	this.group.Add(1)

	this.par_ctlor <- this.index
	this.index++

	go func() {
		defer func() {
			<-this.par_ctlor
			this.group.Done()
		}()

		if err := f(); err != nil {
			this.mutx.Lock()
			defer this.mutx.Unlock()
			this.errs.AddError(err)
		}
	}()
}

func (this *Par) GoV2(f func(args ...interface{}) error, args ...interface{}) {
	this.group.Add(1)

	this.par_ctlor <- this.index
	this.index++

	go func() {
		defer func() {
			<-this.par_ctlor
			this.group.Done()
		}()
		if err := f(args...); err != nil {
			this.mutx.Lock()
			defer this.mutx.Unlock()

			this.errs.AddError(err)
		}
	}()
}

func (this *Par) Wait() *MultiError {
	this.group.Wait()
	return this.errs
}
