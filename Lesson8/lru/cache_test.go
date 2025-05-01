package lru

import (
	"container/list"
	"testing"
)

func TestNewLruCache(t *testing.T) {
	type args struct {
		capacity int
	}
	tests := []struct {
		name string
		args args
		want LRUCache
	}{
		{
			name: "Create cache with capacity 2",
			args: args{capacity: 2},
			want: &lruCacheImpl{
				capacity: 2,
				items:    map[string]*list.Element{},
				order:    list.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLruCache(tt.args.capacity)
			if got == nil {
				t.Fatalf("NewLruCache() returned nil")
			}
			cImpl, ok := got.(*lruCacheImpl)
			if !ok {
				t.Fatalf("Expected *lruCacheImpl, got %T", got)
			}
			if cImpl.capacity != tt.want.(*lruCacheImpl).capacity {
				t.Errorf("Expected capacity %v, got %v", tt.want.(*lruCacheImpl).capacity, cImpl.capacity)
			}
		})
	}
}

func Test_lruCacheImpl_Get(t *testing.T) {
	cache := NewLruCache(2)
	cache.Put("a", "1")
	cache.Put("b", "2")

	tests := []struct {
		name      string
		key       string
		wantValue string
		wantOk    bool
	}{
		{
			name:      "Existing key",
			key:       "a",
			wantValue: "1",
			wantOk:    true,
		},
		{
			name:      "Non-existing key",
			key:       "c",
			wantValue: "",
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := cache.Get(tt.key)
			if got != tt.wantValue || ok != tt.wantOk {
				t.Errorf("Get(%q) = (%q, %v), want (%q, %v)", tt.key, got, ok, tt.wantValue, tt.wantOk)
			}
		})
	}
}

func Test_lruCacheImpl_Put(t *testing.T) {
	cache := NewLruCache(2)

	t.Run("Insert and retrieve", func(t *testing.T) {
		cache.Put("x", "100")
		got, ok := cache.Get("x")
		if !ok || got != "100" {
			t.Errorf("Expected (100, true), got (%v, %v)", got, ok)
		}
	})

	t.Run("Evict least recently used", func(t *testing.T) {
		cache.Put("x", "100")
		cache.Put("y", "200")
		cache.Get("x") // x becomes most recently used
		cache.Put("z", "300")

		if _, ok := cache.Get("y"); ok {
			t.Errorf("Expected y to be evicted")
		}
		if val, ok := cache.Get("x"); !ok || val != "100" {
			t.Errorf("Expected x to stay, got %v", val)
		}
		if val, ok := cache.Get("z"); !ok || val != "300" {
			t.Errorf("Expected z, got %v", val)
		}
	})

	t.Run("Update value and order", func(t *testing.T) {
		cache.Put("x", "new100")
		val, ok := cache.Get("x")
		if !ok || val != "new100" {
			t.Errorf("Expected updated value new100, got %v", val)
		}
	})
}
