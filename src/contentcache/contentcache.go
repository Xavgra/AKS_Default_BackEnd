package contentcache

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
)

type CacheItem struct {
	keyName string
	value   []byte
}

type ContentCache struct {
	cacheData []CacheItem
}

func (cache *ContentCache) newContentCache() *ContentCache {
	return new(ContentCache)
}

func (cache *ContentCache) AddItem(keyName string, filePath string) error {

	itemToCache := new(CacheItem)
	itemToCache.keyName = keyName
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Printf("Unexpected error caching key %s file %s", keyName, filePath)
		return err
	}

	itemToCache.value = content
	cache.cacheData = append(cache.cacheData, *itemToCache)

	log.Printf("Added %s to cache", keyName)
	return nil
}

func (cache *ContentCache) GetItemReader(keyName string) (io.Reader, error) {
	itemCached, ok := cache.searchItem(keyName)
	if !ok {
		err := errors.New("KeyName not found")
		return nil, err
	}
	return bytes.NewReader(itemCached.value), nil
}

func (cache *ContentCache) searchItem(keyName string) (CacheItem, bool) {

	for i := range cache.cacheData {
		if cache.cacheData[i].keyName == keyName {
			return cache.cacheData[i], true
		}
	}
	return *new(CacheItem), false
}
