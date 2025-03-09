package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

const (
	bearerEnvName  = "BEARER"
	urlEnvGetData  = "URLGET"
	urlEnvSaveData = "URLSAVE"
)

type Config interface {
	GetBearerToken() string
	GetUrlGetData() string
	GetUrlSaveData() string
}

type config struct {
	urlGetData, bearerToken, urlSaveData string
}

// Загружаем файл в переменные окружения
func Load(path string) (Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	//получаем данные из конфига
	token := os.Getenv(bearerEnvName)
	if len(token) == 0 {
		return nil, errors.New("token is not found")
	}

	urlGet := os.Getenv(urlEnvGetData)
	if len(urlGet) == 0 {
		return nil, errors.New("urlGet is not found")
	}

	urlSave := os.Getenv(urlEnvSaveData)
	if len(urlSave) == 0 {
		return nil, errors.New("urlSave is not found")
	}

	return &config{
		urlGetData:  urlGet,
		urlSaveData: urlSave,
		bearerToken: token,
	}, nil
}

func (c *config) GetBearerToken() string {
	return c.bearerToken
}

func (c *config) GetUrlGetData() string {
	return c.urlGetData
}

func (c *config) GetUrlSaveData() string {
	return c.urlSaveData
}
