package telegram

import (
	"context"
	"taskbot1/storage"

	//"errors"
	"log"
	//"net/url"
	"strings"
	//"taskbot1/lib/e"
	//"taskbot1/storage"
)

var isLogin bool

const (
	RndCmd      = "/rnd"
	HelpCmd     = "/help"
	StartCmd    = "/start"
	authCmd     = "/auth"
	registerCmd = "/register"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)
	if text == registerCmd {
		err := p.Register(chatID, username)
		if err != nil {
			return p.tg.SendMessage(ctx, chatID, "got message: "+err.Error())
		}

	}
	if text == authCmd {
		err := p.Auth(chatID, username)
		if err != nil {
			return p.tg.SendMessage(ctx, chatID, "got message: "+err.Error())
		}
		isLogin = true
		return p.tg.SendMessage(ctx, chatID, "got message: "+"Succes!")
	}
	return nil
	// return p.tg.SendMessage(ctx, chatID, "got message: "+"Succes!")
	// return p.tg.SendMessage(ctx, chatID, "got message: "+text)

	//if isAddCmd(text) {
	//	return p.savePage(ctx, chatID, text, username)
	//}

	//switch text {
	//case RndCmd:
	//	return p.sendRandom(ctx, chatID, username)
	//case HelpCmd:
	//	return p.sendHelp(ctx, chatID)
	//case StartCmd:
	//	return p.sendHello(ctx, chatID)
	//default:
	//	return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	//}
}
func (p *Processor) Register(chatid int, username string) (err error) {
	pretendent := &storage.User{
		Username: username,
		Chatid:   chatid,
	}
	_, err = p.storage.RetrieveUser(pretendent)
	if err != nil {
		errInSaving := p.storage.SaveUser(pretendent)
		if errInSaving != nil {
			return p.tg.SendMessage(context.Background(), chatid, msgUserExist)
		}
	}
	return p.tg.SendMessage(context.Background(), chatid, msgHello)
}
func (p *Processor) Auth(chatid int, username string) (err error) {
	//check
	pretendent := &storage.User{
		Username: username,
		Chatid:   chatid,
	}
	_, err = p.storage.RetrieveUser(pretendent)
	if err != nil {
		return p.tg.SendMessage(context.Background(), chatid, msgPlsRegister)

	}
	isLogin = true

	return p.tg.SendMessage(context.Background(), chatid, msgHello)
}

//func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string) (err error) {
//	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()
//
//	page := &storage.Page{
//		URL:      pageURL,
//		UserName: username,
//	}
//
//	isExists, err := p.storage.IsExists(ctx, page)
//	if err != nil {
//		return err
//	}
//	if isExists {
//		return p.tg.SendMessage(ctx, chatID, msgAlreadyExists)
//	}
//
//	if err := p.storage.Save(ctx, page); err != nil {
//		return err
//	}
//
//	if err := p.tg.SendMessage(ctx, chatID, msgSaved); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) {
//	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()
//
//	page, err := p.storage.PickRandom(ctx, username)
//	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
//		return err
//	}
//	if errors.Is(err, storage.ErrNoSavedPages) {
//		return p.tg.SendMessage(ctx, chatID, msgNoSavedPages)
//	}
//
//	if err := p.tg.SendMessage(ctx, chatID, page.URL); err != nil {
//		return err
//	}
//
//	return p.storage.Remove(ctx, page)
//}
//
//func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
//	return p.tg.SendMessage(ctx, chatID, msgHelp)
//}
//
//func (p *Processor) sendHello(ctx context.Context, chatID int) error {
//	return p.tg.SendMessage(ctx, chatID, msgHello)
//}
//
//func isAddCmd(text string) bool {
//	return isURL(text)
//}
//
//func isURL(text string) bool {
//	u, err := url.Parse(text)
//
//	return err == nil && u.Host != ""
//}
