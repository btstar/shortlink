package util

import (
	"errors"
	"strings"
)

// 检查Url是否合法
func CheckUrl(url string) error {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return nil
	}
	return errors.New("Url must start with http or https...")
}
