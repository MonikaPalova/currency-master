package coinapi

import (
	"time"
)

const minutesToKeepCache = 5

type cache struct {
	assets  []Asset
	ids     map[string]*Asset
	expires time.Time
}

func newCache() *cache {
	return &cache{assets: []Asset{}, ids: make(map[string]*Asset), expires: time.Now()}
}

func (c *cache) fill(assets []Asset) {
	c.assets = assets
	for _, asset := range assets {
		c.ids[asset.ID] = &asset
	}
	c.expires = time.Now().Add(time.Minute * minutesToKeepCache)
}

func (c cache) getPage(page, size int) AssetPage {
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

func (c cache) getAsset(id string) *Asset {
	if c.isExpired() {
		c.cleanCache()
		return nil
	}

	return c.ids[id]
}

func (c cache) isExpired() bool {
	return c.expires.Before(time.Now())
}

func (c *cache) cleanCache() {
	c.assets = []Asset{}
}
