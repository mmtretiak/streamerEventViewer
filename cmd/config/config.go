package config

// for now simply return mocked structure instead of reading from config file or cloud secrets
func New() Config {
	return Config{
		JWTConfig: JWTConfig{
			Algorithm: "HS256",
			Secret:    "testkeyforsignig",
		},
	}
}

type Config struct {
	JWTConfig JWTConfig
}

type JWTConfig struct {
	Algorithm  string `json:"algorithm"`
	Secret     string `json:"secret"`
	TTLMinutes int64  `json:"ttl_minutes"`
}
