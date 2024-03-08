package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	//"net/url"
	"strings"

	"taskbot1/lib/e"
	"taskbot1/storage"
)

var isLogin bool

const (
	StartCmd    = "/start"
	HelpCmd     = "/help"
	CommandsCmd = "/commands"
	TasksCmd    = "/tasks"
	AddCmd      = "/add"
	RemoveCmd   = "/remove"
	CompleteCmd = "/complete"
	DeadlineCmd = "/deadline"
	authCmd     = "/auth"
	registerCmd = "/register"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	command, content, _ := fetchInput(text)
	if command == "" && content == "" {
		command = text
	}

	switch command {
	// case authCmd:
	// 	return p.Auth(chatID, username)
	// case registerCmd:
	// 	return p.Register(chatID, username)
	case StartCmd:
		return p.tg.SendMessage(ctx, chatID, msgHello)
	case HelpCmd:
		return p.tg.SendMessage(ctx, chatID, msgHelp)
	case CommandsCmd:
		return p.tg.SendMessage(ctx, chatID, msgCommands)
	case AddCmd:

		return p.saveTask(ctx, chatID, content, username)

	case RemoveCmd:

		return p.removeTask(ctx, chatID, content, username)

	case CompleteCmd:

		return p.completeTask(ctx, chatID, content, username)

	case DeadlineCmd:

		return p.deadlineTask(ctx, chatID, content, username)

	case TasksCmd:

		return p.showTasks(ctx, chatID, username)

	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand+": "+command)
	}
}

// func (p *Processor) Register(chatID int, username string) (err error) {
// 	pretendent := &storage.User{
// 		Username: username,
// 		Chatid:   chatID,
// 	}
// 	_, err = p.storage.RetrieveUser(pretendent)
// 	if err != nil {
// 		errInSaving := p.storage.SaveUser(pretendent)
// 		if errInSaving != nil {
// 			return p.tg.SendMessage(context.Background(), chatID, msgUserExist)
// 		}
// 	}
// 	return p.tg.SendMessage(context.Background(), chatID, "Send /auth to start!")
// }
// func (p *Processor) Auth(chatID int, username string) (err error) {
// 	//check
// 	pretendent := &storage.User{
// 		Username: username,
// 		Chatid:   chatID,
// 	}
// 	_, err = p.storage.RetrieveUser(pretendent)
// 	if err != nil {
// 		return p.tg.SendMessage(context.Background(), chatID, msgPlsRegister)

// 	}
// 	isLogin = true

// 	return p.tg.SendMessage(context.Background(), chatID, msgHello)
// }

func (p *Processor) showTasks(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send tasks", err) }()

	tasks, err := p.storage.Tasks(ctx, chatID)
	if err != nil && !errors.Is(err, storage.ErrNoSavedTasks) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedTasks) {
		return p.tg.SendMessage(ctx, chatID, msgNoSavedTasks)
	}

	taskListText := "Task List:\n"
	for _, task := range tasks {
		deadline := ""
		completedStatus := "Not Completed"
		if task.Completed == 1 {
			completedStatus = "Completed"
		}
		if !task.Deadline.IsZero() {
			deadline = fmt.Sprintf("(Deadline : %s)", task.Deadline.Format("2006-06-02 15:04:05"))
		}
		taskListText += fmt.Sprintf("- Task %d: %s \n(Created: %s, %s)\n %s", task.ID, task.Content, task.Created.Format("2006-01-02 15:04:05"), completedStatus, deadline)
	}

	return p.tg.SendMessage(ctx, chatID, taskListText)
}

func (p *Processor) completeTask(ctx context.Context, chatID int, idString string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: complete task:", err) }()

	_, err = strconv.Atoi(idString)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, msgIncorrectInput)
	}
	id, _ := strconv.Atoi(idString)

	isExists, err := p.storage.IsExistsID(ctx, id)
	if err != nil {
		return err
	}

	if !isExists {
		return p.tg.SendMessage(ctx, chatID, msgDoesntExists)
	}
	err = p.storage.Complete(ctx, id)

	if err == nil {
		return p.tg.SendMessage(ctx, chatID, msgCompleted)
	} else {
		return err
	}
}

func (p *Processor) deadlineTask(ctx context.Context, chatID int, text string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: set deadline to task:", err) }()
	textSlice := strings.Split(text, " ")
	var id, days int
	if len(textSlice) == 2 {
		_, err1 := strconv.Atoi(textSlice[0])
		_, err2 := strconv.Atoi(textSlice[1])

		if err1 == nil && err2 == nil {
			id, _ = strconv.Atoi(textSlice[0])
			days, _ = strconv.Atoi(textSlice[1])
		} else {
			return p.tg.SendMessage(ctx, chatID, msgIncorrectInput+"\n Should be /deadline \"<id> <id>\"")
		}
	} else {
		return p.tg.SendMessage(ctx, chatID, msgIncorrectInput+"\n Should be /deadline \"<id> <id>\"")
	}

	isExists, err := p.storage.IsExistsID(ctx, id)
	if err != nil {
		return err
	}

	if !isExists {
		return p.tg.SendMessage(ctx, chatID, msgDoesntExists)
	}

	err = p.storage.Deadline(ctx, id, days)

	if err == nil {
		return p.tg.SendMessage(ctx, chatID, msgDeadline)
	} else {
		return fmt.Errorf("Error while setting deadline: ", err)
	}
}

//func (p *Processor) pastTasks(chatID int, username string) (err error) {
//	defer func() { err = e.WrapIfErr("can't do command: can't send tasks", err) }()
//
//	tasks, err := p.storage.PastTasks()
//	if err != nil && !errors.Is(err, storage.ErrNoPastTasks {
//		return err
//	}
//
//	if errors.Is(err, storage.ErrNoSavedTasks) {
//		return p.tg.SendMessage(chatID, msgNoPastTasks)
//	}
//
//	r hasInpu
//}

func (p *Processor) saveTask(ctx context.Context, chatID int, content string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save task", err) }()

	isExists, err := p.storage.IsExists(ctx, content)
	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(ctx, chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, content, chatID); err != nil {
		return err
	}

	return p.tg.SendMessage(ctx, chatID, msgSaved)
}

func (p *Processor) removeTask(ctx context.Context, chatID int, idString string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: remove task", err) }()

	_, err = strconv.Atoi(idString)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, msgIncorrectInput)
	}
	id, _ := strconv.Atoi(idString)

	isExists, err := p.storage.IsExistsID(ctx, id)
	if err != nil {
		return err
	}

	if !isExists {
		return p.tg.SendMessage(ctx, chatID, msgDoesntExists)
	}

	if err := p.storage.Remove(ctx, id); err != nil {
		return err
	}

	return p.tg.SendMessage(ctx, chatID, msgRemoved)
}

func fetchInput(text string) (string, string, error) {
	pattern := `^/(\w+)\s+"([^"]+)"$`

	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(text)

	// Check if there is a match
	if len(matches) != 3 {
		return "", "", fmt.Errorf("Invalid command format: %s", text)
	}

	command := "/" + matches[1]
	content := matches[2]

	return command, content, nil
}
