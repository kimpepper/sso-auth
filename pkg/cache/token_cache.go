package cache

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/skpr/sso-auth/pkg/types"
)

// TokenCache handles caching oauth2 types.Token.
type TokenCache struct{}

var defaultCacheLocation func() string

func defaultCacheLocationImpl() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".aws", "sso", "cache")
}

func init() {
	defaultCacheLocation = defaultCacheLocationImpl
}

// NewTokenCache creates a new instance.
func NewTokenCache() *TokenCache {
	return &TokenCache{}
}

// Get will return the oauth token from cache.
func (c *TokenCache) Get(startURL string) (*types.Token, error) {
	key, err := getTokenFileName(startURL)
	var token *types.Token

	cacheFile := filepath.Join(defaultCacheLocation(), key)
	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return &types.Token{}, err
	}

	err = json.Unmarshal(data, &token)
	if err != nil {
		return &types.Token{}, err
	}

	return token, nil
}

// Put writes an oauth token to cache.
func (c *TokenCache) Put(token *types.Token) error {

	key, err := getTokenFileName(token.StartURL)
	cacheFile := filepath.Join(defaultCacheLocation(), key)

	// Create parent directory if it doesn't exist.
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		dir := path.Dir(cacheFile)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	data, err := json.Marshal(token)
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
func (c *TokenCache) Invalidate(startURL string) error {
	key, err := getTokenFileName(startURL)
	cacheFile := filepath.Join(defaultCacheLocation(), key)
	err = os.Remove(cacheFile)
	if err != nil {
		return err
	}
	return nil
}

func getTokenFileName(url string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(url))
	if err != nil {
		return "", err
	}
	return strings.ToLower(hex.EncodeToString(hash.Sum(nil))) + ".json", nil
}
