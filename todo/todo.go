// Package todo is a plugin for bawt that creates to do lists per channel
package todo

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gopherworks/bawt"
)

func (p *Plugin) listenTodo() {
	p.bot.Listen(&bawt.Listener{
		Matches:            regexp.MustCompile(`^!todo.*`),
		MessageHandlerFunc: p.handleTodo,
		Name:               "To Do",
		Description:        "Keeps a tab of all your to do's!",
		Commands: []bawt.Command{
			{
				Usage:    "!todo",
				HelpText: "Displays a list of tasks",
			},
			{
				Usage:    "!todo add <some text>",
				HelpText: "Displays a list of tasks",
			},
			{
				Usage:    "!todo scratch <id>",
				HelpText: "Removes a task",
			},
			{
				Usage:    "!todo append <id> <some text>",
				HelpText: "Adds to the end of a task",
			},
		},
	})
}

func (p *Plugin) handleTodo(listen *bawt.Listener, msg *bawt.Message) {
	idFormat := regexp.MustCompile(`^[a-z]{2}$`)
	parts := strings.Split(msg.Match[0], " ")
	if len(parts) == 1 {
		p.listTasks(msg)
		return
	}
	act := parts[1]

	switch act {
	case "add":
		if len(parts) < 2 {
			msg.ReplyMention("Add a task with `!todo add [some text]`")
			return
		}
		p.createTask(msg, strings.Join(parts[2:], " "))

	case "scratch":
		if len(parts) < 3 || !idFormat.MatchString(parts[2]) {
			msg.ReplyMention(fmt.Sprintf("Please %s a task with `!todo %s ID`", act, act))
			return
		}

		p.deleteTask(msg, parts[2], false)

	case "append":
		if len(parts) < 4 || !idFormat.MatchString(parts[2]) {
			msg.ReplyMention(fmt.Sprintf("Please %s a task with `!todo %s ID [more notes]`", act, act))
			return
		}

		p.appendToTask(msg, parts[2], strings.Join(parts[3:], " "))

	case "help":
		p.replyHelp(msg, "")

	default:
		if idFormat.MatchString(act) {
			p.replyHelp(msg, "Wooops, not sure what you wanted.\n")
		} else {
			p.listTasks(msg)
		}
	}
}

func (p *Plugin) detailTask(msg *bawt.Message, id string) {
	todo := p.store.Get(msg.Channel)
	index, err := getTaskIndex(id, todo)
	if err != nil {
		msg.ReplyMention("Task not found...")
		return
	}
	task := todo[index]
	msg.Reply(printTaskDetails(task))
}

func printTaskDetails(task *Task) string {
	return fmt.Sprintf("%s\n> Created %s by <@%s>", task.String(), task.CreatedAt.Format("2006-01-02 15:04:05"), task.CreatedBy)
}

func (p *Plugin) createTask(msg *bawt.Message, content string) {
	todo := p.store.Get(msg.Channel)

	if len(todo) > 600 {
		msg.ReplyMention("Gosh you have over 600 tasks!!! Clean some up first.")
		return
	}

	id := p.generateRandomID(todo)
	task := &Task{
		ID:        id,
		CreatedAt: time.Now(),
		CreatedBy: msg.FromUser.ID,
		Text:      []string{content},
	}
	todo = append(todo, task)
	p.store.Put(msg.Channel, todo)
	msg.ReplyMention("added: " + task.String())
}

func (p *Plugin) appendToTask(msg *bawt.Message, id, text string) {
	todo := p.store.Get(msg.Channel)
	index, err := getTaskIndex(id, todo)
	if err != nil {
		msg.ReplyMention("Task not found...")
		return
	}

	task := todo[index]
	task.Text = append(task.Text, strings.Split(text, " // ")...)
	p.store.Put(msg.Channel, todo)

	msg.ReplyMention("updated " + task.String())
}

func (p *Plugin) listTasks(msg *bawt.Message) {
	todo := p.store.Get(msg.Channel)
	sort.Sort(byID(todo))

	var answer []string
	var toDelete []string
	for _, task := range todo {
		if task.Closed {
			toDelete = append(toDelete, task.ID)
		} else {
			answer = append(answer, task.String())
		}
	}
	if len(toDelete) != 0 {
		p.deleteTask(msg, strings.Join(toDelete, ","), true)
	}
	if len(answer) == 0 {
		msg.ReplyMention("Nothing to do... Coffee time?")
	} else {
		msg.Reply(strings.Join(answer, "\n"))
	}
}

func (p *Plugin) deleteTask(msg *bawt.Message, ids string, silent bool) {
	todo := p.store.Get(msg.Channel)

	parts := strings.Split(msg.Match[0], " ")
	var closingNotes string
	if len(parts) > 3 {
		closingNotes = strings.Join(parts[3:], " ")
	}

	var out []string
	for _, id := range strings.Split(ids, ",") {
		index, err := getTaskIndex(id, todo)
		if err != nil {
			out = append(out, "Task `"+id+"` not found")
			continue
		}

		task := todo[index]
		task.Closed = true
		task.ClosingNote = closingNotes
		todo = append(todo[:index], todo[index+1:]...)

		if silent != true {
			out = append(out, task.String())
		}
	}

	p.store.Put(msg.Channel, todo)

	msg.Reply(strings.Join(out, "\n"))
}

func getTaskIndex(id string, todo Todo) (int, error) {
	for i, task := range todo {
		if task.ID == id {
			return i, nil
		}
	}
	return 0, errors.New("Not found")
}

func (p *Plugin) replyHelp(msg *bawt.Message, extra string) {
	answer := extra + `Commands:` + "```" + `
!todo add [some text]             - add task
!todo                             - list tasks
!todo scratch [id]                - deletes task(s)
!todo append [id] [more stuff]    - append text to a task
!todo help                        - show this help
` + "```"
	msg.Reply(answer)
	return
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (p *Plugin) generateRandomID(todo Todo) string {
	for {
		id := randSeq(2)
		if idInList(id, todo) {
			continue
		}
		return id
	}
}

func idInList(id string, todo Todo) bool {
	for _, task := range todo {
		if task.ID == id {
			return true
		}
	}
	return false
}
