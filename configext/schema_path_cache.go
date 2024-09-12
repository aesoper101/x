package configext

import (
	"crypto/sha256"
	"fmt"
	"github.com/aesoper101/x/jsonschemaext"
	"github.com/dgraph-io/ristretto"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

var schemaPathCacheConfig = &ristretto.Config{
	// Hold up to 25 schemas in cache. Usually we only need one.
	MaxCost:            250,
	NumCounters:        2500,
	BufferItems:        64,
	Metrics:            false,
	IgnoreInternalCost: true,
}

var schemaPathCache, _ = ristretto.NewCache(schemaPathCacheConfig)

func getSchemaPaths(rawSchema []byte, schema *jsonschema.Schema) ([]jsonschemaext.Path, error) {
	key := fmt.Sprintf("%x", sha256.Sum256(rawSchema))
	if val, found := schemaPathCache.Get(key); found {
		if validator, ok := val.([]jsonschemaext.Path); ok {
			return validator, nil
		}
		schemaPathCache.Del(key)
	}

	keys, err := jsonschemaext.ListPathsWithInitializedSchemaAndArraysIncluded(schema)
	if err != nil {
		return nil, err
	}

	schemaPathCache.Set(key, keys, 1)
	schemaPathCache.Wait()
	return keys, nil
}
