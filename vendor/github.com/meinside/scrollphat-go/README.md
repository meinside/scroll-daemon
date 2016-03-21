# Go library for Pimoroni Scroll pHat

Go library for controlling [Pimoroni](https://shop.pimoroni.com/)'s [Scroll pHat](https://shop.pimoroni.com/products/scroll-phat).

Most of the codes were inspired by [scroll-phat](https://github.com/pimoroni/scroll-phat) and [Dave Cheney](http://dave.cheney.net/)'s [i2c](https://github.com/davecheney/i2c/blob/master/i2c.go).

## Install

```bash
$ go get -u github.com/meinside/scrollphat-go
```

## Configuration

You need to enable I2C first:

```bash
$ sudo raspi-config
```

Also, your account should be in 'i2c' group:

```bash
$ sudo usermod -G i2c USERNAME
```

## Sample

```go
// sample.go - sample application for Scroll pHat, 2016.01.07.
package main

import (
	"time"

	scroll "github.com/meinside/scrollphat-go"
)

func main() {
	//scroll.TrimRedundantSpaces = false	// don't trim redundant spaces

	if phat := scroll.New(); phat != nil {
		// set brightneess
		phat.SetBrightness(5)
		time.Sleep(1 * time.Second)

		// turn on all LEDs
		phat.Fill()
		time.Sleep(1 * time.Second)

		// turn off all LEDs
		phat.Clear()
		time.Sleep(1 * time.Second)

		// draw a string (not scrolling)
		phat.DrawBytes(scroll.BytesForString(":-)"))
		time.Sleep(1 * time.Second)

		// draw a string (scrolling)
		phat.Scroll(" 0123456789 ABCDEFGHIJKLMNOPQRSTUVWXYZ abcdefghijklmnopqrstuvwxyz `~!@#$%^&*()-=_+[]{};:\"'<>,.?/\\|", 50)
		time.Sleep(1 * time.Second)

		// draw points (heart-shaped...)
		phat.DrawBytes(scroll.ArraysToBytes([][]uint8{
			{0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0},
			{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
			{0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0},
			{0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		}))
		time.Sleep(1 * time.Second)

		// invert LEDs
		phat.Invert()
		time.Sleep(1 * time.Second)

		// turn off all LEDs
		phat.Clear()

		phat.Close()
	} else {
		panic("could not initialize ScrollPHat")
	}
}
```

Compile it and run like this:

```bash
$ go build sample.go
$ ./sample
```

Then you'll see:

![scroll_phat_sample.gif](https://cloud.githubusercontent.com/assets/185988/12164502/5031a09c-b559-11e5-9c34-c66c502183f6.gif)

## License

MIT

