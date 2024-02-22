package mathx

import (
	"github.com/qinchende/gofast/aid/randx"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxInt(t *testing.T) {
	cases := []struct {
		a      int
		b      int
		expect int
	}{
		{
			a:      0,
			b:      1,
			expect: 1,
		},
		{
			a:      0,
			b:      -1,
			expect: 0,
		},
		{
			a:      1,
			b:      1,
			expect: 1,
		},
	}

	for _, each := range cases {
		each := each
		t.Run(randx.Rand(), func(t *testing.T) {
			actual := MaxInt(each.a, each.b)
			assert.Equal(t, each.expect, actual)
		})
	}
}

func TestMinInt(t *testing.T) {
	cases := []struct {
		a      int
		b      int
		expect int
	}{
		{
			a:      0,
			b:      1,
			expect: 0,
		},
		{
			a:      0,
			b:      -1,
			expect: -1,
		},
		{
			a:      1,
			b:      1,
			expect: 1,
		},
	}

	for _, each := range cases {
		t.Run(randx.Rand(), func(t *testing.T) {
			actual := MinInt(each.a, each.b)
			assert.Equal(t, each.expect, actual)
		})
	}
}
