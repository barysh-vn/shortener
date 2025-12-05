package model

type ShortenerConfig struct {
	Address     *ShortenerAddress
	BaseURL     string
	FilePath    string
	DataBaseDSN string
	Secret      string
}
