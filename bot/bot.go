package bot

import (
	"fmt"
	"os"

	"github.com/MarselBissengaliyev/realtime-chat/downloader/youtube"
	"github.com/MarselBissengaliyev/realtime-chat/queue"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI

	downloadService *youtube.Downloader
	queue           *queue.DownloadQueue
	downloadMsgs    chan *queue.Result

	username    string
	maxDuration int64
}

func NewTelegramBot(token string, maxDownloadTime, maxVideoDuration int64, username string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	downloadService, err := youtube.NewDownloader(maxVideoDuration)
	if err != nil {
		return nil, err
	}

	downloadQueue := queue.NewDownloadQueue(downloadService.Download, maxVideoDuration)

	return &TelegramBot{
		bot:             bot,
		downloadService: downloadService,
		queue:           downloadQueue,
		downloadMsgs:    make(chan *queue.Result),
		maxDuration:     maxVideoDuration,
		username:        username,
	}, nil
}

func (t *TelegramBot) Run(debug bool) error {
	fmt.Println("running")

	t.bot.Debug = debug

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	t.queue.Start(t.downloadMsgs)
	go t.mailoutDownloads()

	updates := t.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		fmt.Println(update)
		t.HandleUpdates(update)
	}

	return nil
}

func (t *TelegramBot) send(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	t.bot.Send(msg)
}

func (t *TelegramBot) sendError(chatID int64) {
	t.send(chatID, "Error occured. Try again.")
}

func (t *TelegramBot) sendAudioFile(chatID int64, filename string) {
	path := "./" + filename
	defer os.Remove(path)

	audioCfg := tgbotapi.NewAudio(chatID, tgbotapi.FilePath(path))
	audioCfg.Caption = "Downloaded via @" + t.username

	_, err := t.bot.Send(audioCfg)
	if err != nil {
		fmt.Printf("error sending message: %s", err.Error())
		t.sendError(chatID)
	}
}

func (t *TelegramBot) Stop() {
	t.queue.Stop()
	close(t.downloadMsgs)
}
