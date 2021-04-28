package cache

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/skpr/sso-auth/pkg/types"
)

func TestTokensCache(t *testing.T) {
	
	t.Cleanup(func() {
		_ = os.RemoveAll(os.TempDir() + "/TestToken")
	})

	expiry := time.Now().UTC().Add(300 * time.Second).Truncate(time.Second)

	startURL := "http://example.com"

	tokenInfo := &types.Token{
		AccessToken: "ABCDEFGHIJKLMNOP1234567890",
		StartURL:    startURL,
		Region:      "ap-southeast-2",
		ExpiresAt:   expiry,
	}

	cache := NewTokenCache()
	defaultCacheLocation = func() string {
		return os.TempDir() + "/TestToken"
	}

	err := cache.Put(tokenInfo)
	assert.Nil(t, err)

	tokenInfo, err = cache.Get(startURL)
	assert.Nil(t, err)

	assert.Equal(t, "ABCDEFGHIJKLMNOP1234567890", tokenInfo.AccessToken, "access_token was set")
	assert.Equal(t, startURL, tokenInfo.StartURL)
	assert.Equal(t, "ap-southeast-2", tokenInfo.Region)
	assert.Equal(t, expiry, tokenInfo.ExpiresAt)
}

func TestHasExpired(t *testing.T) {
	expiry := time.Now().UTC().Add(-300 * time.Second).Truncate(time.Second)
	tokenInfo := &types.Token{
		AccessToken: "ABCDEFGHIJKLMNOP1234567890",
		StartURL:    "http://example.com",
		Region:      "ap-southeast-2",
		ExpiresAt:   expiry,
	}

	assert.True(t, tokenInfo.Expired())
}
