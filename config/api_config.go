package config

type ApiConfig struct {
	DatabaseURL string

	AuthServerPort     string
	AuthVerifyEndpoint string

	ResourceServerPort string
}
