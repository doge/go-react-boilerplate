package envloader

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var loadOnce sync.Once

func LoadDotEnv() {
	loadOnce.Do(func() {
		appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))

		if appEnv == "" || appEnv == "local" {
			_ = godotenv.Load(".env")
		}
	})
}
