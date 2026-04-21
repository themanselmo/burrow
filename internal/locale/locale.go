package locale

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

var strings_ map[string]interface{}

func Load() error {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..", "locales", "en.yml")

	data, err := os.ReadFile(root)
	if err != nil {
		return fmt.Errorf("failed to load locale file: %w", err)
	}

	return yaml.Unmarshal(data, &strings_)
}

func T(key string) string {
	parts := strings.Split(key, ".")
	var cur interface{} = strings_

	for _, part := range parts {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return key
		}
		cur, ok = m[part]
		if !ok {
			return key
		}
	}

	if s, ok := cur.(string); ok {
		return s
	}
	return key
}

func Tf(key string, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}
