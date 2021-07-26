package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
)

type ONLYOFFICE_BOT struct {
	Id           string
	LoggerPrefix string
	P            *Plugin
}

func (bot *ONLYOFFICE_BOT) BOT_CREATE_POST(message string, channelId string) {
	ONLYOFFICE_BOT_POST := model.Post{
		Message:   message,
		ChannelId: channelId,
		UserId:    bot.Id,
	}

	_, creationErr := bot.P.API.CreatePost(&ONLYOFFICE_BOT_POST)
	if creationErr != nil {
		bot.P.API.LogError(ONLYOFFICE_BOT_LOGGER_PREFIX + "Post creation error")
		return
	}

	bot.P.API.LogInfo(bot.LoggerPrefix + "Created a new post")
}

func (bot *ONLYOFFICE_BOT) BOT_CREATE_REPLY(message string, channelId string, parentId string) {
	ONLYOFFICE_BOT_POST := model.Post{
		Message:   message,
		ParentId:  parentId,
		RootId:    parentId,
		ChannelId: channelId,
		UserId:    bot.Id,
	}

	_, creationErr := bot.P.API.CreatePost(&ONLYOFFICE_BOT_POST)
	if creationErr != nil {
		bot.P.API.LogError(ONLYOFFICE_BOT_LOGGER_PREFIX + "Post creation error")
		return
	}
	bot.P.API.LogInfo(bot.LoggerPrefix + "Created a new reply")
}
