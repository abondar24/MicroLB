package main

import (
	"bufio"
	"os"
	"strings"
)

func LoadBackends(configFile string) ([]string, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var backends []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		backend := strings.TrimSpace(scanner.Text())
		if backend != "" {
			backends = append(backends, backend)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return backends, nil
}
