package taskrunner

type Runner struct {
	//3个channel 有各自的作用
	Controller controlChan
	Error      controlChan
	Data       dataChan
	dataSize   int
	longLived  bool //决定Runner是否长期存活
	Dispatcher fn   //fn是函数类型
	Executor   fn
}

//Runner 的构建函数
func NewRunner(size int, longlived bool, d fn, e fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1),
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size), //利用interface模拟泛型
		longLived:  longlived,
		dataSize:   size,
		Dispatcher: d,
		Executor:   e,
	}
}

func (r *Runner) starDispatch() {
	defer func() { //starDispatch 函数执行完退出后会执行这段语句，根据longLived来决定是否释放Runner的空间
		if !r.longLived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}() //这个括号是给该匿名函数传参的

	for {
		select {
		case c := <-r.Controller:
			if c == READY_TO_DISPATCH {
				err := r.Dispatcher(r.Data) //r.dataChan  是一个数据channel,  r.Controller 应该是传递控制标志的channel
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_EXECUTE
				}
			}

			if c == READY_TO_EXECUTE {
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}
		case e := <-r.Error:
			if e == CLOSE {
				return
			}
		default:
		}
	}
}

func (r *Runner) StartAll() {
	r.Controller <- READY_TO_DISPATCH //初始时应该先给channnel传递一个消息，以开启
	r.starDispatch()
}
