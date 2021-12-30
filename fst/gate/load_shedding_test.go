package gate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSheddingStat(t *testing.T) {
	st := createSheddingStat("any")
	for i := 0; i < 3; i++ {
		st.Total()
	}
	for i := 0; i < 5; i++ {
		st.Pass()
	}
	for i := 0; i < 7; i++ {
		st.Drop()
	}
	result := st.reset()
	assert.Equal(t, int64(3), result.total)
	assert.Equal(t, int64(5), result.pass)
	assert.Equal(t, int64(7), result.drop)
}
