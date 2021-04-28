package cache

import (
	"os"
	"testing"
"time"

"github.com/stretchr/testify/assert"

"github.com/skpr/sso-auth/pkg/types"
)

func TestClientCredsCache(t *testing.T) {

	t.Cleanup(func() {
		_ = os.RemoveAll(os.TempDir() + "/TestCreds")
	})

	expiry := time.Now().Truncate(time.Second).Add(300 * time.Second).UTC()
	issued := time.Now().UTC().Truncate(time.Second)
	
	clientCreds := &types.ClientCreds{
		ClientId:              "abcdef123456",
		IssuedAt:              issued,
		ClientSecret:          "xyzabc98765",
		ExpiresAt:             expiry,
	}

	cache := NewClientCredsCache()
	defaultCacheLocation = func() string {
		return os.TempDir() + "/TestCreds"
	}
	region := "ap-southeast-2"
	err := cache.Put(region, clientCreds)
	assert.Nil(t, err)

	clientCreds, err = cache.Get(region)
	assert.Nil(t, err)

	assert.Equal(t, "abcdef123456", clientCreds.ClientId)
	assert.Equal(t, issued, clientCreds.IssuedAt)
	assert.Equal(t, "xyzabc98765", clientCreds.ClientSecret)
	assert.Equal(t, expiry, clientCreds.ExpiresAt)
}

func TestHExpired(t *testing.T) {
	expiry := time.Now().UTC().Add(-300 * time.Second).Truncate(time.Second)
	issued := time.Now().UTC().Truncate(time.Second)

	clientCreds := &types.ClientCreds{
		ClientId:              "abcdef123456",
		IssuedAt:              issued,
		ClientSecret:          "xyzabc98765",
		ExpiresAt:             expiry,
	}

	assert.True(t, clientCreds.Expired())
}
