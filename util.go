package versionedconfig

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck // best-effort close on read

	return io.ReadAll(resp.Body)
}

func getConfigType(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) > 1 {
		cfType := strings.ToLower(ext[1:])
		switch cfType {
		case "yml":
			return "yaml"
		default:
			return cfType
		}
	}

	return ""
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func readConfiguration(filename string) ([]byte, error) {
	switch {
	case filename == "":
		return nil, errors.New("filename not specified")
	case filename == "-":
		return io.ReadAll(os.Stdin)
	case isURL(filename):
		return download(filename)
	default:
		return os.ReadFile(filename)
	}
}
