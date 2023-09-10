package threadpool

type Job func()

type Interface interface {
	Start()
	AddJob(job Job)
	Stop()
}

type ThreadPool struct {
	workers    []*Worker
	queue      chan Job
	readyQueue chan chan Job
}

func NewThreadPool(size int) *ThreadPool {
	return &ThreadPool{
		workers:    make([]*Worker, size),
		queue:      make(chan Job, size),
		readyQueue: make(chan chan Job, size),
	}
}

func (p *ThreadPool) Start() {
	for i := 0; i < len(p.workers); i++ {
		p.workers[i] = newWorker(i, p.readyQueue)
		p.workers[i].start()
	}

	go p.dispatch()
}

func (p *ThreadPool) dispatch() {
	for {
		select {
		case job := <-p.queue:
			worker := <-p.readyQueue
			worker <- job
		}
	}
}

func (p *ThreadPool) AddJob(job Job) {
	p.queue <- job
}

func (p *ThreadPool) Stop() {
	for _, worker := range p.workers {
		worker.stop()
	}

	for _, worker := range p.workers {
		<-worker.stopped
	}
}

type WorkerInterface interface {
	start()
	stop()
}

type Worker struct {
	id         int
	job        chan Job
	readyQueue chan chan Job
	quit       chan bool
	isActive   bool
	isRunning  bool
	stopped    chan bool
}

func newWorker(id int, readyQueue chan chan Job) *Worker {
	return &Worker{
		id:         id,
		job:        make(chan Job),
		quit:       make(chan bool),
		stopped:    make(chan bool),
		readyQueue: readyQueue,
	}
}

func (w *Worker) start() {
	w.isActive = true

	go func() {
		for {
			w.readyQueue <- w.job
			select {
			case job := <-w.job:
				w.isRunning = true
				job()
				w.isRunning = false

			case <-w.quit:
				w.isActive = false
				w.stopped <- true
				return
			}
		}
	}()
}

func (w *Worker) stop() {
	w.quit <- true
}
