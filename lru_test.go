package lru

import (
	"testing"
)

func TestAddGet(t *testing.T) {
	lruCache := NewLruCache(1)
	key := "1"
	expectedValue := "1"

	lruCache.Add(key, expectedValue)

	if v, ok := lruCache.Get(key); ok {
		s, ok := v.(string)
		if !ok {
			t.Errorf("expected string, got %v", v)
		}

		if s != expectedValue {
			t.Errorf("got %s, wanted %s", s, expectedValue)
		}
	} else {
		t.Error("get failed, could not find value")
	}

}

func TestRemove(t *testing.T) {
	lruCache := NewLruCache(1)
	key := "1"
	expectedValue := "1"

	lruCache.Add(key, expectedValue)

	if v, ok := lruCache.Get(key); ok {
		s, ok := v.(string)
		if !ok {
			t.Errorf("expected string, got %v", v)
		}

		if s != expectedValue {
			t.Errorf("got %s, wanted %s", s, expectedValue)
		}
	} else {
		t.Error("get failed, could not find value")
	}

	// remove now

	lruCache.Remove(key)

	if _, ok := lruCache.Get(key); ok {
		t.Error("should not have found element")
	}

}

func TestRemoveOldest(t *testing.T) {
	lruCache := NewLruCache(2)
	keyOldest := "1"
	expectedValueOldest := "1"
	keyNewest := "2"
	expectedValueNewest := "2"

	lruCache.Add(keyOldest, expectedValueOldest)
	lruCache.Add(keyNewest, expectedValueNewest)
	lruCache.RemoveOldest()

	if v, ok := lruCache.Get(keyNewest); ok {
		s, ok := v.(string)
		if !ok {
			t.Errorf("expected string, got %v", v)
		}

		if s != expectedValueNewest {
			t.Errorf("got %s, wanted %s", s, expectedValueNewest)
		}
	} else {
		t.Error("get failed, could not find value")
	}

	// remove now

	if _, ok := lruCache.Get(expectedValueOldest); ok {
		t.Error("should not have found element")
	}

}

func TestLen(t *testing.T) {
	lruCache := NewLruCache(2)
	keyOldest := "1"
	expectedValueOldest := "1"
	keyNewest := "2"
	expectedValueNewest := "2"

	lruCache.Add(keyOldest, expectedValueOldest)
	lruCache.Add(keyNewest, expectedValueNewest)

	if lruCache.Len() != 2 {
		t.Errorf("got %d, expected %d", lruCache.Len(), 2)
	}

}
