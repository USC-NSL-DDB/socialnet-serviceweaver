package perf

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

type PerfThreadState interface{}
type PerfRequest interface{}
type PerfRequestWithTime struct {
	start_us uint64
	req      PerfRequest
}

type PerfAdapter interface {
	CreateGoroutineState() PerfThreadState
	GenRequest(PerfThreadState) PerfRequest
	ServeRequest(PerfThreadState, PerfRequest) bool
}

type Trace struct {
	absl_start_us uint64
	start_us      uint64
	duration      uint64
}

// Closed-loop, possion arrival.
type Perf struct {
	adapter    PerfAdapter
	traces     []Trace
	real_mops_ float64
}

func NewPerf(adapter PerfAdapter) *Perf {
	return &Perf{adapter: adapter}
}

func (p *Perf) Clear() {
	p.traces = p.traces[:0]
	p.real_mops_ = 0
}

func (p *Perf) GenRequests(
	all_reqs [][]PerfRequestWithTime, thread_states []PerfThreadState, num_threads uint32, target_mops float64, duration_us uint64,
) {
	wg := sync.WaitGroup{}

	for i := uint32(0); i < num_threads; i++ {
		wg.Add(1)
		go func(reqs []PerfRequestWithTime, thread_state PerfThreadState) {
			defer wg.Done()
			state := p.adapter.CreateGoroutineState()
			src := rand.NewSource(time.Now().UnixNano())
			rng := rand.New(src)

			curr_us := uint64(0)

			for curr_us < duration_us {
				r := rng.Float64()
				tmp := -math.Log(r) / (target_mops * float64(num_threads))
				interval := uint64(math.Max(1, float64(lround(tmp))))

				req_with_time := PerfRequestWithTime{
					start_us: curr_us, req: p.adapter.GenRequest(state),
				}

				reqs = append(reqs, req_with_time)
				curr_us += uint64(interval)
			}
		}(all_reqs[i], thread_states[i])
	}

	wg.Wait()
}

func (p *Perf) Benchmark(
	all_reqs [][]PerfRequestWithTime, thread_states []PerfThreadState, num_threads uint32, miss_ddl_thresh_us *uint64,
) []Trace {
	all_traces := make([][]Trace, num_threads)
	for i := uint32(0); i < num_threads; i++ {
		all_traces[i] = make([]Trace, len(all_reqs[i]))
	}
	wg := sync.WaitGroup{}

	for i := uint32(0); i < num_threads; i++ {
		wg.Add(1)
		go func(reqs []PerfRequestWithTime, thread_state PerfThreadState, traces []Trace) {
			defer wg.Done()
			start_us := microtime()
			for i, req_with_time := range reqs {
				relative_us := microtime() - start_us
				if req_with_time.start_us > relative_us {
					time.Sleep(time.Duration(req_with_time.start_us-relative_us) * time.Microsecond)
				} else if miss_ddl_thresh_us != nil && req_with_time.start_us+(*miss_ddl_thresh_us) < relative_us {
					continue
				}

				now := microtime()
				trace := Trace{
					absl_start_us: now, start_us: now - start_us, duration: 0,
				}
				ok := p.adapter.ServeRequest(thread_state, req_with_time.req)
				trace.duration = microtime() - start_us - trace.start_us
				if ok {
					traces[i] = trace
				}
			}
		}(all_reqs[i], thread_states[i], all_traces[i])
	}

	wg.Wait()

	gathered_traces := make([]Trace, 0)
	for i := uint32(0); i < num_threads; i++ {
		gathered_traces = append(gathered_traces, all_traces[i]...)
	}
	return gathered_traces
}

func (p *Perf) CreateThreadStates(thread_states []PerfThreadState, num_threads uint32) {
	for i := uint32(0); i < num_threads; i++ {
		thread_states = append(thread_states, p.adapter.CreateGoroutineState())
	}
}

func (p *Perf) Run(
	num_threads uint32, target_mops float64, duration_us uint64, warmup_us uint64, miss_ddl_thresh_us uint64,
) {
	p.RunMultiClients(num_threads, target_mops, duration_us, warmup_us, miss_ddl_thresh_us)
}

func (p *Perf) RunMultiClients(
	num_threads uint32, target_mops float64, duration_us uint64, warmup_us uint64, miss_ddl_thresh_us uint64,
) {
	thread_states := make([]PerfThreadState, 0)
	p.CreateThreadStates(thread_states, num_threads)
	all_warmup_reqs := make([][]PerfRequestWithTime, num_threads)
	all_perf_reqs := make([][]PerfRequestWithTime, num_threads)
	p.GenRequests(all_warmup_reqs, thread_states, num_threads, target_mops, warmup_us)
	p.GenRequests(all_perf_reqs, thread_states, num_threads, target_mops, duration_us)
	p.Benchmark(all_warmup_reqs, thread_states, num_threads, nil)
	// barrier?
	p.traces = p.Benchmark(all_perf_reqs, thread_states, num_threads, &miss_ddl_thresh_us)
	var real_duration_us uint64 = 0
	for _, trace := range p.traces {
		end_us := trace.start_us + trace.duration
		if end_us > real_duration_us {
			real_duration_us = end_us
		}
	}
	p.real_mops_ = float64(len(p.traces)) / float64(real_duration_us)
}

// get current time in microsecond resolution.
func microtime() uint64 {
	return uint64(time.Now().UnixNano() / 1000)
}

// lround rounds a floating-point number to the nearest integer.
func lround(x float64) int64 {
	return int64(math.Round(x))
}
