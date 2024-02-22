package iox

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	const val = `1
2
3
4`
	reader := strings.NewReader(val)
	scanner := NewTextLineScanner(reader)
	var lines []string
	for scanner.Scan() {
		line, err := scanner.Line()
		assert.Nil(t, err)
		lines = append(lines, line)
	}
	assert.EqualValues(t, []string{"1", "2", "3", "4"}, lines)
}
