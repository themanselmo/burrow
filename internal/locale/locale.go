package locale

import (
	_ "embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed en.yml
var localeData []byte

var strings_ map[string]interface{}

func Load() error {
	return yaml.Unmarshal(localeData, &strings_)
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
