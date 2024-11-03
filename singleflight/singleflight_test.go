package singleflight

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestDo 测试 Do 方法的各种场景。
func TestDo(t *testing.T) {
	var g Group
	v, err, _ := g.Do("key", func() (any, error) {
		return "val", nil
	})
	if got, want := fmt.Sprintf("%v (%T)", v, v), "val (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}
func TestDupDo(t *testing.T) {
	var g = Group{}
	var calls int32
	ch := make(chan string)
	fn := func() (any, error) {
		atomic.AddInt32(&calls, 1)
		return <-ch, nil
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			v, err, _ := g.Do("key", fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			if v.(string) != "val" {
				t.Errorf("got %q; want %q", v, "val")
			}
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	ch <- "val"
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("number of calls = %d; want 1", got)
	}

}
