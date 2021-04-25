package ubot

type Configuration struct {
	APIToken   string `json:"api_token"`
	ServerPort string `json:"server_port"`
	WebhookUrl string `json:"webhook_url"`
	WorkerNo   int    `json:"worker_no"`
}
