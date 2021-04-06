package invbot

import (
	"context"
	"testing"

	"github.com/slack-go/slack"
)

func TestUsers(t *testing.T) {
	api := slack.New(config.Token, slack.OptionDebug(true))
	ctx := context.Background()
	users, err := getMCGUsers(ctx, api)
	if err != nil {
		t.Fatal(err)
	}
	for _, user := range users {
		t.Logf("%s: %s", user.Name, user.ID)
	}
}

func TestChannels(t *testing.T) {
	api := slack.New(config.Token, slack.OptionDebug(true))
	ctx := context.Background()
	channels, err := getPrefixChannel(ctx, api, config.Prefix)
	if err != nil {
		t.Fatal(err)
	}
	for _, channel := range channels {
		t.Logf("%s: %s", channel.Name, channel.ID)
	}
}

func TestInvite(t *testing.T) {
	api := slack.New(config.Token, slack.OptionDebug(true))
	ctx := context.Background()
	combatUser := "U011S918AP4"
	randomChannel := "C011Z205A8L"
	err := invite(ctx, api, randomChannel, []string{combatUser})
	if err != nil {
		t.Fatal(err)
	}
}

func TestInviteAllMCG(t *testing.T) {
	api := slack.New(config.Token, slack.OptionDebug(true))
	ctx := context.Background()
	users, channels, err := inviteAllMCG(ctx, api, config.Prefix)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", resultMessage(users, channels))
}

func TestInvitePubSub(t *testing.T) {
	ctx := context.Background()
	m := PubSubMessage{Data: []byte("C0D595XEG")}
	err := InvitePubSub(ctx, m)
	if err != nil {
		t.Fatal(err)
	}
}