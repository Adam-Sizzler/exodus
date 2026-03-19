package assets

import (
	"fmt"
	"os"
	"path/filepath"
)

func DetectPath() (string, error) {
	candidates := []string{
		"/opt/app/frontend",
		filepath.Clean(filepath.Join("frontend", "dist")),
		filepath.Clean(filepath.Join("..", "frontend", "dist")),
		filepath.Clean(filepath.Join("backend", "dev_frontend")),
		filepath.Clean(filepath.Join("..", "backend", "dev_frontend")),
		filepath.Clean("frontend"),
		filepath.Clean(filepath.Join("..", "frontend")),
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}

		info, err := os.Stat(filepath.Join(candidate, "index.html"))
		if err == nil && !info.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("frontend assets directory with index.html not found")
}
