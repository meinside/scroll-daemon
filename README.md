# Go daemon for Pimoroni Scroll pHat

## Install

```bash
$ go get -u github.com/meinside/scroll-daemon/...
```

## Configuration

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

