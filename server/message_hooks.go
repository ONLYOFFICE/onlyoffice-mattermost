package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	var newPost *model.Post = post.Clone()

	return newPost, ""
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	// if !strings.Contains(post.Message, "help") {
	// 	return
	// }

	// p.API.SendEphemeralPost(post.UserId, &model.Post{
	// 	ChannelId: post.ChannelId,
	// 	Message:   "You asked for help? Checkout https://about.mattermost.com/help/",
	// })
}
