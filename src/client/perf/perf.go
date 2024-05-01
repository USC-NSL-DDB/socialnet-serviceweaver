package perf

import (
	"math"
	"math/rand"
	"sort"
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
	AbslStartUs uint64
	StartUs     uint64
	Duration    uint64
}

type TraceFormat int

const (
	kUnsorted TraceFormat = iota
	kSortedByDuration
	kSortedByStart
)

// Closed-loop, possion arrival.
type Perf struct {
	adapter      PerfAdapter
	traces       []Trace
	real_mops_   float64
	trace_format TraceFormat
}

func NewPerf(adapter PerfAdapter) *Perf {
	return &Perf{adapter: adapter}
}

func (p *Perf) Clear() {
	p.traces = p.traces[:0]
	p.real_mops_ = 0
}

func (p *Perf) GenRequests(
	all_reqs *[][]PerfRequestWithTime, thread_states []PerfThreadState, num_threads uint32, target_mops float64, duration_us uint64,
) {
	wg := sync.WaitGroup{}

	for i := uint32(0); i < num_threads; i++ {
		wg.Add(1)
		go func(reqs *[]PerfRequestWithTime, thread_state PerfThreadState) {
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

				*reqs = append(*reqs, req_with_time)
				curr_us += uint64(interval)
			}
		}(&(*all_reqs)[i], thread_states[i])
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
					AbslStartUs: now, StartUs: now - start_us, Duration: 0,
				}
				ok := p.adapter.ServeRequest(thread_state, req_with_time.req)
				trace.Duration = microtime() - start_us - trace.StartUs
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

func (p *Perf) CreateThreadStates(thread_states *[]PerfThreadState, num_threads uint32) {
	for i := uint32(0); i < num_threads; i++ {
		*thread_states = append(*thread_states, p.adapter.CreateGoroutineState())
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
	p.CreateThreadStates(&thread_states, num_threads)
	all_warmup_reqs := make([][]PerfRequestWithTime, num_threads)
	all_perf_reqs := make([][]PerfRequestWithTime, num_threads)
	p.GenRequests(&all_warmup_reqs, thread_states, num_threads, target_mops, warmup_us)
	p.GenRequests(&all_perf_reqs, thread_states, num_threads, target_mops, duration_us)
	p.Benchmark(all_warmup_reqs, thread_states, num_threads, nil)
	// barrier?
	p.traces = p.Benchmark(all_perf_reqs, thread_states, num_threads, &miss_ddl_thresh_us)
	var real_duration_us uint64 = 0
	for _, trace := range p.traces {
		end_us := trace.StartUs + trace.Duration
		if end_us > real_duration_us {
			real_duration_us = end_us
		}
	}
	p.real_mops_ = float64(len(p.traces)) / float64(real_duration_us)
}

func (p *Perf) GetAvgLat() uint64 {
	if p.trace_format != kSortedByDuration {
		sort.Slice(p.traces, func(i, j int) bool {
			return p.traces[i].Duration < p.traces[j].Duration
		})
		p.trace_format = kSortedByDuration
	}

	var sum uint64
	for _, trace := range p.traces {
		sum += trace.Duration
	}
	return sum / uint64(len(p.traces))
}

func (p *Perf) GetNthLats(nth float64) uint64 {
	if p.trace_format != kSortedByDuration {
		sort.Slice(p.traces, func(i, j int) bool {
			return p.traces[i].Duration < p.traces[j].Duration
		})
		p.trace_format = kSortedByDuration
	}

	idx := int(nth / 100.0 * float64(len(p.traces)))
	return p.traces[idx].Duration
}

func (p *Perf) GetTimeseriesNthLats(interval_us uint64, nth float64) []Trace {
	if p.trace_format != kSortedByStart {
		sort.Slice(p.traces, func(i, j int) bool {
			return p.traces[i].StartUs < p.traces[j].StartUs
		})
		p.trace_format = kSortedByStart
	}

	var timeseries []Trace
	var win_duratins []uint64
	curr_win_us := p.traces[0].StartUs
	absl_curr_win_us := p.traces[0].AbslStartUs

	for _, trace := range p.traces {
		if curr_win_us+interval_us < trace.StartUs {
			sort.Slice(win_duratins, func(i, j int) bool {
				return win_duratins[i] < win_duratins[j]
			})
			if len(win_duratins) >= 100 {
				idx := int(nth / 100.0 * float64(len(win_duratins)))
				timeseries = append(timeseries, Trace{
					AbslStartUs: absl_curr_win_us, StartUs: curr_win_us, Duration: win_duratins[idx],
				})
			}
			curr_win_us += interval_us
			absl_curr_win_us += interval_us
			win_duratins = win_duratins[:0]
		}
		win_duratins = append(win_duratins, trace.Duration)
	}

	return timeseries
}

func (p *Perf) GetRealMops() float64 {
	return p.real_mops_
}

func (p *Perf) GetTraces() []Trace {
	return p.traces
}

// get current time in microsecond resolution.
func microtime() uint64 {
	return uint64(time.Now().UnixNano() / 1000)
}

// lround rounds a floating-point number to the nearest integer.
func lround(x float64) int64 {
	return int64(math.Round(x))
}
