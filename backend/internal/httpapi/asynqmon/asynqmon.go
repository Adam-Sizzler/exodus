package asynqmon

import (
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/jobqueue"

	"github.com/hibiken/asynqmon"
)

func NewAsynqmon(cfg *config.BackendConfig) (http.Handler, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	redisOpt := jobqueue.BuildAsynqRedisOpt(cfg)

	basePath := strings.TrimSuffix(cfg.Panel.BasePath, "/")
	if basePath == "" {
		basePath = "/"
	}

	rootPath := "/api/queues"
	if basePath != "/" {
		rootPath = basePath + "/api/queues"
	}

	h := asynqmon.New(asynqmon.Options{
		RootPath:     rootPath,
		RedisConnOpt: redisOpt,
	})

	if basePath == "/" {
		return h, nil
	}

	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = basePath + r.URL.Path
		h.ServeHTTP(w, r)
	})

	return wrapped, nil
}
