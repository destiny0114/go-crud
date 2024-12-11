package rate

import "time"

type Capacity int

type Limiter struct {
	capacity        Capacity
	leak_rate       int
	bucket          []time.Time
	lastRequestTime time.Time
}

func NewRateLimiter(c Capacity, r int) *Limiter {
	return &Limiter{
		capacity:        c,
		leak_rate:       r,
		bucket:          make([]time.Time, 0, c),
		lastRequestTime: time.Now(),
	}
}

func (lim *Limiter) IsRequestAllow() bool {
	now := time.Now()
	elapsed_time := now.Sub(lim.lastRequestTime).Seconds()

	if leaked := int(elapsed_time * float64(lim.leak_rate)); leaked > 0 {
		for i := 0; i < min(leaked, len(lim.bucket)); i++ {
			_, lim.bucket = lim.bucket[0], lim.bucket[1:]
		}
		lim.lastRequestTime = now
	}

	if len(lim.bucket) < int(lim.capacity) {
		lim.bucket = append(lim.bucket, now)
		return true
	}
	return false
}
