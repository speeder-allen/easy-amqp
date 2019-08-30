package pool_test

import (
	"github.com/speeder-allen/easy-amqp/pool"
	"gotest.tools/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	var number int64 = 0
	p := pool.NewPool(3)
	assert.Equal(t, p.Cap(), int64(3))
	var wg sync.WaitGroup
	wg.Add(6)
	for i := 0; i < 6; i++ {
		go func(idx int) {
			defer wg.Done()
			p.Get()
			atomic.AddInt64(&number, 1)
		}(i)
	}
	time.Sleep(time.Microsecond * 700)
	assert.Equal(t, number, int64(3))
	p.Put()
	time.Sleep(time.Microsecond * 700)
	assert.Equal(t, number, int64(4))
	p.Put()
	p.Put()
	time.Sleep(time.Microsecond * 700)
	assert.Equal(t, number, int64(6))
	wg.Wait()
}
