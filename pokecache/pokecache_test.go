package pokecache

import (
	"testing"
	"time"

	"github.com/DreamyMemories/pokedex-cli/types"
	"github.com/stretchr/testify/assert"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	apiResponse := types.ApiResponse{
		Next:     "abc",
		Previous: "abcd",
		Results:  make([]types.LocationArea, 0),
	}

	cache.Add("testKey", apiResponse)

	retrievedValue, cached := cache.Get("testKey")
	if !cached {
		t.Errorf("Expected to find key that was added")
	}

	assert.Equal(t, apiResponse, retrievedValue, "Retrieved Value should match each other")
}

func TestCacheExpiration(t *testing.T) {
	cache := NewCache(1 * time.Second)

	apiResponse := types.ApiResponse{
		Next:     "abc",
		Previous: "abcd",
		Results:  make([]types.LocationArea, 0),
	}

	cache.Add("testKey", apiResponse)
	time.Sleep(2 * time.Second)

	_, found := cache.Get("testkey")
	assert.False(t, found, "Expected key to be expired and item deleted")
}
