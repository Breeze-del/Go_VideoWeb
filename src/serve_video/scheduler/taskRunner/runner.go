package taskRunner

type Runner struct {
	// 决定是接受还是处理数据
	Controller controlChan
	// 出现错误
	Error controlChan
	// 存放数据
	Data dataChan
	// 数据缓存大小
	dataSize int
	// 是否在任务结束后关闭
	longLived bool
	// 分发数据函数
	Dispatcher fn
	// 处理数据函数
	Executor fn
}

func NewRunner(size int, longlived bool, d fn, e fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1),
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size),
		dataSize:   size,
		longLived:  longlived,
		Dispatcher: d,
		Executor:   e,
	}
}

func (r *Runner) startDispatch() {
	defer func() {
		if !r.longLived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()

	for {
		select {
		case con := <-r.Controller:
			if con == READY_TO_DISPATCH {
				err := r.Dispatcher(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_EXECUTE
				}
			}
			if con == READY_TO_EXECUTE {
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}
		case err := <-r.Error:
			if err == CLOSE {
				return
			}
		default:

		}
	}
}

func (r *Runner) startAll() {
	// 不是长久的就重新创建通道
	if !r.longLived {
		r.Controller = make(chan string, 1)
		r.Error = make(chan string, 1)
		r.Data = make(chan interface{}, r.dataSize)
	}
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}
