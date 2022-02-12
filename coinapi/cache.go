package coinapi

import (
	"time"
)

const MINUTES_TO_KEEP_CACHE = 5

type Cache struct {
	assets  []Asset
	expires time.Time
}

func NewCache() *Cache {
	return &Cache{[]Asset{}, time.Now()}
}

func (c *Cache) Fill(assets []Asset) {
	c.assets = assets
	c.expires = time.Now().Add(time.Minute * MINUTES_TO_KEEP_CACHE)
}

func (c Cache) GetPage(page, size int) AssetPage {
	if c.isExpired() {
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

func (c Cache) isExpired() bool {
	return c.expires.Before(time.Now())
}

func (c *Cache) cleanCache() {
	c.assets = []Asset{}
}
