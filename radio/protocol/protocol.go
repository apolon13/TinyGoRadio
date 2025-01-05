package protocol

// HighLow is description of a single pule, which consists of a high signal
// whose duration is "high" times the base pulse length, followed
// by a low signal lasting "low" times the base pulse length.
// Thus, the pulse overall lasts (high+low)*pulseLength
type HighLow struct {
	High int16
	Low  int16
}

// Protocol describes how zero and one bits are encoded into high/low pulses
type Protocol struct {
	pulseLength int16
	syncFactor  HighLow
	zero        HighLow
	one         HighLow
	inverted    bool
}

func New(pulseLength int16, syncFactor, zero, one HighLow, inverted bool) Protocol {
	return Protocol{
		pulseLength,
		syncFactor,
		zero,
		one,
		inverted,
	}
}

func (p Protocol) SyncFactor() HighLow {
	return p.syncFactor
}

func (p Protocol) Zero() HighLow {
	return p.zero
}

func (p Protocol) One() HighLow {
	return p.one
}

func (p Protocol) Inverted() bool {
	return p.inverted
}

func (p Protocol) PulseLength() int16 {
	return p.pulseLength
}

func (p Protocol) delimiter(startSignalTiming int64) int64 {
	syncLengthInPulses := p.syncFactor.High
	if p.syncFactor.Low > p.syncFactor.High {
		syncLengthInPulses = p.syncFactor.Low
	}
	return startSignalTiming / int64(syncLengthInPulses)
}

func (p Protocol) isZero(timing, next, delimiter, tolerance int64) bool {
	return (timing-(delimiter*int64(p.zero.High))) < tolerance && (next-delimiter*int64(p.zero.Low)) < tolerance
}

func (p Protocol) isOne(timing, next, delimiter, tolerance int64) bool {
	return (timing-(delimiter*int64(p.one.High))) < tolerance && (next-delimiter*int64(p.one.Low)) < tolerance
}

func (p Protocol) Decode(timings []int64) int64 {
	delimiter := p.delimiter(timings[0])
	tolerance := delimiter * 60 / 100
	var code int64 = 0
	firstData := 1
	if p.inverted {
		firstData = 2
	}

	for i := firstData; i < len(timings)-1; i += 2 {
		timing := timings[i]
		next := timings[i+1]
		code <<= 1
		if p.isZero(timing, next, delimiter, tolerance) {
		} else if p.isOne(timing, next, delimiter, tolerance) {
			code |= 1
		} else {
			return 0
		}
	}

	return code
}
