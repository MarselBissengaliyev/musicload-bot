package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (t *TelegramBot) HandleUpdates(update tgbotapi.Update) {
	if m := update.Message; m != nil {
		if m.IsCommand() && m.Command() == "start" {
			t.send(m.Chat.ID, "Hi there! Send me a link to video you want extract music from.")
			return
		}

		if t.downloadService.IsValidURL(m.Text) {
			t.send(m.Chat.ID, "Wait for the link to process..")
			t.queue.Enqueue(m)
			return
		}

		t.send(m.Chat.ID, "Invalid message text. I'm waiting for youtube link")
	}
}
