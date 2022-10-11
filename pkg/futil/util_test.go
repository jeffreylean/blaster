package futil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	b, err := ReadFile("file:///Users/jeffreylean/project/personal/blaster/examples/payload.json")

	assert.NoError(t, err)
	assert.NotNil(t, b)
	fmt.Println(string(b))
}
