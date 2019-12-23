package bobo

import (
	"io"
	"sync"
	"time"

	"github.com/jpillora/overseer/fetcher"
)

type MultiFetcher struct {
	List           []fetcher.Interface
	GlobalInterval time.Duration

	result          chan fetchResult
	runningStatusMu sync.RWMutex
	runningStatus   map[int]bool
}

func (f *MultiFetcher) Init() error {
	const minInterval = 1 * time.Second
	if f.GlobalInterval < minInterval {
		f.GlobalInterval = minInterval
	}
	f.runningStatus = make(map[int]bool, len(f.List))
	f.result = make(chan fetchResult)

	for i, v := range f.List {
		f.updateRunningStatus(i, false)
		if err := v.Init(); err != nil {
			return err
		}
	}
	return nil
}

func (f *MultiFetcher) Fetch() (io.Reader, error) {
	time.Sleep(f.GlobalInterval)

	for i, v := range f.List {
		if f.isRunning(i) {
			continue
		}

		go func(i int, v fetcher.Interface) {
			f.updateRunningStatus(i, true)
			f.result <- doFetch(v)
			f.updateRunningStatus(i, false)
		}(i, v)
	}

	data := f.fetchResult()
	return data.Reader, data.Err
}

func (f *MultiFetcher) fetchResult() fetchResult {
	return <-f.result
}

func (f *MultiFetcher) isRunning(num int) bool {
	f.runningStatusMu.RLock()
	defer f.runningStatusMu.RUnlock()
	return f.runningStatus[num]
}

func (f *MultiFetcher) updateRunningStatus(num int, status bool) {
	f.runningStatusMu.Lock()
	f.runningStatus[num] = status
	f.runningStatusMu.Unlock()
}

func doFetch(f fetcher.Interface) fetchResult {
	r, err := f.Fetch()
	return fetchResult{
		Reader: r,
		Err:    err,
	}
}

type fetchResult struct {
	Reader io.Reader
	Err    error
}
