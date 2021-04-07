package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type Job struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
}

func doWork(job *Job, db *gorm.DB) {
	if job.Session.State.User.ID == job.Message.Author.ID {
		return
	}

	myUser := "@" + job.Session.State.User.Username
	cmds := []struct {
		verb string
		call func(*Context)
	}{
		{"announce", SetAnnounce},
		{"karma", CheckKarma},
	}

	msg := job.Message.ContentWithMentionsReplaced()
	for _, c := range cmds {
		header := fmt.Sprintf("%s %s ", myUser, c.verb)
		if strings.HasPrefix(msg, header) {
			c.call(&Context{job, db, header})
			return
		}
	}

	ModKarma(&Context{job, db, ""})
}

func worker(workQueue <-chan Job, cancel <-chan struct{}, db *gorm.DB) {
	for {
		select {
		case <-cancel:
			return
		case job := <-workQueue:
			doWork(&job, db)
		}
	}
}
