package receiver

import (
	"time"

	"github.com/apolon13/TinyGoRadio/radio/protocol"
	"github.com/apolon13/TinyGoRadio/radio/protocol/receive"
)

const (
	MinDuration            = 10000
	DiffBetweenTwoTransmit = 200
	RequiredRepeatCount    = 2
)

type Decoder interface {
	Decode(timings []int64) int64
}

type Config struct {
	minDuration            int64
	diffBetweenTwoTransmit int64
	requiredRepeatCount    int8
}

func DefaultConfig() *Config {
	return &Config{
		minDuration:            MinDuration,
		diffBetweenTwoTransmit: DiffBetweenTwoTransmit,
		requiredRepeatCount:    RequiredRepeatCount,
	}
}

func NewConfig(minDuration, diffBetweenTwoTransmit int64, requiredRepeatCount int8) *Config {
	return &Config{
		minDuration:            minDuration,
		diffBetweenTwoTransmit: diffBetweenTwoTransmit,
		requiredRepeatCount:    requiredRepeatCount,
	}
}

type Receiver struct {
	repeatCount   int8
	changesCount  int8
	lastInterrupt int64
	timings       []int64
	protocols     []Decoder
	config        Config
}

func NewDefaultReceiver(config *Config) Receiver {

	if config == nil {
		config = DefaultConfig()
	}

	return Receiver{
		timings: make([]int64, receive.MaxChangesCount),
		config:  *config,
		protocols: []Decoder{
			receive.New(protocol.HighLow{High: 1, Low: 31}, protocol.HighLow{High: 1, Low: 3}, protocol.HighLow{High: 3, Low: 1}, false),
			receive.New(protocol.HighLow{High: 1, Low: 10}, protocol.HighLow{High: 1, Low: 2}, protocol.HighLow{High: 2, Low: 1}, false),
			receive.New(protocol.HighLow{High: 30, Low: 71}, protocol.HighLow{High: 4, Low: 11}, protocol.HighLow{High: 9, Low: 6}, false),
			receive.New(protocol.HighLow{High: 1, Low: 6}, protocol.HighLow{High: 1, Low: 3}, protocol.HighLow{High: 3, Low: 1}, false),
			receive.New(protocol.HighLow{High: 6, Low: 14}, protocol.HighLow{High: 1, Low: 2}, protocol.HighLow{High: 2, Low: 1}, false),
			receive.New(protocol.HighLow{High: 23, Low: 1}, protocol.HighLow{High: 1, Low: 2}, protocol.HighLow{High: 2, Low: 1}, true),
			receive.New(protocol.HighLow{High: 2, Low: 62}, protocol.HighLow{High: 1, Low: 6}, protocol.HighLow{High: 6, Low: 1}, false),
			receive.New(protocol.HighLow{High: 3, Low: 130}, protocol.HighLow{High: 7, Low: 16}, protocol.HighLow{High: 3, Low: 16}, false),
			receive.New(protocol.HighLow{High: 130, Low: 7}, protocol.HighLow{High: 16, Low: 7}, protocol.HighLow{High: 16, Low: 3}, true),
			receive.New(protocol.HighLow{High: 18, Low: 1}, protocol.HighLow{High: 3, Low: 1}, protocol.HighLow{High: 1, Low: 3}, true),
			receive.New(protocol.HighLow{High: 36, Low: 1}, protocol.HighLow{High: 1, Low: 2}, protocol.HighLow{High: 2, Low: 1}, true),
			receive.New(protocol.HighLow{High: 36, Low: 1}, protocol.HighLow{High: 1, Low: 2}, protocol.HighLow{High: 2, Low: 1}, true),
		},
	}
}

func NewReceiverWithProtocols(protocols []Decoder, config *Config) Receiver {
	if config == nil {
		config = DefaultConfig()
	}

	return Receiver{
		protocols: protocols,
		config:    *config,
		timings:   make([]int64, receive.MaxChangesCount),
	}
}

func (r *Receiver) handleTimestamp(timestamp int64) int64 {
	duration := r.durationBetweenEvents(timestamp)
	var code int64

	//Long prepare signal
	if duration > r.config.minDuration {
		if r.repeatCount == 0 || (duration-r.timings[receive.PrepareSignalIndex]) < r.config.diffBetweenTwoTransmit {
			// This long signal is close in length to the long signal which
			// started the previously recorded timings; this suggests that
			// it may indeed by a a gap between two transmissions (we assume
			// here that a sender will send the signal multiple times,
			// with roughly the same gap between them).
			r.repeatCount++
			if r.repeatCount == r.config.requiredRepeatCount && r.changesCount > receive.MinChangesCount {
				for p := 0; p < len(r.protocols); p++ {
					code = r.protocols[p].Decode(r.timings[:r.changesCount])
					if code != 0 {
						//Receive succeeded for protocol p
						break
					}
				}
				//Reset repeat counter, wait new signals with MinDuration
				r.resetRepeat()
			}
		}

		if r.repeatCount > r.config.requiredRepeatCount {
			r.resetRepeat()
		}

		//Clear counter if receive long signal
		r.resetChanges()
	}

	//detect slice of timings overflow
	if r.changesCount >= receive.MaxChangesCount {
		r.resetChanges()
		r.resetRepeat()
	}

	r.newChange(duration)
	r.updateLastEventTime(timestamp)
	return code
}

func (r *Receiver) Listen() int64 {
	return r.handleTimestamp(time.Now().UnixMicro())
}

func (r *Receiver) newChange(duration int64) {
	r.timings[r.changesCount] = duration
	r.changesCount++
}

func (r *Receiver) resetRepeat() {
	r.repeatCount = 0
}

func (r *Receiver) resetChanges() {
	r.changesCount = 0
}

func (r *Receiver) durationBetweenEvents(timestamp int64) int64 {
	return timestamp - r.lastInterrupt
}

func (r *Receiver) updateLastEventTime(timestamp int64) {
	r.lastInterrupt = timestamp
}
