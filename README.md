# Go daemon for Pimoroni Scroll pHat

With this daemon, you can scroll strings on your [Pimoroni](https://shop.pimoroni.com/) [Scroll pHat](https://shop.pimoroni.com/products/scroll-phat) through command line or Telegram messages.

## Install

```bash
$ go get -u github.com/meinside/scroll-daemon/...
```

## Configuration

### Config file

```bash
$ cd $GOPATH/src/github.com/meinside/scroll-daemon
$ cp config.json.sample config.json
$ vi config.json
```

and edit values to yours:

```
{
	"api_token": "0123456789:abcdefghijklmnopqrstuvwyz-x-0a1b2c3d4e",
	"local_port": 59991,
	"available_ids": [
		"telegram_id_1",
		"telegram_id_2",
		"telegram_id_3"
	],
	"phat_brightness": 5,
	"phat_scroll_delay": 75,
	"phat_rotate_180degrees": false,
	"telegram_monitor_interval": 1,
	"is_verbose": false
}
```

### systemd

```bash
$ sudo cp systemd/scroll-daemon.service /lib/systemd/system/
$ sudo vi /lib/systemd/system/scroll-daemon.service
```

and edit **User**, **Group**, **WorkingDirectory** and **ExecStart** values.

It will launch automatically on boot with:

```bash
$ sudo systemctl enable scroll-daemon.service
```

and will start with:

```bash
$ sudo systemctl start scroll-daemon.service
```

## Send messages/commands through Telegram

![scroll_daemon_telegram](https://cloud.githubusercontent.com/assets/185988/12233597/63451f38-b8ab-11e5-8aa8-f90c8023698c.png)

## Send messages/commands through command line

```bash
# scroll current time
$ scroll-cli -t

# scroll ip addresses
$ scroll-cli -i

# scroll a long message
$ scroll-cli "show me the money $$$"
$ scroll-cli `uname -a`
```


## License

MIT

