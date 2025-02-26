package main

import (
	"context"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	RunSocketMode()
}

func resultMessage(users []slack.User, channels []slack.Channel) string {
	msg := "完了!\nMCGs: "
	for i, user := range users {
		if i != 0 {
			msg += ", "
		}
		msg += user.Name
	}
	msg += "\nChannels: "
	for i, channel := range channels {
		if i != 0 {
			msg += ", "
		}
		msg += channel.Name
	}
	return msg
}

func RunSocketMode() {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		log.Fatal("SLACK_APP_TOKEN is required")
	}
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("SLACK_BOT_TOKEN is required")
	}
	prefix := os.Getenv("PREFIX_CHANNEL")
	if prefix == "" {
		log.Fatal("PREFIX_CHANNEL is required")
	}

	api := slack.New(
		os.Getenv("SLACK_BOT_TOKEN"),
		slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
	client := socketmode.New(
		api,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	_, authTestErr := api.AuthTest()
	if authTestErr != nil {
		log.Fatalf("SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
	}

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					continue
				}
				client.Ack(*evt.Request)
				if cmd.Command != "/inviteallmcg" {
					continue
				}
				channelID := cmd.ChannelID
				_, _, err := api.PostMessage(channelID, slack.MsgOptionText("MCGを招待します", false))
				if err != nil {
					log.Printf("PostMessage error: %v", err)
					continue
				}
				continue
				users, channels, err := inviteAllMCG(context.Background(), api, prefix)
				if err != nil {
					log.Printf("inviteAllMCG error: %v", err)
					api.PostMessage(channelID, slack.MsgOptionText("エラーが発生しました", false))
					continue
				}
				msg := resultMessage(users, channels)
				_, _, err = api.PostMessage(channelID, slack.MsgOptionText(msg, false))
				if err != nil {
					log.Printf("PostMessage error: %v", err)
				}
			}
		}
	}()

	if err := client.Run(); err != nil {
		log.Fatalf("client.Run: %v", err)
	}
}
