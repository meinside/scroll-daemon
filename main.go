package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/meinside/scroll-daemon/lib"

	"github.com/meinside/rpi-tools/status"

	scroll "github.com/meinside/scrollphat-go"
	bot "github.com/meinside/telegram-bot-go"
)

const (
	HttpCommandPath = "/"

	QueueSize = 1

	LocationLivePeriodSeconds = 60
)

// variables
var pHat *scroll.ScrollPHat
var apiToken string
var localPort int
var availableIds []string
var pHatBrightness byte
var pHatScrollDelay uint
var pHatRotateOrNot bool
var telegramMonitorInterval uint
var isVerbose bool
var queue chan string
var allKeyboards [][]bot.KeyboardButton

func init() {
	// read variables from config file
	if conf, err := lib.GetConfig(); err == nil {
		apiToken = conf.ApiToken
		localPort = conf.LocalPort
		availableIds = conf.AvailableIds
		pHatBrightness = conf.PHatBrightness
		pHatScrollDelay = conf.PHatScrollDelay
		pHatRotateOrNot = conf.PHatRotate180Degrees
		telegramMonitorInterval = conf.TelegramMonitorInterval
		isVerbose = conf.IsVerbose

		// all keyboard buttons
		allKeyboards = [][]bot.KeyboardButton{
			[]bot.KeyboardButton{
				bot.KeyboardButton{
					Text: "/" + lib.CommandTime,
				},
				bot.KeyboardButton{
					Text: "/" + lib.CommandIP,
				},
				bot.KeyboardButton{
					Text: "/" + lib.CommandLocation,
				},
				bot.KeyboardButton{
					Text: "/" + lib.CommandHelp,
				},
			},
		}

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
func processUpdate(b *bot.Bot, update bot.Update) bool {
	// check username
	var userId string
	if update.Message.From.Username == nil {
		log.Printf("*** Not allowed (no user name): %s", update.Message.From.FirstName)
		return false
	}
	userId = *update.Message.From.Username
	if !isAvailableId(userId) {
		log.Printf("*** Id not allowed: %s", userId)
		return false
	}

	if update.HasMessage() {
		txt := *update.Message.Text

		if isVerbose {
			log.Printf("received telegram message: %s", txt)
		}

		if strings.HasPrefix(txt, "/") { // if it is command,
			var options map[string]interface{} = map[string]interface{}{
				"reply_markup": bot.ReplyKeyboardMarkup{
					Keyboard:       allKeyboards,
					ResizeKeyboard: true,
				},
			}

			if strings.HasPrefix(txt, "/"+lib.CommandStart) {
				message := lib.MessageStart

				// send message
				if sent := b.SendMessage(update.Message.Chat.Id, message, options); !sent.Ok {
					log.Printf("*** Failed to send message: %s", *sent.Description)
				}
			} else if strings.HasPrefix(txt, "/"+lib.CommandTime) {
				time := lib.GetTimeString()

				queue <- time

				if sent := b.SendMessage(update.Message.Chat.Id, time, options); !sent.Ok {
					log.Printf("*** Failed to send message: %s", *sent.Description)
				}
			} else if strings.HasPrefix(txt, "/"+lib.CommandIP) {
				ip := strings.Join(status.IpAddresses(), ", ")

				queue <- ip

				if sent := b.SendMessage(update.Message.Chat.Id, ip, options); !sent.Ok {
					log.Printf("*** Failed to send message: %s", *sent.Description)
				}
			} else if strings.HasPrefix(txt, "/"+lib.CommandLocation) {
				if extIp, err := status.ExternalIpAddress(); err == nil {
					if geoInfo, err := status.GeoLocation(extIp); err == nil {
						options["live_period"] = LocationLivePeriodSeconds

						if sent := b.SendLocation(
							update.Message.Chat.Id,
							geoInfo.Location.Latitude,
							geoInfo.Location.Longitude,
							options,
						); !sent.Ok {
							log.Printf("*** Failed to send location: %s", *sent.Description)
						}
					} else {
						if sent := b.SendMessage(
							update.Message.Chat.Id,
							fmt.Sprintf("Failed to get geo location: %s", err),
							options,
						); !sent.Ok {
							log.Printf("*** Failed to send message: %s", *sent.Description)
						}
					}
				} else {
					if sent := b.SendMessage(
						update.Message.Chat.Id,
						fmt.Sprintf("Failed to get external ip address: %s", err),
						options,
					); !sent.Ok {
						log.Printf("*** Failed to send message: %s", *sent.Description)
					}
				}
			} else if strings.HasPrefix(txt, "/"+lib.CommandHelp) {
				// send message
				message := lib.MessageHelp
				if sent := b.SendMessage(update.Message.Chat.Id, message, options); !sent.Ok {
					log.Printf("*** Failed to send message: %s", *sent.Description)
				}
			} else {
				log.Printf("*** No such command: %s", txt)

				// send message
				message := fmt.Sprintf("No such command: %s", txt)
				if sent := b.SendMessage(update.Message.Chat.Id, message, options); !sent.Ok {
					log.Printf("*** Failed to send message: %s", *sent.Description)
				}
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
		log.Printf("received http command: %s, value: %s", command, value)
	}

	if command != "" { // if it is command,
		if command == lib.CommandTime {
			queue <- lib.GetTimeString()
		} else if command == lib.CommandIP {
			queue <- strings.Join(status.IpAddresses(), ", ")
		} else {
			log.Printf("*** No such command: %s", command)
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
		scroll.IsFlippedHorizontally = pHatRotateOrNot
		scroll.IsFlippedVertically = pHatRotateOrNot
		pHat.SetBrightness(pHatBrightness)

		// wait for channel
		go func() {
			for {
				select {
				case s := <-queue:
					pHat.Scroll(s, pHatScrollDelay)
				}
			}
		}()

		// start web server
		go func() {
			log.Printf("Starting local web server on port: %d", localPort)

			http.HandleFunc(HttpCommandPath, httpHandler)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", localPort), nil); err != nil {
				panic(err.Error())
			}
		}()

		// monitor for new telegram updates
		if me := client.GetMe(); me.Ok { // get info about this bot
			log.Printf("Launching bot: @%s (%s)", *me.Result.Username, me.Result.FirstName)

			// delete webhook (getting updates will not work when wehbook is set up)
			if unhooked := client.DeleteWebhook(); unhooked.Ok {
				// wait for new updates
				client.StartMonitoringUpdates(0, int(telegramMonitorInterval), func(b *bot.Bot, update bot.Update, err error) {
					if err == nil {
						if update.Message != nil {
							processUpdate(b, update)
						}
					} else {
						log.Printf("*** Error while receiving update (%s)", err.Error())
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
