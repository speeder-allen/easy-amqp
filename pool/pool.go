package pool

var defaultCapacity int64 = 10

type Pool struct {
	capacity int64
	cc       chan uint8
}

func (p *Pool) Get() {
	p.cc <- 1
}

func (p *Pool) Put() {
	<-p.cc
}

func (p *Pool) Cap() int64 {
	return p.capacity
}

func (p *Pool) Close() {
	close(p.cc)
}

func NewPool(capacity int64) *Pool {
	if capacity < 0 {
		capacity = defaultCapacity
	}
	return &Pool{
		capacity: capacity,
		cc:       make(chan uint8, capacity),
	}
}
