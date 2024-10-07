package telegram

import (
	"app/parser"
	"log"
	"os"
	"strings"
	"time"

	"app/clients/telegram"
	"app/config"
	"app/storage"
)

const HelpCmd = "/help"

func (p *Processor) doCmd(text string, chatId int, userName string, photos []telegram.Photo) error {
	text = strings.TrimSpace(text)

	// Обработка текстовых команд
	switch text {
	case HelpCmd:
		return p.sendHelp(chatId)
	default:
	}

	if isMessageForMe(text) == false {
		return nil
	}

	log.Printf("got new command from '%s'", userName)
	var pictureUrl string
	var pictureErr error

	// Проверка наличия фото
	if len(photos) > 0 {
		pictureUrl, pictureErr = p.tg.GetFileURL(photos[len(photos)-1].FileID)
		if pictureErr != nil {
			log.Printf("Error handling photos: %v", pictureErr)
		}
	}

	return p.savePage(chatId, text, userName, pictureUrl)
}

func isMessageForMe(text string) bool {
	return strings.Contains(text, "#save")
}

func (p *Processor) savePage(chatId int, text string, userName string, pictureUrl string) (err error) {
	defer func() {
		if err != nil {
			if message, exists := parser.ValidationMessages[err.Error()]; exists {
				//Если сообщение существует, отправляем его
				if err := p.tg.SendMessage(chatId, msgParserError+message); err != nil {
					return
				}
			}
			if err = p.tg.SendMessage(chatId, msgRecommendedFormat); err != nil {
				return
			}
		}
	}()

	page := &storage.Page{
		Text:       text,
		UserName:   userName,
		CreatedAt:  time.Now(),
		PictureUrl: pictureUrl,
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatId, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatId int) error {
	err := p.tg.SendMessage(chatId, msgHelp)
	if err != nil {
		return err
	}

	voiceFilePath := config.Get().HelpVoicePath
	if _, err := os.Stat(voiceFilePath); os.IsNotExist(err) {
		return nil
	}

	return p.tg.SendVoice(chatId, voiceFilePath)
}
