package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadBackends(t *testing.T) {
	content := "back1\nback2\nback3"

	file, err := os.Create("backends.txt")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.Write([]byte(content))
	assert.NoError(t, err)
	file.Close()

	backends, err := LoadBackends(file.Name())
	assert.NoError(t, err)

	expected := []string{"back1", "back2", "back3"}

	assert.Equal(t, expected, backends)

}
