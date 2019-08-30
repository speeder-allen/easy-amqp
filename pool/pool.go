package pool

var defaultCapacity int64 = 10

type pool struct {
	capacity int64
	cc       chan uint8
}

func (p *pool) Get() {
	p.cc <- 1
}

func (p *pool) Put() {
	<-p.cc
}

func (p *pool) Cap() int64 {
	return p.capacity
}

func (p *pool) Close() {
	close(p.cc)
}

func NewPool(capacity int64) *pool {
	if capacity < 0 {
		capacity = defaultCapacity
	}
	return &pool{
		capacity: capacity,
		cc:       make(chan uint8, capacity),
	}
}
