package mapping

type (
	Valuer interface {
		Value(key string) (any, bool)
	}

	MapValuer map[string]any
)

func (mv MapValuer) Value(key string) (any, bool) {
	v, ok := mv[key]
	return v, ok
}
