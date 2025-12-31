package config

type EnvConfig struct {
	DBUser              string
	DBPass              string
	JWTSecret           string
	StripeSecretKey     string
	StripeWebhookSecret string
	EmailSenderHost     string
	EmailSenderPort     string
	EmailSenderUser     string
	EmailSenderPass     string
	EmailSenderFrom     string
}
