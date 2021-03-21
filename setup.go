package invbot

type configuration struct {
	Token  string `json:"SLACK_TOKEN" required:"true"`
	Secret string `json:"SLACK_SIGNING_SECRET" required:"true"`
	Prefix string `json:"PREFIX_CHANNEL" required:"true"`
}
