package base

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

const msize = 10000

func BenchmarkStrMap(b *testing.B) {
	m := map[string]string{}
	check := []string{}
	b.StopTimer()
	for i := 0; i < msize; i++ {
		for {
			v := fmt.Sprintf("%d", rand.Int63())
			if _, ok := m[v]; ok {
				continue
			}
			m[v] = v
			if i%10 == 0 {
				check = append(check, v)
			}
			break
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range check {
			_ = m[c]
		}
	}
}

func BenchmarkIntMap(b *testing.B) {
	m := map[int64]string{}
	check := []int64{}
	b.StopTimer()
	for i := 0; i < msize; i++ {
		for {
			v := rand.Int63()
			if _, ok := m[v]; ok {
				continue
			}
			m[v] = fmt.Sprintf("%d", v)
			if i%10 == 0 {
				check = append(check, v)
			}
			break
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range check {
			_ = m[c]
		}
	}
}

func BenchmarkSyncIntMap(b *testing.B) {
	m := NewSyncMap[int64, string]()
	check := []int64{}
	b.StopTimer()
	for i := 0; i < msize; i++ {
		for {
			v := rand.Int63()

			if _, ok := m.Load(v); ok {
				continue
			}
			m.Store(v, fmt.Sprintf("%d", v))
			if i%10 == 0 {
				check = append(check, v)
			}
			break
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range check {
			_, _ = m.Load(c)
		}
	}
}

func BenchmarkSyncStrMap(b *testing.B) {
	m := NewSyncMap[string, string]()
	check := []string{}
	b.StopTimer()
	for i := 0; i < msize; i++ {
		for {
			v := fmt.Sprintf("%d", rand.Int63())
			if _, ok := m.Load(v); ok {
				continue
			}
			m.Store(v, v)
			if i%10 == 0 {
				check = append(check, v)
			}
			break
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range check {
			_, _ = m.Load(c)
		}
	}
}

func BenchmarkGoSyncIntMap(b *testing.B) {
	m := &sync.Map{}
	check := []int64{}
	b.StopTimer()
	for i := 0; i < msize; i++ {
		for {
			v := rand.Int63()
			if _, ok := m.Load(v); ok {
				continue
			}
			m.Store(v, fmt.Sprintf("%d", v))
			if i%10 == 0 {
				check = append(check, v)
			}
			break
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range check {
			_, _ = m.Load(c)
		}
	}
}

func BenchmarkGoSyncStrMap(b *testing.B) {
	m := &sync.Map{}
	check := []string{}
	b.StopTimer()
	for i := 0; i < msize; i++ {
		for {
			v := fmt.Sprintf("%d", rand.Int63())
			if _, ok := m.Load(v); ok {
				continue
			}
			m.Store(v, v)
			if i%10 == 0 {
				check = append(check, v)
			}
			break
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range check {
			_, _ = m.Load(c)
		}
	}
}
