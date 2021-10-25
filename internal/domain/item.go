package domain

import (
	"github.com/ovargas/api-go/item/v1"
)

type Item struct {
	Id          string
	Name        string
	Description string
	ImageUrl    string
}

func (i *Item) Equals(item *Item) bool {
	return i.Id == item.Id
}

func (i *Item) ToPB() *item.Item {
	return &item.Item{
		Id:          i.Id,
		Name:        i.Name,
		Description: i.Description,
		ImageUrl:    i.ImageUrl,
	}
}

type ItemCriteria struct {
	Ids         []string
	Name        string
	Description string
	Page
}

type ItemList []*Item

type ItemPage struct {
	Items        ItemList
	TotalRecords int32
}

func (il ItemList) ToPB() []*item.Item {
	var items []*item.Item
	for _, i := range il {
		items = append(items, i.ToPB())
	}
	return items
}
