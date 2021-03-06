package coinapi

import (
	"time"
)

const minutesToKeepCache = 30

// Cache object which keeps information about assets received from external api.
type Cache struct {
	assets  []Asset
	ids     map[string]int
	expires time.Time
}

// Cache constructor.
func NewCache() *Cache {
	return &Cache{assets: []Asset{}, ids: map[string]int{}, expires: time.Now().Add(-time.Hour)}
}

// Clears cache and adds assets.
func (c *Cache) Fill(assets []Asset) {
	c.assets = assets
	for idx, asset := range assets {
		c.ids[asset.ID] = idx
	}
	c.expires = time.Now().Add(time.Minute * minutesToKeepCache)
}

// Gets specific page from cache.
// If page is negative or after the last page, returns an empty page
func (c Cache) GetPage(page, size int) AssetPage {
	if c.IsExpired() {
		c.cleanCache()
		return AssetPage{[]Asset{}, page, size, 0}
	}

	from := (page - 1) * size
	to := from + size
	total := len(c.assets)
	if from > total || from < 0 {
		return AssetPage{[]Asset{}, page, size, total}
	}
	if to > total {
		to = total
	}
	return AssetPage{c.assets[from:to], page, size, total}
}

// Gets asset by id.
// If cache is expired or asset is not in cache, returns nil
func (c Cache) GetAsset(id string) *Asset {
	if c.IsExpired() {
		c.cleanCache()
		return nil
	}
	pos, ok := c.ids[id]
	if !ok {
		return nil
	}
	return &c.assets[pos]
}

// Returns if cache is expired.
func (c Cache) IsExpired() bool {
	return c.expires.Before(time.Now())
}

func (c *Cache) cleanCache() {
	c.assets = []Asset{}
	c.ids = map[string]int{}
	c.expires = time.Now()
}
