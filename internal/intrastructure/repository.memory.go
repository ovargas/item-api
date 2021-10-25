package intrastructure

import (
	"context"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/google/uuid"
	"github.com/ovargas/item-api/internal/domain"
	"strings"
	"sync"
)

var (
	mu    sync.RWMutex
	items = make(map[string]*domain.Item)
)

type ItemMemoryRepository struct {
}

func NewItemMemoryRepository() *ItemMemoryRepository {
	return &ItemMemoryRepository{}
}

func (r *ItemMemoryRepository) Get(ctx context.Context, id string) (*domain.Item, error) {
	mu.RLock()
	defer mu.RUnlock()
	result, ok := items[id]
	if !ok {
		return nil, domain.ErrItemNotFound
	}

	return result, nil
}

func (r *ItemMemoryRepository) Fetch(ctx context.Context, criteria domain.ItemCriteria) (*domain.ItemPage, error) {
	mu.RLock()
	defer mu.RUnlock()

	query := From(items).
		Where(func(i interface{}) bool {
			keyValue := i.(KeyValue)
			item := keyValue.Value.(*domain.Item)
			return containsIds(item, criteria.Ids) &&
				withName(item, criteria.Name) &&
				containsDescription(item, criteria.Description)
		})

	tr := query.Count()
	var list domain.ItemList

	query.Skip(int((criteria.Page.Number - 1) * criteria.Size)).
		Take(int(criteria.Size)).
		Select(func(i interface{}) interface{} {
			return i.(KeyValue).Value.(*domain.Item)
		}).
		ToSlice(&list)
	return &domain.ItemPage{
		Items:        list,
		TotalRecords: int32(tr),
	}, nil
}

func (r *ItemMemoryRepository) Create(ctx context.Context, item *domain.Item) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	if item.Id == "" {
		item.Id = uuid.New().String()
	}

	items[item.Id] = item

	return item.Id, nil

}

func (r *ItemMemoryRepository) Update(ctx context.Context, item *domain.Item) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := items[item.Id]; !ok {
		return domain.ErrItemNotFound
	}

	items[item.Id] = item

	return nil
}

func (r *ItemMemoryRepository) Delete(ctx context.Context, id string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := items[id]; !ok {
		return domain.ErrItemNotFound
	}

	delete(items, id)
	return nil
}

func containsIds(item *domain.Item, ids []string) bool {
	if len(ids) == 0 {
		return true
	}
	return From(ids).Contains(item.Id)
}

func withName(item *domain.Item, name string) bool {
	if name == "" {
		return true
	}
	return item.Name == name
}

func containsDescription(item *domain.Item, description string) bool {
	return strings.Index(item.Description, description) > -1
}
