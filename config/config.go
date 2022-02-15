package config

const (
	user     = "root"
	password = ""
)

type Mysql struct {
	User     string
	Password string
}

// MySql configuration
func NewMysql() *Mysql {
	return &Mysql{user, password}
}

const (
	assetsUrl    = "https://rest.coinapi.io/v1/assets"
	apiKeyHeader = "X-CoinAPI-Key"
	apiKey       = "D8096E91-86D8-4998-B5B8-C785CE5D58AD"
)

// External Coin API configuration
type CoinAPI struct {
	AssetsUrl    string
	ApiKeyHeader string
	ApiKey       string
}

func NewCoinAPI() *CoinAPI {
	return &CoinAPI{assetsUrl, apiKeyHeader, apiKey}
}

const (
	host = "localhost"
	port = "7777"
)

// Application configuration
type App struct {
	Host string
	Port string
}

func NewApp() *App {
	return &App{host, port}
}
