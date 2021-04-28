package cache

import (
	"encoding/json"
	"fmt"
	"github.com/skpr/sso-auth/pkg/types"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	clientIDPrefix = "skpr-client-id"
)

// ClientCredsCache handles caching oauth2 client credentials.
type ClientCredsCache struct{}

// NewClientCredsCache creates a new instance.
func NewClientCredsCache() *ClientCredsCache {
	return &ClientCredsCache{}
}

// Get will return the client creds from cache.
func (c *ClientCredsCache) Get(region string) (*types.ClientCreds, error) {
	key := getClientCredsFilename(region)
	cacheFile := filepath.Join(defaultCacheLocation(), key)

	var clientCreds *types.ClientCreds

	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return &types.ClientCreds{}, err
	}

	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return &types.ClientCreds{}, err
	}

	err = json.Unmarshal(data, &clientCreds)
	if err != nil {
		return &types.ClientCreds{}, err
	}

	return clientCreds, nil
}

// Put writes client credentials to cache.
func (c *ClientCredsCache) Put(region string, clientCreds *types.ClientCreds) error {
	key := getClientCredsFilename(region)
	cacheFile := filepath.Join(defaultCacheLocation(), key)
	// Create parent directory if it doesn't exist.
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		dir := path.Dir(cacheFile)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	data, err := json.Marshal(clientCreds)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(cacheFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Invalidate the cache file.
func (c *ClientCredsCache) Invalidate(region string) error {
	key := getClientCredsFilename(region)
	cacheFile := filepath.Join(defaultCacheLocation(), key)
	err := os.Remove(cacheFile)
	if err != nil {
		return err
	}
	return nil
}

// getClientCredsFilename returns the filename for the client credentials.
func getClientCredsFilename(region string) string {
	return fmt.Sprintf("%s-%s.json", clientIDPrefix, region)
}
