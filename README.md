## Description

Build Your Own 315/433 MHz AM Receiver and Transmitter with TinyGo.
The library y is based on the [rс-switch](https://github.com/sui77/rc-switch)
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
	"github.com/apolon13/TinyGoRadio/radio"
	"machine"
)

func main() {
	receiver := radio.NewDefaultReceiver(nil)
	pin := machine.GPIO6
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	pin.SetInterrupt(machine.PinToggle, func(pin machine.Pin) {
		if code := receiver.Listen(); code != 0 {
			println(code)
		}
	})
}
```

### Create your custom protocol
```go
package main

import (
	"github.com/apolon13/TinyGoRadio/radio"
	"machine"
)

type CustomProtocol struct {

}

func (p CustomProtocol) Decode(timings []int64) int64  {
	//handle your timings and return code
	return 0
}

func main() {
	customProtocol := CustomProtocol{}
	receiver := radio.NewReceiverWithProtocols([]radio.Protocol{customProtocol}, nil)

	pin := machine.GPIO6
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	pin.SetInterrupt(machine.PinToggle, func(pin machine.Pin) {
		if code := receiver.Listen(); code != 0 {
			println(code)
		}
	})
}
```

### Custom receiver config
```go
package main

import (
	"github.com/apolon13/TinyGoRadio/radio"
	"machine"
)


func main() {
	var minStartSignalDuration int64 = 1000
	var diffBetweenTwoTransmit int64 = 300
	var requiredRepeatCount int8 = 3

	receiver := radio.NewDefaultReceiver(radio.NewConfig(minStartSignalDuration, diffBetweenTwoTransmit, requiredRepeatCount), )
	pin := machine.GPIO6
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	pin.SetInterrupt(machine.PinToggle, func(pin machine.Pin) {
		if code := receiver.Listen(); code != 0 {
			println(code)
		}
	})
}

```