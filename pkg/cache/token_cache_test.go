package cache

import (
	"github.com/skpr/sso-auth/pkg/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokensCache(t *testing.T) {

	expiry := time.Now().UTC().Add(300 * time.Second).Truncate(time.Second)

	tokenInfo := &types.TokenInfo{
		AccessToken: "ABCDEFGHIJKLMNOP1234567890",
		StartURL:    "http://example.com",
		Region:      "ap-southeast-2",
		ExpiresAt:   expiry,
	}

	cache := NewTokenCache("/tmp/skpr/cache.json")
	err := cache.Put(tokenInfo)
	assert.Nil(t, err)

	tokenInfo, err = cache.Get()
	assert.Nil(t, err)

	assert.Equal(t, "ABCDEFGHIJKLMNOP1234567890", tokenInfo.AccessToken, "access_token was set")
	assert.Equal(t, "http://example.com", tokenInfo.StartURL)
	assert.Equal(t, "ap-southeast-2", tokenInfo.Region)
	assert.Equal(t, expiry, tokenInfo.ExpiresAt)
}

func TestHasExpired(t *testing.T) {
	expiry := time.Now().UTC().Add(-300 * time.Second).Truncate(time.Second)
	tokenInfo := &types.TokenInfo{
		AccessToken: "ABCDEFGHIJKLMNOP1234567890",
		StartURL:    "http://example.com",
		Region:      "ap-southeast-2",
		ExpiresAt:   expiry,
	}

	assert.True(t, tokenInfo.Expired())
}
