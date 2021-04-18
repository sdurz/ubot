package ubot

import "encoding/json"

type Configuration struct {
	APIToken   string            `json:"api_token"`
	ServerPort string            `json:"server_port"`
	WebhookUrl string            `json:"webhook_url"`
	LongPoll   bool              `json:"long_poll"`
	Params     map[string]string `json:"params"`
}

func (c *Configuration) Parse(data []byte) error {
	return json.Unmarshal(data, c)
}
