package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	cache := New[string, any]()

	t.Run("simple", func(t *testing.T) {
		v := 2
		cache.Set("simple", v, time.Hour)
		t.Log(cache.items["simple"].value)

		assert.Equal(t, v, cache.items["simple"].value)
	})
	t.Run("a map!!!", func(t *testing.T) {
		v := make(map[string]float32)
		v["ayo"] = 123
		cache.Set("map", v, time.Hour)

		assert.Equal(t, v, cache.items["map"].value)
	})
}

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		duration time.Duration
		wantBool bool
	}{
		{
			name:     "simple",
			value:    123,
			duration: time.Hour,
			wantBool: true,
		},
		{
			name:     "map",
			value:    make(map[string]float32),
			duration: time.Hour,
			wantBool: true,
		},
		{
			name:     "expired",
			value:    123,
			duration: time.Nanosecond,
			wantBool: false,
		},
	}

	cache := New[string, any]()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache.Set(test.name, test.value, test.duration)

			got, ok := cache.Get(test.name)

			assert.Equal(t, test.wantBool, ok)
			assert.Equal(t, test.value, got)
		})
	}
}
