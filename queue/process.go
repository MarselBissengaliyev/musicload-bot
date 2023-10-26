package queue

func (q *DownloadQueue) startProccess(results chan *Result) {
	for {
		msg := <-q.queue
		go q.downloadAndSend(msg, results)
	}
}

func (q *DownloadQueue) downloadAndSend(m *message, results chan *Result) {
	q.doneWg.Add(1)
	defer q.doneWg.Done()

	result, err := q.handler(q.maxProcessTime, m.url)
	results <- &Result{
		ChatID:   m.chatID,
		Filename: result,
		Err:      err,
	}
}
