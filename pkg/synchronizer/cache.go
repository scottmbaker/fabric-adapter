// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

// Cache implements a cache of data that is pushed to the core.

package synchronizer

import (
	"fmt"
	"reflect"
)

const (
	// CacheModelNetConfig is the modelName to use when caching fabric to the onos
	CacheModelNetConfig = "netconfig"
)

// CacheCheck returns true if (modelName, modelId) exists in the cache and the contents have not
// changed.
func (s *Synchronizer) CacheCheck(modelName string, modelID string, contents interface{}) bool {
	key := fmt.Sprintf("%s-%s", modelName, modelID)
	entry, okay := s.cache[key]
	if !okay {
		return false
	}

	return reflect.DeepEqual(entry, contents) // (entry == contents)
}

// CacheUpdate updates the contents of (modelName, modelID) in the cache with new contents
func (s *Synchronizer) CacheUpdate(modelName string, modelID string, contents interface{}) {
	key := fmt.Sprintf("%s-%s", modelName, modelID)
	s.cache[key] = contents
}

// CacheInvalidate removes all entries in the cache
func (s *Synchronizer) CacheInvalidate() {
	s.cache = map[string]interface{}{}
}

// CacheDelete removes a single entry from the cache
func (s *Synchronizer) CacheDelete(modelName string, modelID string) {
	key := fmt.Sprintf("%s-%s", modelName, modelID)

	// delete does not crash if the key does not exist
	delete(s.cache, key)
}