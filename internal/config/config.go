package config

type Config struct {
	DB          DBConfig `json:"DB"`
	CORSOrigins []string `json:"CORS_ORIGINS"`
}

type DBConfig struct {
	Name string `json:"NAME"`
	Host string `json:"HOST"`
	Port uint   `json:"PORT"`
	SSL  string `json:"SSL"`
}
