package apiclient

import (
	"fmt"
	"sync"
	"time"
)

// We maintain an in-memory cache cache of rendered dashboards, to reduce rendering load
// from multiple clients. The memory footprint is about 300KB for a typical dashboard
// rendered on 1920x1080 screen size. The pods should be scaled accordingly depending on
// how many dashboards on how many resolutions are used by clients.

type cacheEntry struct {
	autoFitPanel bool
	timeRendered time.Time
	payload      []byte
}

var (
	cacheMutex      = sync.RWMutex{}
	dashboardCaches = map[string]cacheEntry{}
)

func cacheKey(dashboardURL string, height, width int, autoFitPanel bool) string {
	return fmt.Sprintf("%s:%d:%d:%t", dashboardURL, height, width, autoFitPanel)
}

func cacheRender(dashboardURL string, height, width int, autoFitPanel bool, dashboardPayload []byte) cacheEntry {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	entry := cacheEntry{
		timeRendered: time.Now(),
		payload:      dashboardPayload,
	}

	dashboardCaches[cacheKey(dashboardURL, height, width, autoFitPanel)] = entry
	return entry
}

func getCachedRender(dashboardURL string, height, width int, autoFitPanel bool) (bool, cacheEntry) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	entry, ok := dashboardCaches[cacheKey(dashboardURL, height, width, autoFitPanel)]
	if !ok {
		// Cache miss, needs to re-render
		return false, cacheEntry{}
	}
	if time.Now().Sub(entry.timeRendered) > cacheValidity {
		// Cache stale, needs to re-render. The stale cache will be overwritten when re-rendered
		return false, cacheEntry{}
	}

	// Cache hit and valid
	return true, entry
}
