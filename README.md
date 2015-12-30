# ANTBOT

A telegram bot create time-lapse video and publish to channel. It uses `ffmpeg` to do the record and encode job. You can find out the command line arguments passed to `ffmpeg` in [grab.go](https://github.com/Patrolavia/antbot/blob/master/grab.go).

## Synopsis

You'll need a working `ffmpeg` in your path, and ensure it can record video from your cam.

#### Publish to Telegram channel

You'll need a Telegram bot token, refer to Telgram official site to find out how you get one. Store the token in a file, say `token` in current directory.

```sh
# Generally, we use v4l to grab cam in linux
antbot -ch mychannelname -f v4l2 -i /dev/your_cam -seg 3600 -size 640x480 -spf 1 -t token
```

#### Publish via scp

You'll need a working ssh client in path. `Antbot` will not ask password or accept unknown host, so you have to authorize via key and add correct entry in your `known_hosts` file.

```sh
antbot -f v4l2 -i /dev/your_cam -seg 3600 -size 640x480 -spf 1 -scp you@example.com:path/to/web/ant.mp4
```

See `antbot -h` for detailed usage.

## License

Any version of GPL, LGPL or AGPL.
