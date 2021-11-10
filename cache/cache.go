package cache

type Cache interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

type NullCache struct {}
func (c NullCache) Get(_ string) (interface{}, bool) {
	return nil, false
}
func (c NullCache) Set(_ string, _ interface{}) {}