package cache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/skpr/sso-auth/pkg/types"
)

// TokenCache handles caching oauth2 types.TokenInfo.
type TokenCache struct {
	cacheFile string
}

// NewTokenCache creates a new instance.
func NewTokenCache(filename string) *TokenCache {
	return &TokenCache{
		cacheFile: filename,
	}
}

// Get will return the oauth token from cache.
func (c *TokenCache) Get() (*types.TokenInfo, error) {

	var tokenInfo *types.TokenInfo

	if _, err := os.Stat(c.cacheFile); os.IsNotExist(err) {
		return &types.TokenInfo{}, err
	}

	data, err := ioutil.ReadFile(c.cacheFile)
	if err != nil {
		return &types.TokenInfo{}, err
	}

	err = json.Unmarshal(data, &tokenInfo)
	if err != nil {
		return &types.TokenInfo{}, err
	}

	return tokenInfo, nil
}

// Put writes an oauth token to cache.
func (c *TokenCache) Put(token *types.TokenInfo) error {

	// Create parent directory if it doesn't exist.
	if _, err := os.Stat(c.cacheFile); os.IsNotExist(err) {
		dir := path.Dir(c.cacheFile)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.cacheFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Delete the cache file.
func (c *TokenCache) Delete(token *types.TokenInfo) error {
	err := os.Remove(c.cacheFile)
	if err != nil {
		return err
	}
	return nil
}
