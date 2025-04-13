package cache

type Cache[K comparable, V any] interface {
	Add(key K, value V)
	Get(key K) (V, bool)
	Delete(key K)
}
