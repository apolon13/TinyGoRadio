## Description

Build Your Own 315/433 MHz AM Receiver and Transmitter on TinyGo and golang.
The library is based on the [rс-switch](https://github.com/sui77/rc-switch), but with some minor modifications.
Currently it only supports the receiver mode.

## Install

```sh
go get github.com/apolon13/TinyGoRadio/radio
```

## Examples

### Simple usage
```go
package main

import (
	"github.com/apolon13/TinyGoRadio/radio/receiver"
	"machine"
)

func main() {
	r := receiver.NewDefaultReceiver(nil)
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

func (p CustomProtocol) Decode(timings []int64) int64 {
	//handle your timings and return code
	return 0
}

func main() {
	customProtocol := CustomProtocol{}
	r := receiver.NewReceiverWithProtocols([]receiver.Decoder{customProtocol}, nil)
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

	r := receiver.NewDefaultReceiver(receiver.NewConfig(minStartSignalDuration, diffBetweenTwoTransmit, requiredRepeatCount))
	pin := machine.GPIO6
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	pin.SetInterrupt(machine.PinToggle, func(pin machine.Pin) {
		if code := r.Listen(); code != 0 {
			println(code)
		}
	})
}
```