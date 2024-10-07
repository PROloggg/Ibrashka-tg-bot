package telegram

import (
    "errors"

    "app/clients/telegram"
    "app/events"
    "app/lib/e"
    "app/storage"
)

type Processor struct {
    tg *telegram.Client
    offset int
    storage storage.Storage
}

type Meta struct {
    ChatId int
    UserName string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New (client *telegram.Client, storage storage.Storage) *Processor {
    return &Processor {
        tg: client,
        storage: storage,
    }
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
    updates,err := p.tg.Updates(p.offset, limit)
    if err != nil {
        return nil,e.Wrap("can't get events", err)
    }

    if len(updates) == 0 {
        return nil, nil
    }

    res := make([]events.Event, 0, len(updates))

    for _, updateItem := range updates {
        res = append(res, event(updateItem))
    }

    p.offset = updates[len(updates) - 1].UpdateID + 1

    return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
        case events.Message:
            return p.processMessage(event)
        default:
            return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatId, meta.UserName, event.Photos); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(updateItem telegram.Update) events.Event {
    typeUpdate := fetchType(updateItem)

    res := events.Event{
        Type: typeUpdate,
        Text: fetchText(updateItem),
        Photos: updateItem.Message.Photo,
    }

    if typeUpdate == events.Message {
        res.Meta = Meta {
            ChatId: updateItem.Message.Chat.ID,
            UserName: updateItem.Message.From.Username,
        }
        res.Photos = updateItem.Message.Photo
    }

    return res
}

func fetchType(updateItem telegram.Update) events.Type {
    if updateItem.Message == nil {
        return events.Unknown
    }

    return events.Message
}

func fetchText(updateItem telegram.Update) string {
    if updateItem.Message == nil {
        return ""
    }

    if updateItem.Message.Text != "" {
        return updateItem.Message.Text
    }

    return updateItem.Message.Caption
}