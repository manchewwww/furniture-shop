package config

type Config struct {
	DB           DBConfig `json:"DB"`
	CORS_ORIGINS []string `json:"CORS_ORIGINS"`
}

type DBConfig struct {
	NAME string `json:"NAME"`
	HOST string `json:"HOST"`
	PORT uint   `json:"PORT"`
	SSL  string `json:"SSL"`
}
