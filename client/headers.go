package client

type Headers interface {
	Set(key, value string)
	Keys() []string
	Get(key string) string
	Value() map[string]string
}

type HeadersStore map[string]string

func NewHeaders(h map[string]string) Headers {
	if h == nil {
		h = make(HeadersStore)
	}
	return HeadersStore(h)
}

func (h HeadersStore) Set(key, value string) {
	h[key] = value
}

func (h HeadersStore) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

func (h HeadersStore) Get(key string) string {
	return h[key]
}

func (h HeadersStore) Value() map[string]string {
	return h
}
