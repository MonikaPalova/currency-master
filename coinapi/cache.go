package coinapi

import (
	"time"
)

const minutesToKeepCache = 5

// Cache object which keeps information about assets received from external api
type Cache struct {
	assets  []Asset
	ids     map[string]int
	expires time.Time
}

// Cache constructor
func NewCache() *Cache {
	return &Cache{assets: []Asset{}, ids: make(map[string]int), expires: time.Now()}
}

// clears cache and adds assets
func (c *Cache) Fill(assets []Asset) {
	c.assets = assets
	for idx, asset := range assets {
		c.ids[asset.ID] = idx
	}
	c.expires = time.Now().Add(time.Minute * minutesToKeepCache)
}

// gets specific page from cache
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

// gets asset by id
func (c Cache) GetAsset(id string) *Asset {
	if c.IsExpired() {
		c.cleanCache()
		return nil
	}
	return &c.assets[c.ids[id]]
}

// returns if cache is expired
func (c Cache) IsExpired() bool {
	return c.expires.Before(time.Now())
}

func (c *Cache) cleanCache() {
	c.assets = []Asset{}
}
