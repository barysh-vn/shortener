package model

type ShortenerConfig struct {
	Address  *ShortenerAddress
	BaseURL  string
	FilePath string
	DbDSN    string
}
