package invbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/slack-go/slack"
)

var projectID = os.Getenv("GCP_PROJECT")

// client is a global Pub/Sub client, initialized once per instance.
var client *pubsub.Client

func init() {
	// err is pre-declared to avoid shadowing client.
	var err error

	// client is initialized with context.Background() because it should
	// persist between function invocations.
	client, err = pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
}

type attachment struct {
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	ImageURL  string `json:"image_url"`
}

// Message is the a Slack message event.
// see https://api.slack.com/docs/message-formatting
type Message struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []attachment `json:"attachments"`
}

var config = setup()

func InviteAllMCG(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Couldn't read request body: %v", err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", 405)
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Couldn't parse form", 400)
		log.Fatalf("ParseForm: %v", err)
	}

	// Reset r.Body as ParseForm depletes it by reading the io.ReadCloser.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	result, err := verifyWebHook(r, config.Secret)
	if err != nil {
		log.Fatalf("verifyWebhook: %v", err)
	}
	if !result {
		log.Fatalf("signatures did not match.")
	}

	if len(r.Form["text"]) == 0 {
		log.Fatalf("empty text in form")
	}

	msg := r.Form["text"][0]
	topicName := "triger-invbot"
	m := &pubsub.Message{Data: []byte(msg)}
	client.Topic(topicName).Publish(r.Context(), m)

	res := &Message{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("Trying to invite"),
		Attachments:  []attachment{},
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(res); err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func InvitePubSub(ctx context.Context, m PubSubMessage) error {
	log.Printf("Invite")
	api := slack.New(config.Token, slack.OptionDebug(true))
	inviteAllMCG(context.Background(), api, config.Prefix)
	return nil
}
