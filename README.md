## Description

Build Your Own 315/433 MHz AM Receiver and Transmitter on TinyGo and golang.
The library is based on the [rс-switch](https://github.com/sui77/rc-switch), but with some minor modifications.
Supports the receiver and transmit mode.

## Install

```sh
go get github.com/apolon13/TinyGoRadio/radio
```

## Examples

### Shortest receiver
```go
package main

import (
	"machine"

	"github.com/apolon13/TinyGoRadio/radio/receiver"
)

func main() {
	r := receiver.NewReceiver(nil)
	r.ListenByPin(machine.GPIO6, func(code int64) {
		println(code)
	})
}
```

### Simple usage
```go
package main

import (
	"machine"

	"github.com/apolon13/TinyGoRadio/radio/receiver"
)

func main() {
	r := receiver.NewReceiver(nil)
	pin := machine.GPIO6
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	pin.SetInterrupt(machine.PinToggle, func(pin machine.Pin) {
		if code := r.Listen(); code != 0 {
			println(code)
		}
	})
}
```

### Create your custom protocol decoder
```go
package main

import (
	"github.com/apolon13/TinyGoRadio/radio/receiver"
	
	"machine"
)

type CustomProtocol struct {
}

func (c CustomProtocol) Decode(timings []int64) int64 {
	//handle your timings and return code
	return 0
}

func main() {
	customProtocol := CustomProtocol{}
	r := receiver.NewReceiverWithProtocols([]receiver.Decodable{customProtocol}, nil)
	pin := machine.GPIO6
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	pin.SetInterrupt(machine.PinToggle, func(pin machine.Pin) {
		if code := r.Listen(); code != 0 {
			println(code)
		}
	})
}

```

### Custom receiver config
```go
package main

import (
	"github.com/apolon13/TinyGoRadio/radio/receiver"
	"machine"
)


func main() {
	var minStartSignalDuration int64 = 1000
	var diffBetweenTwoTransmit int64 = 300
	var requiredRepeatCount int8 = 3

	r := receiver.NewReceiver(receiver.NewConfig(minStartSignalDuration, diffBetweenTwoTransmit, requiredRepeatCount))
	pin := machine.GPIO6
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	pin.SetInterrupt(machine.PinToggle, func(pin machine.Pin) {
		if code := r.Listen(); code != 0 {
			println(code)
		}
	})
}
```

### Transmit data
```go
package main

import (
	"machine"
	"time"

	"github.com/apolon13/TinyGoRadio/radio/protocol"
	"github.com/apolon13/TinyGoRadio/radio/transmitter"
)


func main() {
	transmitter := transmitter.NewTransmitter(nil)
	pin := machine.GPIO6
	for {
		transmitter.Send(15858700, pin, protocol.New(
			350,
			protocol.HighLow{High: 1, Low: 31},
			protocol.HighLow{High: 1, Low: 3},
			protocol.HighLow{High: 3, Low: 1},
			false,
		))
	}
}
```