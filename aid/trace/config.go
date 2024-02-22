package trace

// TraceName represents the tracing name.
const TraceName = "GoFast"

// A Config is a opentelemetry config.
type Config struct {
	Name     string  `v:""`
	Endpoint string  `v:""`
	Sampler  float64 `v:"def=1.0"`
	Batcher  string  `v:"def=jaeger,enum=jaeger|zipkin|grpc"`
}
