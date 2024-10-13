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