package worker

import "sync"

// Pool is a simple goroutine worker pool. Submit tasks and Stop when done.
type Pool struct {
	tasks chan func()
	wg    sync.WaitGroup
}

// NewPool creates a pool with the given number of workers.
func NewPool(workers int) *Pool {
	if workers <= 0 {
		workers = 4
	}
	p := &Pool{tasks: make(chan func(), workers*4)}
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for t := range p.tasks {
				safeRun(t)
			}
		}()
	}
	return p
}

func safeRun(f func()) {
	defer func() {
		_ = recover()
	}()
	f()
}

// Submit enqueues a task. The pool must be stopped later to wait for completion.
func (p *Pool) Submit(f func()) {
	if f == nil {
		return
	}
	p.tasks <- f
}

// Stop closes the tasks channel and waits for workers to exit.
func (p *Pool) Stop() {
	close(p.tasks)
	p.wg.Wait()
}
