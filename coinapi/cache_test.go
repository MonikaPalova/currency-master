package coinapi

import (
	"reflect"
	"testing"
	"time"
)

func TestCache_Fill(t *testing.T) {
	type want struct {
		assets []Asset
		ids    map[string]int
	}
	type args struct {
		assets []Asset
	}
	tests := []struct {
		name string
		want want
		args args
	}{
		{"empty cache", want{[]Asset{}, map[string]int{}}, args{[]Asset{}}},
		{"fill cache", want{[]Asset{{ID: "id1"}, {ID: "id2"}}, map[string]int{"id1": 0, "id2": 1}}, args{[]Asset{{ID: "id1"}, {ID: "id2"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache()

			c.Fill(tt.args.assets)

			if c.IsExpired() {
				t.Error("Cache should not be expired immediately after being refilled")
			}
			if !reflect.DeepEqual(tt.want.assets, c.assets) {
				t.Errorf("The cache assets were not filled properly. Got %v, wanted %v", c.assets, tt.want.assets)
			}
			if !reflect.DeepEqual(tt.want.ids, c.ids) {
				t.Errorf("The cache ids were not filled properly. Got %v, wanted %v", c.ids, tt.want.ids)
			}
		})
	}
}

func TestCache_GetPage(t *testing.T) {
	type fields struct {
		assets  []Asset
		ids     map[string]int
		expires time.Time
	}
	type args struct {
		page int
		size int
	}
	a1 := Asset{ID: "id1"}
	a2 := Asset{ID: "id2"}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AssetPage
	}{
		{"page 0", fields{[]Asset{a1, a2}, map[string]int{a1.ID: 0, a2.ID: 1}, time.Now().Add(time.Hour * 1)}, args{0, 1}, AssetPage{[]Asset{}, 0, 1, 2}},
		{"page after last", fields{[]Asset{a1, a2}, map[string]int{a1.ID: 0, a2.ID: 1}, time.Now().Add(time.Hour * 1)}, args{15, 1}, AssetPage{[]Asset{}, 15, 1, 2}},
		{"page not full", fields{[]Asset{a1, a2}, map[string]int{a1.ID: 0, a2.ID: 1}, time.Now().Add(time.Hour * 1)}, args{1, 3}, AssetPage{[]Asset{a1, a2}, 1, 3, 2}},
		{"size 0", fields{[]Asset{a1, a2}, map[string]int{a1.ID: 0, a2.ID: 1}, time.Now().Add(time.Hour * 1)}, args{1, 0}, AssetPage{[]Asset{}, 1, 0, 2}},
		{"page 1 size 1", fields{[]Asset{a1, a2}, map[string]int{a1.ID: 0, a2.ID: 1}, time.Now().Add(time.Hour * 1)}, args{1, 1}, AssetPage{[]Asset{a1}, 1, 1, 2}},
		{"page 2 size 1", fields{[]Asset{a1, a2}, map[string]int{a1.ID: 0, a2.ID: 1}, time.Now().Add(time.Hour * 1)}, args{2, 1}, AssetPage{[]Asset{a2}, 2, 1, 2}},
		{"expired", fields{[]Asset{a1, a2}, map[string]int{a1.ID: 0, a2.ID: 1}, time.Now().Add(-time.Hour * 1)}, args{1, 1}, AssetPage{[]Asset{}, 1, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				assets:  tt.fields.assets,
				ids:     tt.fields.ids,
				expires: tt.fields.expires,
			}
			if got := c.GetPage(tt.args.page, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.GetPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_GetAsset(t *testing.T) {
	type fields struct {
		assets  []Asset
		ids     map[string]int
		expires time.Time
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Asset
	}{
		{"get existing asset", fields{[]Asset{{ID: "id1"}}, map[string]int{"id1": 0}, time.Now().Add(time.Hour * 1)}, args{"id1"}, &Asset{ID: "id1"}},
		{"get not existing asset", fields{[]Asset{{ID: "id1"}}, map[string]int{"id1": 0}, time.Now().Add(time.Hour * 1)}, args{"id2"}, nil},
		{"get not existing asset empty cache", fields{[]Asset{}, map[string]int{}, time.Now().Add(time.Hour * 1)}, args{"id1"}, nil},
		{"expired cache", fields{[]Asset{{ID: "id1"}}, map[string]int{"id1": 0}, time.Now().Add(-time.Hour * 1)}, args{"id1"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				assets:  tt.fields.assets,
				ids:     tt.fields.ids,
				expires: tt.fields.expires,
			}
			if got := c.GetAsset(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.GetAsset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_IsExpired(t *testing.T) {
	type fields struct {
		expires time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"true", fields{time.Now().Add(-1 * time.Hour)}, true},
		{"false", fields{time.Now().Add(1 * time.Hour)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				expires: tt.fields.expires,
			}
			if got := c.IsExpired(); got != tt.want {
				t.Errorf("Cache.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
