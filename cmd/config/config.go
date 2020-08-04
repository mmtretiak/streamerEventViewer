package config

// for now simply return mocked structure instead of reading from config file or cloud secrets
func New() Config {
	return Config{
		JWTConfig: JWTConfig{
			Algorithm: "HS256",
			Secret:    "testkeyforsignig",
		},
		DB: DB{
			User:     "postgres",
			Password: "postgres",
			Dialect:  "postgres",
			DBName:   "develop",
		},
	}
}

type Config struct {
	JWTConfig   JWTConfig   `json:"jwt_config"`
	OauthConfig OauthConfig `json:"oauth_config"`
	DB          DB          `json:"db"`
}

type JWTConfig struct {
	Algorithm  string `json:"algorithm"`
	Secret     string `json:"secret"`
	TTLMinutes int64  `json:"ttl_minutes"`
}

// required scope
// user:read:email
//
type OauthConfig struct {
	RedirectURL  string   `json:"redirect_url"`
	Scopes       []string `json:"scopes"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
}

type DB struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	Dialect  string `json:"dialect"`
}
