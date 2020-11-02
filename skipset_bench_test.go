package skipset

import (
	"math"
	"sync"
	"testing"
)

const initsize = 1 << 10 // for `contains` `1Delete9Insert90Contains` `1Range9Delete90Insert900Contains`
const randN = math.MaxUint32

func BenchmarkInsert(b *testing.B) {
	b.Run("skipset", func(b *testing.B) {
		l := NewInt64()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Insert(int64(fastrandn(randN)))
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Store(int64(fastrandn(randN)), nil)
			}
		})
	})
}

func BenchmarkContains100Hits(b *testing.B) {
	b.Run("skipset", func(b *testing.B) {
		l := NewInt64()
		for i := 0; i < initsize; i++ {
			l.Insert(int64(i))
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = l.Contains(int64(fastrandn(initsize)))
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize; i++ {
			l.Store(int64(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(int64(fastrandn(initsize)))
			}
		})
	})
}

func BenchmarkContains50Hits(b *testing.B) {
	const rate = 2
	b.Run("skipset", func(b *testing.B) {
		l := NewInt64()
		for i := 0; i < initsize*rate; i++ {
			if fastrandn(rate) == 0 {
				l.Insert(int64(i))
			}
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = l.Contains(int64(fastrandn(initsize * rate)))
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize*rate; i++ {
			if fastrandn(rate) == 0 {
				l.Store(int64(i), nil)
			}
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(int64(fastrandn(initsize * rate)))
			}
		})
	})
}

func BenchmarkContainsNoHits(b *testing.B) {
	b.Run("skipset", func(b *testing.B) {
		l := NewInt64()
		invalid := make([]int64, 0, initsize)
		for i := 0; i < initsize*2; i++ {
			if i%2 == 0 {
				l.Insert(int64(i))
			} else {
				invalid = append(invalid, int64(i))
			}
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = l.Contains(invalid[fastrandn(uint32(len(invalid)))])
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		invalid := make([]int64, 0, initsize)
		for i := 0; i < initsize*2; i++ {
			if i%2 == 0 {
				l.Store(int64(i), nil)
			} else {
				invalid = append(invalid, int64(i))
			}
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(invalid[fastrandn(uint32(len(invalid)))])
			}
		})
	})
}

func Benchmark50Insert50Contains(b *testing.B) {
	b.Run("skipset", func(b *testing.B) {
		l := NewInt64()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 5 {
					l.Insert(int64(fastrandn(randN)))
				} else {
					l.Contains(int64(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 5 {
					l.Store(int64(fastrandn(randN)), nil)
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
}

func Benchmark30Insert70Contains(b *testing.B) {
	b.Run("skipset", func(b *testing.B) {
		l := NewInt64()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 3 {
					l.Insert(int64(fastrandn(randN)))
				} else {
					l.Contains(int64(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 3 {
					l.Store(int64(fastrandn(randN)), nil)
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
}

func Benchmark1Delete9Insert90Contains(b *testing.B) {
	b.Run("skipset", func(b *testing.B) {
		l := NewInt64()
		for i := 0; i < initsize; i++ {
			l.Insert(int64(i))
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(100)
				if u == 1 {
					l.Insert(int64(fastrandn(randN)))
				} else if u == 2 {
					l.Delete(int64(fastrandn(randN)))
				} else {
					l.Contains(int64(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize; i++ {
			l.Store(int64(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(100)
				if u == 1 {
					l.Store(int64(fastrandn(randN)), nil)
				} else if u == 2 {
					l.Delete(int64(fastrandn(randN)))
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
}

func Benchmark1Range9Delete90Insert900Contains(b *testing.B) {
	b.Run("skipset", func(b *testing.B) {
		l := NewInt64()
		for i := 0; i < initsize; i++ {
			l.Insert(int64(i))
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(1000)
				if u == 0 {
					l.Range(func(i int, score int64) bool {
						return true
					})
				} else if u > 10 && u < 20 {
					l.Delete(int64(fastrandn(randN)))
				} else if u >= 100 && u < 190 {
					l.Insert(int64(fastrandn(randN)))
				} else {
					l.Contains(int64(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize; i++ {
			l.Store(int64(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(1000)
				if u == 0 {
					l.Range(func(key, value interface{}) bool {
						return true
					})
				} else if u > 10 && u < 20 {
					l.Delete(fastrandn(randN))
				} else if u >= 100 && u < 190 {
					l.Store(fastrandn(randN), nil)
				} else {
					l.Load(fastrandn(randN))
				}
			}
		})
	})
}
