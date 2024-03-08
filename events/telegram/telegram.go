package telegram

import (
	"context"
	"errors"

	"fmt"
	"os"
	"encoding/binary"

	"taskbot1/clients/telegram"
	"taskbot1/events"
	"taskbot1/lib/e"
	"taskbot1/storage"
)



type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	p := &Processor{
		tg:      client,
		storage: storage,
	}

	// Set the offset by fetching it from the file
	offset, err := p.fetchOffset()
	if err != nil {
		fmt.Println("Error fetching offset:", err)
		offset = 0 // Set a default offset value if needed
	}

	p.offset = offset

	return p
}

func (p *Processor) Fetch(ctx context.Context, limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(ctx, p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
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

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func (p *Processor) saveOffset() error {
	file, err := os.Create("offset")
	if err != nil {
		return fmt.Errorf("Error while saving offset to file: %v", err)
	}
	defer file.Close()

	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, int64(p.offset))

	_, err = file.Write(buf)
	if err != nil {
		return fmt.Errorf("Error while writing offset to file: %v", err)
	}

	return nil
}

func (p *Processor) fetchOffset() (int, error) {
	file, err := os.Open("offset")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Read bytes from the file
	buf := make([]byte, binary.MaxVarintLen64)
	_, err = file.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("Error while reading the offset file: %v", err)
	}

	// Parse bytes into integer
	offset, _ := binary.Varint(buf)
	if err != nil {
		return 0, fmt.Errorf("Error while parsing buffer: %v", err)
	}

	return int(offset), nil
}
