package bot

import (
	"fmt"

	"github.com/MarselBissengaliyev/realtime-chat/downloader"
)

func (t *TelegramBot) mailoutDownloads() {
	for {
		res := <-t.downloadMsgs
		if res.Err != nil {
			fmt.Printf(res.Err.Error())
			if res.Err == downloader.ErrDurationLimitExceeded {
				t.send(res.ChatID, fmt.Sprintf("Can't download video longer than %d minutes", t.maxDuration))
				continue
			}

			fmt.Printf("mailoutDownloads() error: %s\n", res.Err.Error())

			go t.sendError(res.ChatID)
			continue
		}

		fmt.Printf("send result: %+v\n", res)

		go t.sendAudioFile(res.ChatID, res.Filename)
	}
}
