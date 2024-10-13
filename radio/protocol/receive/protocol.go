package receive

import "github.com/apolon13/TinyGoRadio/radio/protocol"

const (
	MaxChangesCount    = 67
	MinChangesCount    = 7
	PrepareSignalIndex = 0
)

// Protocol describes how zero and one bits are encoded into high/low pulses
type Protocol struct {
	syncFactor     protocol.HighLow
	zero           protocol.HighLow
	one            protocol.HighLow
	invertedSignal bool
}

func New(syncFactor, zero, one protocol.HighLow, invertedSignal bool) Protocol {
	return Protocol{
		syncFactor,
		zero,
		one,
		invertedSignal,
	}
}

func (p Protocol) delay(startSignalTiming int64) int64 {
	syncLengthInPulses := p.syncFactor.High
	if p.syncFactor.Low > p.syncFactor.High {
		syncLengthInPulses = p.syncFactor.Low
	}

	return startSignalTiming / int64(syncLengthInPulses)
}

func (p Protocol) isZero(timing, next, delay, tolerance int64) bool {
	return (timing-(delay*int64(p.zero.High))) < tolerance && (next-delay*int64(p.zero.Low)) < tolerance
}

func (p Protocol) isOne(timing, next, delay, tolerance int64) bool {
	return (timing-(delay*int64(p.one.High))) < tolerance && (next-delay*int64(p.one.Low)) < tolerance
}

func (p Protocol) Decode(timings []int64) int64 {
	delay := p.delay(timings[PrepareSignalIndex])
	tolerance := delay * 60 / 100
	var code int64 = 0
	firstData := 1
	if p.invertedSignal {
		firstData = 2
	}

	for i := firstData; i < len(timings)-1; i += 2 {
		timing := timings[i]
		next := timings[i+1]
		code <<= 1
		if p.isZero(timing, next, delay, tolerance) {
		} else if p.isOne(timing, next, delay, tolerance) {
			code |= 1
		} else {
			return 0
		}
	}

	return code
}
