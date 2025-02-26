package main

import (
	"context"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

func getMCGUsers(ctx context.Context, api *slack.Client) ([]slack.User, error) {
	users, err := api.GetUsersContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "GetUsersContext")
	}
	res := make([]slack.User, 0)
	for _, user := range users {
		mcg := user.IsRestricted && !user.IsUltraRestricted && !user.IsBot
		if mcg {
			res = append(res, user)
		}
	}
	return res, nil
}

func getPrefixChannel(ctx context.Context, api *slack.Client, prefix string) ([]slack.Channel, error) {
	cursor := ""
	for i := 0; i < 30; i++ {
		res := make([]slack.Channel, 0)
		channels, new_cursor, err := api.GetConversationsContext(ctx,
			&slack.GetConversationsParameters{
				ExcludeArchived: true,
				Limit:           999,
				Cursor:          cursor,
			})
		cursor = new_cursor
		if err != nil {
			return nil, errors.Wrapf(err, "GetChannels")
		}
		for _, channel := range channels {
			if !channel.IsChannel {
				continue
			}
			if strings.HasPrefix(channel.Name, prefix) {
				res = append(res, channel)
			}
		}
		if cursor == "" {
			return res, nil
		}
	}
	return nil, errors.New("Too many channels")
}

func invite(ctx context.Context, api *slack.Client, channelID string, userID []string) error {
	if _, _, _, err := api.JoinConversationContext(ctx, channelID); err != nil {
		return errors.Wrapf(err, "JoinConversation")
	}
	if _, err := api.InviteUsersToConversationContext(ctx, channelID, userID...); err != nil {
		log.Printf("InviteUsersToConversation is failed: %v (%s <- %v)", err, channelID, userID)
	}
	return nil
}

func inviteAllMCG(ctx context.Context, api *slack.Client, prefix string) ([]slack.User, []slack.Channel, error) {
	users, err := getMCGUsers(ctx, api)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "getMCGUsers")
	}
	userIDs := []string{}
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	channels, err := getPrefixChannel(ctx, api, prefix)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "getPrefixChannels")
	}
	for _, channel := range channels {
		invite(ctx, api, channel.ID, userIDs)
	}
	return users, channels, nil
}

func postMessage(api *slack.Client, ctx context.Context, channel string, msg string) error {
	if _, _, _, err := api.JoinConversationContext(ctx, channel); err != nil {
		return errors.Wrapf(err, "JoinConversation")
	}
	_, _, err := api.PostMessageContext(ctx, string(channel), slack.MsgOptionText(msg, false))
	if err != nil {
		return errors.Wrapf(err, "PostMessage")
	}
	return nil
}
