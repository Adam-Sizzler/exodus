package asynqmon

import (
	"fmt"
	"net/http"

	"exodus/internal/config"
	"exodus/internal/jobqueue"

	"github.com/hibiken/asynqmon"
)

func NewAsynqmon(cfg *config.BackendConfig) (http.Handler, error) {
	return NewAsynqmonWithRootPath(cfg, "/api/backend-tools/queues")
}

func NewAsynqmonWithRootPath(cfg *config.BackendConfig, routePath string) (http.Handler, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	redisOpt := jobqueue.BuildAsynqRedisOpt(cfg)

	basePath := cfg.Backend.Trimmed()
	rootPath := routePath
	if basePath != "" {
		rootPath = basePath + routePath
	}

	h := asynqmon.New(asynqmon.Options{
		RootPath:     rootPath,
		RedisConnOpt: redisOpt,
	})

	if basePath == "" {
		return h, nil
	}

	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = basePath + r.URL.Path
		h.ServeHTTP(w, r)
	})

	return wrapped, nil
}
