package cache

import (
	"encoding/json"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
)

type CacheType struct {
	Data    []byte `json:"data"`
	Created int64  `json:"created"`
	Expires int64  `json:"expires"`
}

func CacheGet(ldb *leveldb.DB, key string) ([]byte, error) {

	data, err := ldb.Get([]byte(key), nil)

	if err != nil {
		if err != leveldb_errors.ErrNotFound {
			return nil, err
		}

		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var cache CacheType
	err = json.Unmarshal(data, &cache)

	if err != nil {
		return nil, nil
	}

	secs := time.Now().Unix()

	if cache.Expires > 0 && cache.Expires <= secs {
		return nil, nil
	}

	return cache.Data, nil
}

func CacheSet(ldb *leveldb.DB, key string, value string, expires int64) error {
	cache := CacheType{Data: []byte(value), Created: time.Now().Unix(), Expires: 0}

	if expires > 0 {
		cache.Expires = cache.Created + expires
	}

	json_string, err := json.Marshal(cache)

	err = ldb.Put([]byte(key), []byte(json_string), nil)

	return err
}
