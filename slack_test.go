package invbot

import (
	"context"
	"testing"

	"github.com/slack-go/slack"
)

func TestUsers(t *testing.T) {
	config := setup()
	api := slack.New(config.Token, slack.OptionDebug(true))
	ctx := context.Background()
	users := getMCGUsers(ctx, api)
	for _, user := range users {
		t.Logf("%s: %s", user.Name, user.ID)
	}
}

func TestChannels(t *testing.T) {
	config := setup()
	api := slack.New(config.Token, slack.OptionDebug(true))
	ctx := context.Background()
	channels := getPrefixChannel(ctx, api, config.Prefix)
	for _, channel := range channels {
		t.Logf("%s: %s", channel.Name, channel.ID)
	}
}

func TestInvite(t *testing.T) {
	config := setup()
	api := slack.New(config.Token, slack.OptionDebug(true))
	ctx := context.Background()
	combatUser := "U011S918AP4"
	randomChannel := "C011Z205A8L"
	invite(ctx, api, randomChannel, []string{combatUser})
}

func TestInviteAllMCG(t *testing.T) {
	config := setup()
	api := slack.New(config.Token, slack.OptionDebug(true))
	ctx := context.Background()
	inviteAllMCG(ctx, api, config.Prefix)
}
