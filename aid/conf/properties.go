package conf

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/qinchende/gofast/aid/iox"
)

// ConfigKVError represents a configuration error message.
type ConfigKVError struct {
	error
	message string
}

// ConfigKV interface provides the means to access configuration.
type ConfigKV interface {
	GetString(key string) string
	SetString(key, value string)
	GetInt(key string) int
	SetInt(key string, value int)
	ToString() string
}

// ConfigKV config is a key/value pair based configuration structure.
type baseConfigKV struct {
	kvs  map[string]string
	lock sync.RWMutex
}

// Loads the kvs into a kvs configuration instance.
// Returns an error that indicates if there was a problem loading the configuration.
func LoadConfigKV(filename string) (ConfigKV, error) {
	lines, err := iox.ReadTextLines(filename, iox.WithoutBlank(), iox.OmitWithPrefix("#"))
	if err != nil {
		return nil, err
	}

	raw := make(map[string]string)
	for i := range lines {
		pair := strings.Split(lines[i], "=")
		if len(pair) != 2 {
			// invalid property format
			return nil, &ConfigKVError{
				message: fmt.Sprintf("invalid property format: %s", pair),
			}
		}

		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])
		raw[key] = value
	}

	return &baseConfigKV{
		kvs: raw,
	}, nil
}

func (config *baseConfigKV) GetString(key string) string {
	config.lock.RLock()
	ret := config.kvs[key]
	config.lock.RUnlock()

	return ret
}

func (config *baseConfigKV) SetString(key, value string) {
	config.lock.Lock()
	config.kvs[key] = value
	config.lock.Unlock()
}

func (config *baseConfigKV) GetInt(key string) int {
	config.lock.RLock()
	// default 0
	value, _ := strconv.Atoi(config.kvs[key])
	config.lock.RUnlock()

	return value
}

func (config *baseConfigKV) SetInt(key string, value int) {
	config.lock.Lock()
	config.kvs[key] = strconv.Itoa(value)
	config.lock.Unlock()
}

// Dumps the configuration internal map into a string.
func (config *baseConfigKV) ToString() string {
	config.lock.RLock()
	ret := fmt.Sprintf("%s", config.kvs)
	config.lock.RUnlock()

	return ret
}

// Returns the error message.
func (configError *ConfigKVError) Error() string {
	return configError.message
}

// Builds a new kvs configuration structure
func NewProperties() ConfigKV {
	return &baseConfigKV{
		kvs: make(map[string]string),
	}
}
