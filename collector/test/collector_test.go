package test

import (
	"sync"
	"testing"

	"github.com/phantom-atom/crawler/collector"
	"github.com/phantom-atom/crawler/collector/typed"
)

func TestChannelCollectorElemToPtr(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()

	ch := make(chan *int, 4)
	var c collector.Collector = typed.NewCollector(ch)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		for i := 0; i < 10; i++ {
			if err := c.Collect(i); err != nil {
				t.Error(err)
			}
		}
		close(ch)
	}()
	wg.Wait()

	for v := range ch {
		t.Log(*v)
	}
}

func TestChannelCollectorPtrToElem(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()

	ch := make(chan int, 4)
	var c collector.Collector = typed.NewCollector(ch)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		for i := 0; i < 10; i++ {
			if err := c.Collect(&i); err != nil {
				t.Error(err)
			}
		}
		close(ch)
	}()
	wg.Wait()

	for v := range ch {
		t.Log(*v)
	}
}
