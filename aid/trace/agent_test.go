package trace

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartAgent(t *testing.T) {
	//logx.Disable()

	const (
		endpoint1 = "localhost:1234"
		endpoint2 = "remotehost:1234"
		endpoint3 = "localhost:1235"
	)
	c1 := Config{
		Name: "foo",
	}
	c2 := Config{
		Name:     "bar",
		Endpoint: endpoint1,
		Batcher:  kindJaeger,
	}
	c3 := Config{
		Name:     "any",
		Endpoint: endpoint2,
		Batcher:  kindZipkin,
	}
	c4 := Config{
		Name:     "bla",
		Endpoint: endpoint3,
		Batcher:  "otlp",
	}
	c5 := Config{
		Name:     "grpc",
		Endpoint: endpoint3,
		Batcher:  "grpc",
	}

	StartAgent(c1)
	StartAgent(c1)
	StartAgent(c2)
	StartAgent(c3)
	StartAgent(c4)
	StartAgent(c5)

	lock.Lock()
	defer lock.Unlock()

	// because remotehost cannot be resolved
	assert.Equal(t, 3, len(agents))
	_, ok := agents[""]
	assert.True(t, ok)
	_, ok = agents[endpoint1]
	assert.True(t, ok)
	_, ok = agents[endpoint2]
	assert.False(t, ok)
}
