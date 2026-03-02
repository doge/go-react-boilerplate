package models

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`

		CORS struct {
			AllowedOrigins []string `yaml:"allowed_origins"`
		} `yaml:"cors"`
	} `yaml:"server"`

	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
		IP       string `yaml:"ip"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	Auth struct {
		JWTSecret         string `yaml:"jwt_secret"`
		Issuer            string `yaml:"issuer"`
		AccessTTLMinutes  int    `yaml:"access_ttl_minutes"`
		RefreshTTLDays    int    `yaml:"refresh_ttl_days"`
		RefreshCookieName string `yaml:"refresh_cookie_name"`
		SecureCookies     bool   `yaml:"secure_cookies"`
		RefreshSameSite   string `yaml:"refresh_same_site"`
	} `yaml:"auth"`
}

func (c Config) Address() string {
	host := c.Server.Host
	if host == "" {
		host = "127.0.0.1"
	}

	port := c.Server.Port
	if port == "" {
		port = "8000"
	}

	return host + ":" + port
}

func (c Config) AllowedOrigins() []string {
	if len(c.Server.CORS.AllowedOrigins) > 0 {
		return c.Server.CORS.AllowedOrigins
	}

	return []string{
		"http://127.0.0.1:3000",
		"http://localhost:3000",
		"http://127.0.0.1:5173",
		"http://localhost:5173",
	}
}

func (c Config) IsAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range c.AllowedOrigins() {
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}

	return false
}

func (c Config) Validate() error {
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("auth.jwt_secret is required")
	}
	if c.Auth.RefreshCookieName == "" {
		return fmt.Errorf("auth.refresh_cookie_name is required")
	}
	if c.Auth.RefreshSameSite == "" {
		c.Auth.RefreshSameSite = "strict"
	}
	if c.Auth.AccessTTLMinutes <= 0 {
		return fmt.Errorf("auth.access_ttl_minutes must be > 0")
	}
	if c.Auth.RefreshTTLDays <= 0 {
		return fmt.Errorf("auth.refresh_ttl_days must be > 0")
	}

	return nil
}

func (c Config) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		c.Database.Username, c.Database.Password,
		c.Database.IP, c.Database.Port,
	)
}

func LoadConfig(configName string) *Config {
	var cfg *Config

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(cwd + "/config/" + configName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	applyEnvOverrides(cfg)

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	return cfg
}

func applyEnvOverrides(cfg *Config) {
	if v := strings.TrimSpace(os.Getenv("SERVER_HOST")); v != "" {
		cfg.Server.Host = v
	}
	if v := strings.TrimSpace(os.Getenv("SERVER_PORT")); v != "" {
		cfg.Server.Port = v
	}
	if v := strings.TrimSpace(os.Getenv("JWT_SECRET")); v != "" {
		cfg.Auth.JWTSecret = v
	}
	if v := strings.TrimSpace(os.Getenv("JWT_ISSUER")); v != "" {
		cfg.Auth.Issuer = v
	}
	if v := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); v != "" {
		cfg.Server.CORS.AllowedOrigins = splitCSV(v)
	}
}

func splitCSV(value string) []string {
	raw := strings.Split(value, ",")
	result := make([]string, 0, len(raw))
	for _, part := range raw {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}
