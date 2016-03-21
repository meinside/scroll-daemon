// Go library for Pimoroni's Scroll pHat (Rasbperry Pi)

package scroll

import (
	"log"
	"math"
	"sync"
	"time"

	i2c "github.com/meinside/scrollphat-go/i2c"
)

const (
	i2cAddr  = 0x60
	mode5x11 = 0x03

	// commands
	CmdSetMode       = 0x00
	CmdSetBrightness = 0x19
	CmdPaintPixel    = 0x01

	// 11 x 5 LEDs
	NumRows = 5
	NumCols = 11
)

const (
	minimumDelay = 10 // 10 ms
)

type ScrollPHat struct {
	sync.Mutex
	i2c *i2c.I2C

	// internal pixels array
	pixels []byte
}

var IsFlippedVertically = false   // leds are flipped vertically or not
var IsFlippedHorizontally = false // leds are flipped horizontally or not

// Get a new ScrollPHat
func New() *ScrollPHat {
	if i, err := i2c.New(i2cAddr); err == nil {
		nu := ScrollPHat{
			i2c:    i,
			pixels: make([]byte, NumCols),
		}

		if _, err := nu.write(CmdSetMode, []byte{mode5x11}); err != nil {
			log.Printf("*** failed to set mode: %s\n", err)
		}

		return &nu
	} else {
		log.Printf("*** failed to open i2c: %s\n", err)
	}

	return nil
}

// Write given command and data to i2c
func (s ScrollPHat) write(cmd byte, data []byte) (n int, err error) {
	var bytes []byte = []byte{cmd}
	bytes = append(bytes, data...)

	return s.i2c.Write(bytes)
}

// Set LED brightness
func (s ScrollPHat) SetBrightness(brightness byte) {
	s.Lock()
	defer s.Unlock()

	if _, err := s.write(CmdSetBrightness, []byte{brightness}); err != nil {
		log.Printf("[SetBrightness] error writing to i2c: %s\n", err)
	}
}

// Draw given pixel bytes
func (s ScrollPHat) draw(bytes []byte) {
	if bytes == nil {
		log.Printf("[draw] given bytes is nil")
		return
	}

	// reverse bit orders
	if IsFlippedVertically {
		var b byte
		for i := 0; i < len(bytes); i++ {
			b = bytes[i]
			bytes[i] = ((b & 0x01) << 4) | ((b & 0x02) << 2) | (b & 0x04) | ((b & 0x08) >> 2) | ((b & 0x10) >> 4)
		}
	}

	// reverse byte orders
	if IsFlippedHorizontally {
		reversed := []byte{}
		for i := len(bytes) - 1; i >= 0; i-- {
			reversed = append(reversed, bytes[i])
		}
		bytes = reversed
	}

	bytes = append(bytes, 0xFF) // need trailing 0xFF for drawing

	if _, err := s.write(CmdPaintPixel, bytes); err != nil {
		log.Printf("[draw] error writing to i2c: %s\n", err)
	}
}

// Alter pixels in the memory and draw them
func (s ScrollPHat) DrawBytes(bytes []byte) {
	s.Lock()
	defer s.Unlock()

	copy(s.pixels, bytes)
	s.draw(s.pixels)
}

// Scroll given string with delay, altering pixels in the memory for each frame
func (s ScrollPHat) Scroll(str string, delayMs uint) {
	if delayMs < minimumDelay {
		delayMs = minimumDelay
	}

	bytes := BytesForString(str)
	padding := make([]byte, NumCols-1) // append padding for scrolling from right edge
	bytes = append(padding, bytes...)

	s.Lock()
	defer s.Unlock()

	var limit int
	for offset, _ := range bytes {
		// copy bytes
		limit = int(math.Min(float64(len(bytes)-offset), float64(NumCols)))
		for i := 0; i < NumCols; i++ {
			if i < limit {
				s.pixels[i] = bytes[offset+i]
			} else {
				s.pixels[i] = 0x00
			}
		}

		s.draw(s.pixels)

		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
}

// Alter given pixel in the memory and draw them
func (s ScrollPHat) DrawPixel(col, row uint, isOn bool) {
	s.Lock()
	defer s.Unlock()

	if isOn {
		s.pixels[col] |= (1 << row)
	} else {
		s.pixels[col] ^= (1 << row)
	}
	s.draw(s.pixels)
}

// Toggle given pixel in the memory and draw them
func (s ScrollPHat) TogglePixel(col, row uint) {
	s.Lock()
	defer s.Unlock()

	if s.pixels[col]&(1<<row) > 0 {
		s.pixels[col] ^= (1 << row)
	} else {
		s.pixels[col] |= (1 << row)
	}
	s.draw(s.pixels)
}

// Clear pixels for both memory and LED
func (s ScrollPHat) Clear() {
	s.Lock()
	defer s.Unlock()

	// clear bytes
	for i := 0; i < NumCols; i++ {
		s.pixels[i] = 0x00
	}
	s.draw(s.pixels)
}

// Fill pixels for both memory and LED
func (s ScrollPHat) Fill() {
	var fill byte
	for i := 0; i < NumRows; i++ {
		fill |= (1 << uint(i))
	}

	s.Lock()
	defer s.Unlock()

	for i := 0; i < NumCols; i++ {
		s.pixels[i] = fill
	}
	s.draw(s.pixels)
}

// Invert pixels for both memory and LED
func (s ScrollPHat) Invert() {
	s.Lock()
	defer s.Unlock()

	for i := 0; i < NumCols; i++ {
		for j := 0; j < NumRows; j++ {
			if s.pixels[i]&(1<<uint(j)) > 0 {
				s.pixels[i] ^= (1 << uint(j))
			} else {
				s.pixels[i] |= (1 << uint(j))
			}
		}
	}
	s.draw(s.pixels)
}

// Close
func (s ScrollPHat) Close() error {
	s.Lock()
	defer s.Unlock()

	return s.i2c.Close()
}
