package config

type EnvConfig struct {
	DBUser              string
	DBPass              string
	JWTSecret           string
	StripeSecretKey     string
	StripeWebhookSecret string
}
