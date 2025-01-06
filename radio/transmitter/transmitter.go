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
	pin    machine.Pin
}

func NewTransmitter(pin machine.Pin, config *Config) Transmitter {
	if config == nil {
		config = DefaultConfig()
	}
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return Transmitter{*config, pin}
}

func (t Transmitter) SendCode(code int64, protocol Informative) {
	binaryView := strconv.FormatInt(code, 2)
	for range t.config.requiredRepeatCount {
		for _, b := range binaryView {
			switch b {
			case '0':
				t.pulse(protocol, protocol.Zero())
			default:
				t.pulse(protocol, protocol.One())
			}
		}
		t.pulse(protocol, protocol.SyncFactor())
	}
	t.pin.Low()
}

func (t Transmitter) pulse(protocol Informative, desc protocol.HighLow) {
	switch protocol.Inverted() {
	case true:
		t.pin.Low()
	default:
		t.pin.High()
	}
	time.Sleep(time.Duration(protocol.PulseLength()*desc.High) * time.Microsecond)
	switch protocol.Inverted() {
	case true:
		t.pin.High()
	default:
		t.pin.Low()
	}
	time.Sleep(time.Duration(protocol.PulseLength()*desc.Low) * time.Microsecond)
}
