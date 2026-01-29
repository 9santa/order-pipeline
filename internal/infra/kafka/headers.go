package kafka

import (
	"github.com/segmentio/kafka-go"
)

func HeaderGetInt(headers []kafka.Header, key string) (int, bool) {
	for _, h := range headers {
		if h.Key == key {
			n := 0
			for _, c := range h.Value {
				if c < '0' || c > '9' {
					return 0, false
				}
				n = n*10 + int(c-'0')
			}
			return n, true
		}
	}
	return 0, false
}

func HeaderSet(headers []kafka.Header, key string, value []byte) []kafka.Header {
	out := make([]kafka.Header, 0, len(headers)+1)
	replaced := false
	for _, h := range headers {
		// If key already is in the header, replace the value
		if h.Key == key {
			out = append(out, kafka.Header{Key: key, Value: value})
			replaced = true
		} else {
			out = append(out, h)
		}
	}
	// If key wasn't in the header, add it
	if !replaced {
		out = append(out, kafka.Header{Key: key, Value: value})
	}
	return out
}
