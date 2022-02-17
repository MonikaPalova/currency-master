package config

import "time"

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
	return &Mysql{User: user, Password: password}
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

	usersApiV1        = "/api/v1/users"
	assetsApiV1       = "/api/v1/assets"
	userAssetsApiV1   = "/api/v1/users/{username}/assets"
	acquisitionsApiV1 = "/api/v1/acquisitions"
)

// Application configuration
type App struct {
	Host string
	Port string

	UsersApiV1        string
	AssetsApiV1       string
	UserAssetsApiV1   string
	AcquisitionsApiV1 string
}

func NewApp() *App {
	return &App{Host: host, Port: port, UserAssetsApiV1: userAssetsApiV1, UsersApiV1: usersApiV1, AssetsApiV1: assetsApiV1, AcquisitionsApiV1: acquisitionsApiV1}
}

const (
	sessionIDLength    = 64
	sessionDuration    = time.Hour
	sessionCookieName  = "CURRENCY-MASTER-SESSION-ID"
	usernameCookieName = "CURRENCY-MASTER-USERNAME"
)

type Session struct {
	SessionIDLength    int
	SessionDuration    time.Duration
	SessionCookieName  string
	UsernameCookieName string
}

func NewSession() *Session {
	return &Session{SessionIDLength: sessionIDLength, SessionDuration: sessionDuration, SessionCookieName: sessionCookieName, UsernameCookieName: usernameCookieName}
}
