package iox

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountLines(t *testing.T) {
	const val = `1
2
3
4`
	file, err := ioutil.TempFile(os.TempDir(), "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(val)
	file.Close()
	lines, err := CountLines(file.Name())
	assert.Nil(t, err)
	assert.Equal(t, 4, lines)
}
