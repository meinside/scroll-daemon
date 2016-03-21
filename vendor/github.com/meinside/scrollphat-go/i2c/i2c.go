// referenced: https://github.com/davecheney/i2c/blob/master/i2c.go
package i2c

import (
	"os"
	"syscall"
)

const (
	i2cSlave = 0x0703
)

type I2C struct {
	rc *os.File
}

// Open a new I2C connection for newer Raspberry Pis (use "/dev/i2c-0" for older ones)
func New(addr uint8) (*I2C, error) {
	if f, err := os.OpenFile("/dev/i2c-1", os.O_RDWR, 0600); err != nil {
		return nil, err
	} else {
		if err := ioctl(f.Fd(), i2cSlave, uintptr(addr)); err != nil {
			return nil, err
		}
		return &I2C{f}, nil
	}
}

// Write buf bytes to i2c
func (i2c *I2C) Write(buf []byte) (int, error) {
	return i2c.rc.Write(buf)
}

// Close
func (i2c *I2C) Close() error {
	return i2c.rc.Close()
}

func ioctl(fd, cmd, arg uintptr) (err error) {
	if _, _, e := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0); e != 0 {
		err = e
	}
	return err
}
