package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	env_pkg "serialization/internal"
	"sync"
)

var logger *slog.Logger
var once sync.Once

func createLogger(env env_pkg.EnvType) *slog.Logger {
	writer := os.Stdout

	switch env {
	case env_pkg.ENV_DEBUG:
		return slog.New(
			slog.NewTextHandler(writer, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	case env_pkg.ENV_TESTING:
		return slog.New(
			slog.NewTextHandler(writer, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	case env_pkg.ENV_PREPROD:
		return slog.New(
			slog.NewTextHandler(writer, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	case env_pkg.ENV_PROD:
		return slog.New(
			slog.NewTextHandler(writer, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	}

	log.Fatalf("неизвестное окружение %s", env)
	return nil
}

func FailOnError(msg string, err error) {
	if err != nil {
		error_msg := fmt.Sprintf("%s: %s", msg, err)
		logger.Error(error_msg)
		panic(error_msg)
	}
}

func Init(env env_pkg.EnvType) *slog.Logger {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LstdFlags | log.Llongfile)

	once.Do(func() {
		logger = createLogger(env)
	})
	return logger
}
