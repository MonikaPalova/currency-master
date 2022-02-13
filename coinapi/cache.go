package coinapi

import (
	"time"
)

const minutesToKeepCache = 5

type Cache struct {
	assets  []Asset
	ids     map[string]*Asset
	expires time.Time
}

func NewCache() *Cache {
	return &Cache{assets: []Asset{}, ids: make(map[string]*Asset), expires: time.Now()}
}

func (c *Cache) Fill(assets []Asset) {
	c.assets = assets
	for _, asset := range assets {
		c.ids[asset.ID] = &asset
	}
	c.expires = time.Now().Add(time.Minute * minutesToKeepCache)
}

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

func (c Cache) GetAsset(id string) *Asset {
	if c.IsExpired() {
		c.cleanCache()
		return nil
	}

	return c.ids[id]
}

func (c Cache) IsExpired() bool {
	return c.expires.Before(time.Now())
}

func (c *Cache) cleanCache() {
	c.assets = []Asset{}
}
