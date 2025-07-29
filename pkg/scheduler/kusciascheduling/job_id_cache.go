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
	"sync"
	"time"

	"github.com/secretflow/kuscia/pkg/utils/nlog"
)

const (
	// Default cleanup interval for expired cache entries
	defaultCleanupInterval = time.Hour

	// Default maximum cache duration (7 days)
	defaultMaxCacheDuration = 7 * 24 * time.Hour
)

// JobIDCache manages job ID to node name mappings with expiration
type JobIDCache struct {
	// Fields reordered for better memory alignment
	maxCacheTime  time.Duration
	cache         map[string]cachedEntry
	mu            sync.RWMutex
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

// cachedEntry holds a node name and its creation timestamp
type cachedEntry struct {
	nodeName     string
	creationTime time.Time
}

// NewJobIDCache creates a new JobIDCache instance
func NewJobIDCache() *JobIDCache {
	cache := &JobIDCache{
		cache:        make(map[string]cachedEntry),
		maxCacheTime: defaultMaxCacheDuration,
		stopChan:     make(chan struct{}),
	}

	// Start cleanup goroutine
	go cache.runCleanupLoop()

	return cache
}

// Set adds or updates a job ID to node mapping
func (jc *JobIDCache) Set(jobID, nodeName string) {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	jc.cache[jobID] = cachedEntry{
		nodeName:     nodeName,
		creationTime: time.Now(),
	}
	nlog.Debugf("Set job ID %s to node %s in cache", jobID, nodeName)
}

// Get retrieves the node name for a given job ID
func (jc *JobIDCache) Get(jobID string) (string, bool) {
	jc.mu.RLock()
	defer jc.mu.RUnlock()

	entry, ok := jc.cache[jobID]
	if !ok {
		return "", false
	}

	// Check if the entry has expired
	if time.Since(entry.creationTime) > jc.maxCacheTime {
		// Entry expired, we'll remove it in the next cleanup cycle
		return "", false
	}

	nlog.Debugf("Get node %s for job ID %s from cache", entry.nodeName, jobID)
	return entry.nodeName, true
}

// Delete removes a job ID from the cache
func (jc *JobIDCache) Delete(jobID string) {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	delete(jc.cache, jobID)
	nlog.Debugf("Delete job ID %s from cache", jobID)
}

// runCleanupLoop periodically cleans up expired entries
func (jc *JobIDCache) runCleanupLoop() {
	jc.cleanupTicker = time.NewTicker(defaultCleanupInterval)
	defer jc.cleanupTicker.Stop()

	for {
		select {
		case <-jc.cleanupTicker.C:
			jc.cleanup()
		case <-jc.stopChan:
			return
		}
	}
}

// cleanup removes expired entries from the cache
func (jc *JobIDCache) cleanup() {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	now := time.Now()

	for jobID, entry := range jc.cache {
		if now.Sub(entry.creationTime) > jc.maxCacheTime {
			delete(jc.cache, jobID)
			nlog.Debugf("Cleaned up expired job ID %s from cache", jobID)
		}
	}
}

// Close stops the cleanup goroutine
func (jc *JobIDCache) Close() {
	close(jc.stopChan)
}
