package invbot

import (
	"context"
	"log"
	"strings"

	"github.com/slack-go/slack"
)

func getMCGUsers(ctx context.Context, api *slack.Client) []slack.User {
	users, err := api.GetUsersContext(ctx)
	if err != nil {
		log.Fatalf("GetUsers was failed: %v", err)
	}
	res := make([]slack.User, 0)
	for _, user := range users {
		mcg := user.IsRestricted && !user.IsUltraRestricted && !user.IsBot
		if mcg {
			res = append(res, user)
		}
	}
	return res
}

func getPrefixChannel(ctx context.Context, api *slack.Client, prefix string) []slack.Channel {
	cursor := ""
	for i := 0; i < 30; i++ {
		res := make([]slack.Channel, 0)
		channels, new_cursor, err := api.GetConversationsContext(ctx,
			&slack.GetConversationsParameters{
				ExcludeArchived: "true",
				Limit: 999,
				Cursor: cursor,
			})
		cursor = new_cursor
		if err != nil {
			log.Fatalf("GetChannels is failed: %v", err)
		}
		for _, channel := range channels {
			if ! channel.IsChannel {
				continue
			}
			if strings.HasPrefix(channel.Name, prefix) {
				res = append(res, channel)
			}
		}
		if cursor == "" {
			return res
		}
	}
	panic("Too many Channels")
}

func invite(ctx context.Context, api *slack.Client, channelID string, userID []string) {
	if _, _, _, err := api.JoinConversationContext(ctx, channelID); err != nil {
		log.Printf("JoinConversation is failed: %v", err)
	}
	if _, err := api.InviteUsersToConversationContext(ctx, channelID, userID...); err != nil {
		log.Printf("InviteUsersToConversation is failed: %v (%s <- %v)", err, channelID, userID)
	}
}

func inviteAllMCG(ctx context.Context, api *slack.Client, prefix string) {
	users := getMCGUsers(ctx, api)
	userIDs := make([]string, 0)
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	channels := getPrefixChannel(ctx, api, prefix)
	for _, channel := range channels {
		invite(ctx, api, channel.ID, userIDs)
	}
}