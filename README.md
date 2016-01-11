# Go daemon for Pimoroni Scroll pHat

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

## License

MIT

