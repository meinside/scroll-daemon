package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	ConfigFilename = "../config.json"

	TimeFormat = "%H:%M" // format for 'date' command
)

// struct for config file
type Config struct {
	ApiToken                string   `json:"api_token"`
	LocalPort               int      `json:"local_port"`
	AvailableIds            []string `json:"available_ids"`
	PHatBrightness          byte     `json:"phat_brightness"`
	PHatScrollDelay         uint     `json:"phat_scroll_delay"` // in msec
	PHatRotate180Degrees    bool     `json:"phat_rotate_180degrees"`
	TelegramMonitorInterval uint     `json:"telegram_monitor_interval"` // in sec
	IsVerbose               bool     `json:"is_verbose"`
}

// Read config
func GetConfig() (config Config, err error) {
	_, filename, _, _ := runtime.Caller(0) // = __FILE__

	if file, err := ioutil.ReadFile(filepath.Join(path.Dir(filename), ConfigFilename)); err == nil {
		var conf Config
		if err := json.Unmarshal(file, &conf); err == nil {
			return conf, nil
		} else {
			return Config{}, err
		}
	} else {
		return Config{}, err
	}
}

// Get current time
func GetTimeString() string {
	if out, err := exec.Command("date", fmt.Sprintf("+%s", TimeFormat)).Output(); err == nil {
		return strings.TrimSpace(string(out))
	}
	return ""
}

// Get current IPs
func GetIPString() string {
	return strings.Join(getIPStringAddresses(), ", ")
}

// Get IP addresses
//
// http://play.golang.org/p/BDt3qEQ_2H
func getIPStringAddresses() []string {
	addrStrings := []string{}
	if ifaces, err := net.Interfaces(); err == nil {
		for _, iface := range ifaces {
			// skip
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}

			if addrs, err := iface.Addrs(); err == nil {
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					if ip == nil || ip.IsLoopback() {
						continue
					}
					ip = ip.To4()
					if ip == nil {
						continue
					}

					addrStrings = append(addrStrings, ip.String())
				}
			}
		}
	}

	return addrStrings
}
