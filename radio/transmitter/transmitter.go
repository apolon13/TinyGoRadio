package transmitter

import (
	"machine"
	"strconv"
	"time"

	"github.com/apolon13/TinyGoRadio/radio/protocol"
)

type Informative interface {
	One() protocol.HighLow
	Zero() protocol.HighLow
	SyncFactor() protocol.HighLow
	Inverted() bool
	PulseLength() int16
}

const (
	requiredRepeatCount = 2
)

type Config struct {
	requiredRepeatCount int8
}

func NewConfig(requiredRepeatCount int8) Config {
	return Config{requiredRepeatCount: requiredRepeatCount}
}

func DefaultConfig() *Config {
	return &Config{
		requiredRepeatCount: requiredRepeatCount,
	}
}

type Transmitter struct {
	config Config
}

func NewTransmitter(config *Config) Transmitter {
	if config == nil {
		config = DefaultConfig()
	}
	return Transmitter{*config}
}

func (t Transmitter) Send(code int, pin machine.Pin, protocol Informative) {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	binaryView := strconv.FormatInt(int64(code), 2)
	for range t.config.requiredRepeatCount {
		for _, b := range binaryView {
			switch b {
			case '0':
				t.transmit(pin, protocol, protocol.Zero())
			default:
				t.transmit(pin, protocol, protocol.One())
			}
		}
		t.transmit(pin, protocol, protocol.SyncFactor())
	}
	pin.Low()
}

func (t Transmitter) transmit(pin machine.Pin, protocol Informative, pulse protocol.HighLow) {
	switch protocol.Inverted() {
	case true:
		pin.Low()
	default:
		pin.High()
	}
	time.Sleep(time.Duration(protocol.PulseLength()*pulse.High) * time.Microsecond)
	switch protocol.Inverted() {
	case true:
		pin.High()
	default:
		pin.Low()
	}
	time.Sleep(time.Duration(protocol.PulseLength()*pulse.Low) * time.Microsecond)
}
