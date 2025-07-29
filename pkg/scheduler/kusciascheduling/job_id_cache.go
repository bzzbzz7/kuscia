// Copyright 2025 Ant Group Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kusciascheduling

import (
	"container/list"
	"sync"
	"time"

	"github.com/secretflow/kuscia/pkg/utils/nlog"
)

const (
	// Default maximum cache duration (10 days)
	defaultMaxCacheDuration = 10 * 24 * time.Hour

	// Threshold for cleanup check
	cleanupThreshold = 1000
)

// JobIDCache manages job ID to node name mappings with expiration and insertion order
type JobIDCache struct {
	maxCacheTime time.Duration
	cache        map[string]*list.Element // map to list element for O(1) access
	orderList    *list.List               // maintains insertion order
	mu           sync.Mutex
}

// cachedEntry holds a node name and its creation timestamp
type cachedEntry struct {
	jobID        string
	nodeName     string
	creationTime time.Time
}

// NewJobIDCache creates a new JobIDCache instance
func NewJobIDCache() *JobIDCache {
	return &JobIDCache{
		cache:        make(map[string]*list.Element),
		orderList:    list.New(),
		maxCacheTime: defaultMaxCacheDuration,
	}
}

// Set adds or updates a job ID to node mapping
// Also cleans up expired entries from the front of the list when size exceeds threshold
func (jc *JobIDCache) Set(jobID, nodeName string) {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	now := time.Now()

	entry := &cachedEntry{
		jobID:        jobID,
		nodeName:     nodeName,
		creationTime: now,
	}

	// If already exists, update it and move to back
	if elem, ok := jc.cache[jobID]; ok {
		elem.Value = entry
		jc.orderList.MoveToBack(elem)
		nlog.Debugf("Updated job ID %s to node %s in cache", jobID, nodeName)
		return
	}

	// Otherwise, insert new
	elem := jc.orderList.PushBack(entry)
	jc.cache[jobID] = elem
	nlog.Debugf("Inserted job ID %s to node %s in cache", jobID, nodeName)

	// Check if we need to cleanup expired entries
	if jc.orderList.Len() > cleanupThreshold {
		jc.cleanupExpired(now)
	}
}

// Get retrieves the node name for a given job ID
func (jc *JobIDCache) Get(jobID string) (string, bool) {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	// Try to get the requested entry
	elem, ok := jc.cache[jobID]
	if !ok {
		return "", false
	}

	entry := elem.Value.(*cachedEntry)

	nlog.Debugf("Get node %s for job ID %s from cache", entry.nodeName, jobID)
	return entry.nodeName, true
}

// Delete removes a job ID from the cache
func (jc *JobIDCache) Delete(jobID string) {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	if elem, ok := jc.cache[jobID]; ok {
		delete(jc.cache, jobID)
		jc.orderList.Remove(elem)
		nlog.Debugf("Deleted job ID %s from cache", jobID)
	}
}

// cleanupExpired removes expired entries from the front of the list
// Only called when cache size exceeds cleanupThreshold
func (jc *JobIDCache) cleanupExpired(now time.Time) {
	for jc.orderList.Len() > 0 {
		front := jc.orderList.Front()
		entry := front.Value.(*cachedEntry)

		// Check if the oldest entry is expired
		if now.Sub(entry.creationTime) <= jc.maxCacheTime {
			// Oldest entry is not expired, stop cleaning
			break
		}

		// Remove expired entry
		delete(jc.cache, entry.jobID)
		jc.orderList.Remove(front)
		nlog.Debugf("Removed expired job ID %s from cache during Set", entry.jobID)

		// Continue to check next oldest entry
	}
}
