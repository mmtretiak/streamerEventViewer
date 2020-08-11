package config

// for now simply return mocked structure instead of reading from config file or cloud secrets
func New() Config {
	return Config{
		JWTConfig: JWTConfig{
			Algorithm:  "HS256",
			Secret:     "testkeyforsignig",
			TTLMinutes: 10080, // 7 days for now
		},
		DB: DB{
			User:     "mmtretiak",
			Password: "postgres",
			Dialect:  "postgres",
			DBName:   "d23ai2sqpp4msf\n",
			Port:     5432,
			Host:     "localhost",
		},
		OauthConfig: OauthConfig{
			RedirectURI:  "https://stark-escarpment-52058.herokuapp.com/redirect",
			Scopes:       []string{"user:read:email", "clips:edit"},
			ClientID:     "k9rce279ezyjl3tafvhvza6pvj55cb",
			ClientSecret: "c60qyyfj3tijxsnau2mqy9ecpqfcl2",
		},
		Server: Server{
			Port: 8080,
			Host: "localhost",
		},
		Jobs: Jobs{
			ViewUpdaterJob: ViewUpdaterJob{
				Schedule: "58 8 * * *",
			},
		},
	}
}

type Config struct {
	JWTConfig   JWTConfig   `json:"jwt_config"`
	OauthConfig OauthConfig `json:"oauth_config"`
	DB          DB          `json:"db"`
	Server      Server      `json:"server"`
	Jobs        Jobs        `json:"jobs"`
}

type JWTConfig struct {
	Algorithm  string `json:"algorithm"`
	Secret     string `json:"secret"`
	TTLMinutes int64  `json:"ttl_minutes"`
}

// required scope
// user:read:email
// clips:edit
type OauthConfig struct {
	RedirectURI  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
}

type DB struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	Dialect  string `json:"dialect"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type Server struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type Jobs struct {
	ViewUpdaterJob ViewUpdaterJob `json:"view_updater_job"`
}

type ViewUpdaterJob struct {
	Schedule string `json:"schedule"`
}
