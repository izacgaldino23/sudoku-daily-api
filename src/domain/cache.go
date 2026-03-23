package domain

type (
	Cache interface {
		Get(key string) (any, bool)
		Set(key string, value any)
		Flush()
	}
)
