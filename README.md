## Description

Build Your Own 315/433 MHz AM Receiver and Transmitter with TinyGo.
The librar y is based on the [rс-switch](https://github.com/sui77/rc-switch)
Currently it only supports the receiver mode.

## Install

```sh
go get github.com/apolon13/TinyGoRadio/radio
```

## Examples

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