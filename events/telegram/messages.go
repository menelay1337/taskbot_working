package telegram

const msgHelp = `I can save and keep your tasks. In any time you can fetch past and present tasks, mark as completed and set deadlines.

/add "content" command to save task.
/tasks command, command to see your tasks.
/commands command to obtain available options.
/help command to get help with the bot.
`

const msgHello = "Hi there \n\n" + msgHelp

const msgCommands = `
/start starts the bot
/help obtain the help with the bot
/commands obtain the commands
/tasks outputs all available tasks
/add "<Content>" add task with some <Content>
/remove "<id>" removes specific task with some <id>
/complete "<id>" marks task with <id> as completed.
`

const (
	msgNoSavedTasks = "You have no saved task"
	msgNoPastTasks = "There are no past tasks"
	msgSaved = "Task is saved."
	msgRemoved = "Task is removed."
	msgCompleted = "Task is completed."
	
	msgDoesntExists = "Task doesn't exist."
	msgUnknownCommand = "Unknown command"
	msgAlreadyExists = "You have already have this task in your list"
	msgIncorrectInput = "Incorrept input"
)
