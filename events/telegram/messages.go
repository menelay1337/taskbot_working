package telegram

const msgHelp = `I can save and keep your tasks. In any time you can fetch past and present tasks, mark as completed and set deadlines.

/add "content" command to save task.
/tasks command, command to see your tasks.
/commands command to obtain available options.
/help command to get help with the bot.
If you are the first time here - send /register
`

const msgHello = "Hi, that's Task bot \n\n" + msgHelp

const msgCommands = `
/start starts the bot
/help obtain the help with the bot
/commands obtain the commands
/tasks outputs all available tasks
/add "<Content>" add task with some <Content>
/remove "<id>" removes specific task with some <id>
/complete "<id>" marks task with <id> as completed.
/deadline "<id> <days>" sets deadline to task with <id> for <days> days.

`

const (
	msgSaved = "Task is saved."
	msgRemoved = "Task is removed."
	msgCompleted = "Task is completed."
	msgDeadline = "Deadline set to task."
	
	msgNoPastTasks = "There are no past tasks"
	msgNoSavedTasks = "You have no saved task"
	msgDoesntExists = "Task doesn't exist."
	msgUnknownCommand = "Unknown command"
	msgAlreadyExists  = "You have already have this task in your list"
	msgIncorrectInput = "Incorrept input"
	msgPlsRegister    = "Please, register! /register"
	msgUserExist      = "User already exists! If it's you, send /auth"
)
