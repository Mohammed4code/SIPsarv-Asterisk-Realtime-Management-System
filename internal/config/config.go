package config

import (
	"os"
	"sync"
)

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Asterisk AsteriskConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	mu       sync.RWMutex
}

type AsteriskConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	mu       sync.RWMutex
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		cfg = &Config{
			Server: ServerConfig{
				Port: getEnv("PORT", ":8080"),
			},
			DB: DBConfig{
				Host:     getEnv("DB_HOST", "127.0.0.1"),
				Port:     getEnv("DB_PORT", "3306"),
				Name:     getEnv("DB_NAME", "asterisk"),
				User:     getEnv("DB_USER", "asterisk"),
				Password: getEnv("DB_PASS", "123456"),
			},
			Asterisk: AsteriskConfig{
				Host:     getEnv("ASTERISK_HOST", "127.0.0.1"),
				Port:     getEnv("ASTERISK_PORT", "5038"),
				User:     getEnv("AMI_USER", "admin"),
				Password: getEnv("AMI_PASS", "password"),
			},
		}
	})
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// دوال آمنة للقراءة والكتابة
func (c *AsteriskConfig) Get() (host, port, user, password string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Host, c.Port, c.User, c.Password
}

func (c *AsteriskConfig) Set(host, port, user, password string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if host != "" {
		c.Host = host
	}
	if port != "" {
		c.Port = port
	}
	if user != "" {
		c.User = user
	}
	if password != "" && password != "********" {
		c.Password = password
	}
}

func (c *DBConfig) Get() (host, port, name, user, password string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Host, c.Port, c.Name, c.User, c.Password
}

func (c *DBConfig) Set(host, port, name, user, password string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if host != "" {
		c.Host = host
	}
	if port != "" {
		c.Port = port
	}
	if name != "" {
		c.Name = name
	}
	if user != "" {
		c.User = user
	}
	if password != "" && password != "********" {
		c.Password = password
	}
}