package queue

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandleFunc func(maxVideoDuration int64, videoId string) (string, error)

type message struct {
	url    string
	chatID int64
}

type Result struct {
	ChatID   int64
	Filename string
	Err      error
}

type DownloadQueue struct {
	queue  chan *message
	doneWg *sync.WaitGroup

	maxProcessTime int64
	handler        HandleFunc
}

func NewDownloadQueue(downloadService HandleFunc, maxProcessTime int64) *DownloadQueue {
	return &DownloadQueue{
		queue:          make(chan *message),
		doneWg:         new(sync.WaitGroup),
		maxProcessTime: maxProcessTime,
		handler:        downloadService,
	}
}

func (q *DownloadQueue) Start(results chan *Result) {
	go q.startProccess(results)
}

func (q *DownloadQueue) Enqueue(m *tgbotapi.Message) {
	msg := q.toMessage(m)
	q.queue <- msg
}

func (q *DownloadQueue) toMessage(m *tgbotapi.Message) *message {
	return &message{
		chatID: m.Chat.ID,
		url:    m.Text,
	}
}

func (q *DownloadQueue) Stop() {
	q.doneWg.Wait()
	close(q.queue)
}