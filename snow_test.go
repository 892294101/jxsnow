package jxsnow

import (
	"testing"
	"time"
)

func TestGenerator_Generate(t *testing.T) {
	ch := make(chan int64, 1000)
	m := make(map[int64]bool)
	go func() {
		for i := range ch {
			if m[i] {
				t.Error(i)
			} else {
				m[i] = true
			}
		}
	}()

	g1, _ := NewGenerator(1)
	g2, _ := NewGenerator(2)
	g3, _ := NewGenerator(3)
	g4, _ := NewGenerator(4)
	g5, _ := NewGenerator(5)
	g6, _ := NewGenerator(6)

	for i := 0; i < 1000; i++ {
		switch i % 6 {
		case 0:
			go func() {
				g2, _ := g2.Generate()
				ch <- g2
			}()
		case 1:
			go func() {
				g3, _ := g3.Generate()
				ch <- g3
			}()
		case 2:
			go func() {
				g4, _ := g4.Generate()
				ch <- g4
			}()
		case 3:
			go func() {
				g5, _ := g5.Generate()
				ch <- g5
			}()
		case 4:
			go func() {
				g6, _ := g6.Generate()
				ch <- g6
			}()
		case 5:
			go func() {
				g1, _ := g1.Generate()
				ch <- g1
			}()
		default:
			t.Fatal(i)
		}
	}

	time.Sleep(time.Second * 30)
}

func BenchmarkGenerator_Generate(b *testing.B) {
	g, _ := NewGenerator(1)
	for i := 0; i < b.N; i++ {
		g.Generate()
	}
}
