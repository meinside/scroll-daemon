package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/meinside/scroll-daemon/lib"
)

var showTime bool
var showIp bool

func init() {
	flag.BoolVar(&showTime, "t", false, "show current time")
	flag.BoolVar(&showIp, "i", false, "show ip addresses")
}

func printUsage() {
	fmt.Printf(`* Usage:

	# scroll a message
	$ %[1]s [strings to show]

	# scroll current time
	$ %[1]s -t

	# scroll ip addresses
	$ %[1]s -i
`, os.Args[0])
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		// read variables from config file
		if conf, err := lib.GetConfig(); err == nil {
			localPort := conf.LocalPort

			flag.Parse()

			var cmd, val string
			if showTime {
				cmd, val = lib.CommandTime, ""
			} else if showIp {
				cmd, val = lib.CommandIP, ""
			} else {
				cmd, val = "", strings.Join(args, " ")
			}

			if _, err := http.PostForm(fmt.Sprintf("http://localhost:%d", localPort), url.Values{
				lib.ParamCommand: {cmd},
				lib.ParamValue:   {val},
			}); err != nil {
				fmt.Println(fmt.Errorf("*** %s", err))
			}
		} else {
			fmt.Println(fmt.Errorf("*** %s", err))
		}
	} else {
		printUsage()
	}
}
