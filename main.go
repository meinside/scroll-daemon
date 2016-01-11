package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/meinside/scroll-daemon/lib"

	scroll "github.com/meinside/scrollphat-go"
	bot "github.com/meinside/telegram-bot-go"
)

const (
	HttpCommandPath = "/"

	QueueSize              = 1
	MonitorIntervalSeconds = 3

	ScrollDelayMsecs = 50
)

// variables
var pHat *scroll.ScrollPHat
var apiToken string
var localPort int
var availableIds []string
var pHatBrightness byte
var isVerbose bool
var queue chan string

func init() {
	// read variables from config file
	if conf, err := lib.GetConfig(); err == nil {
		apiToken = conf.ApiToken
		localPort = conf.LocalPort
		availableIds = conf.AvailableIds
		pHatBrightness = conf.PHatBrightness
		isVerbose = conf.IsVerbose

		// initialize other variables
		queue = make(chan string, QueueSize)
	} else {
		panic(err.Error())
	}
}

// check if given Telegram id is available
func isAvailableId(id string) bool {
	for _, v := range availableIds {
		if v == id {
			return true
		}
	}
	return false
}

// for processing incoming update from Telegram
func processUpdate(client *bot.Bot, update bot.Update) bool {
	// check username
	var userId string
	if update.Message.From.Username == nil {
		log.Printf("*** Not allowed (no user name): %s\n", *update.Message.From.FirstName)
		return false
	}
	userId = *update.Message.From.Username
	if !isAvailableId(userId) {
		log.Printf("*** Id not allowed: %s\n", userId)
		return false
	}

	if update.HasMessage() {
		txt := *update.Message.Text

		if isVerbose {
			log.Printf("received telegram message: %s\n", txt)
		}

		if strings.HasPrefix(txt, "/") { // if it is command,
			if strings.HasPrefix(txt, "/"+lib.CommandStart) {
				message := lib.MessageStart
				var options map[string]interface{} = map[string]interface{}{
					"reply_markup": bot.ReplyKeyboardMarkup{
						Keyboard: [][]string{
							[]string{"/" + lib.CommandTime, "/" + lib.CommandIP},
						},
					},
				}

				// send message
				if sent := client.SendMessage(update.Message.Chat.Id, &message, options); !sent.Ok {
					log.Printf("*** Failed to send message: %s\n", *sent.Description)
				}
			} else if strings.HasPrefix(txt, "/"+lib.CommandTime) {
				queue <- lib.GetTimeString()
			} else if strings.HasPrefix(txt, "/"+lib.CommandIP) {
				queue <- lib.GetIPString()
			} else {
				log.Printf("*** No such command: %s\n", txt)
			}
		} else { // otherwise,
			queue <- txt
		}
	}

	return false
}

// for processing incoming request through HTTP
var httpHandler = func(w http.ResponseWriter, r *http.Request) {
	command := r.FormValue(lib.ParamCommand)
	value := r.FormValue(lib.ParamValue)

	if isVerbose {
		log.Printf("received http command: %s, value: %s\n", command, value)
	}

	if command != "" { // if it is command,
		if command == lib.CommandTime {
			queue <- lib.GetTimeString()
		} else if command == lib.CommandIP {
			queue <- lib.GetIPString()
		} else {
			log.Printf("*** No such command: %s\n", command)
		}
	} else { // otherwise,
		queue <- value
	}
}

func main() {
	client := bot.NewClient(apiToken)
	client.Verbose = isVerbose

	pHat = scroll.New()
	if pHat == nil {
		panic("Failed to initialize Scroll pHat")
	} else {
		// setup pHat
		pHat.SetBrightness(pHatBrightness)

		// wait for channel
		go func() {
			for {
				select {
				case s := <-queue:
					pHat.Scroll(s, ScrollDelayMsecs)
				default:
					// do nothing
				}
			}
		}()

		// start web server
		go func() {
			log.Printf("Starting local web server on port: %d\n", localPort)

			http.HandleFunc(HttpCommandPath, httpHandler)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", localPort), nil); err != nil {
				panic(err.Error())
			}
		}()

		// monitor for new telegram updates
		if me := client.GetMe(); me.Ok { // get info about this bot
			log.Printf("Launching bot: @%s (%s)\n", *me.Result.Username, *me.Result.FirstName)

			// delete webhook (getting updates will not work when wehbook is set up)
			if unhooked := client.DeleteWebhook(); unhooked.Ok {
				// wait for new updates
				client.StartMonitoringUpdates(0, MonitorIntervalSeconds, func(update bot.Update, err error) {
					if err == nil {
						if update.Message != nil {
							processUpdate(client, update)
						}
					} else {
						log.Printf("*** Error while receiving update (%s)\n", err.Error())
					}
				})
			} else {
				panic("Failed to delete webhook")
			}
		} else {
			panic("Failed to get info of the bot")
		}
	}
}
