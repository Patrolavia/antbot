# ANTBOT

自動錄製縮時影片並上傳到 Telegram 頻道的 bot。

使用 `Antbot` 可能需要對 `ffmpeg` 有一定的認識，因為`Antbot` 會使用 `ffmpeg` 來錄製和轉檔。你可以在 [grab.go](https://github.com/Patrolavia/antbot/blob/master/grab.go) 裡找到相關的參數。

## 概論

你的電腦裡必須先安裝、設定好 `ffmpeg`，確定它可以從你的視訊攝影機取得影像。

#### 發佈到 Telegram 頻道

你需要一個 Telegram bot 的 token 才能使用 Telegram 相關的功能。把 token 存到一個文字檔裡，假設叫 `token`。

```sh
# 在 linux 裡通常我們是用 v4l 來取得視訊攝影機的影像
antbot -ch mychannelname -f v4l2 -i /dev/your_cam -seg 3600 -size 640x480 -spf 1 -t token
```

#### 用 scp 發佈到伺服器上

你得先安裝和設定好 ssh 程式。由於 `Antbot` 不會跟你要密碼，所以你也得設定好 ssh 讓它可以用金鑰登入。

```sh
antbot -f v4l2 -i /dev/your_cam -seg 3600 -size 640x480 -spf 1 -scp you@example.com:path/to/web/ant.mp4
```

`antbot -h` 可以查看其他的設定。

## License

Any version of GPL, LGPL or AGPL.
