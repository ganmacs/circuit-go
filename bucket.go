package circuit

import (
	"sync"
)

type Bucket struct {
	sucess      uint64
	fail        uint64
	totalSucess uint64
	totalFail   uint64
	request     uint64
	mutex       sync.Mutex
}

func NewBucket() *Bucket {
	return &Bucket{}
}

func (b *Bucket) Rate() float64 {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return float64(b.fail) / float64(b.request)
}

func (b *Bucket) Sucess() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.sucess++
	b.totalSucess++
	b.request++
}

func (b *Bucket) Fail() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.fail++
	b.totalFail++
	b.request++
}

func (b *Bucket) Reset() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.fail = 0
	b.sucess = 0
	b.request = 0
}
