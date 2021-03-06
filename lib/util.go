package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
