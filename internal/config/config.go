package config

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type ShortenerConfig struct {
	Address ShortenerAddress
	URL     string
}

type ShortenerAddress struct {
	Host string
	Port int
}

func (c *ShortenerAddress) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *ShortenerAddress) Set(value string) error {
	values := strings.Split(value, ":")
	if len(values) != 2 {
		return fmt.Errorf("invalid shortener URL: %s", value)
	}

	port, err := strconv.Atoi(values[1])
	if err != nil {
		return err
	}

	c.Host = values[0]
	c.Port = port
	return nil
}

var (
	ShortenerAddressConfig = &ShortenerAddress{
		Host: "localhost",
		Port: 8080,
	}
	ResultURL string
)

func ParseFlags() {
	flag.Var(ShortenerAddressConfig, "a", "Shortener address (host:port)")
	flag.StringVar(&ResultURL, "b", "http://"+ShortenerAddressConfig.String()+"/", "Shortener result URL")
	flag.Parse()
}

func GetShortenerConfig() *ShortenerConfig {
	return &ShortenerConfig{
		Address: *ShortenerAddressConfig,
		URL:     ResultURL,
	}
}
