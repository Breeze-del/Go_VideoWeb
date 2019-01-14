package taskRunner

import "time"

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

func (w *Worker) startWorker() {
	for {
		select {
		case <-w.ticker.C:
			go w.runner.startAll()
		}
	}
}

func Start() {
	// start video file cleaning
	r := NewRunner(4, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(10, r)
	go w.startWorker()
}
