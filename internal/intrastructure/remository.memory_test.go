package intrastructure

import (
	"context"
	"github.com/ovargas/item-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository(t *testing.T) {

	repository := &ItemMemoryRepository{}

	var id string
	t.Run("create", func(t *testing.T) {
		itemId, err := repository.Create(context.TODO(), &domain.Item{
			Name:        "Dummy Item",
			Description: "A dummy item created for testing",
			ImageUrl:    "http://localhost/images/dummy.jpg",
		})

		assert.NoError(t, err)
		assert.NotNil(t, itemId)
		assert.NotEmpty(t, itemId)

		id = itemId
	})

	t.Run("get", func(t *testing.T) {
		item, err := repository.Get(context.TODO(), id)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, id, item.Id)
		assert.Equal(t, "Dummy Item", item.Name)
		assert.Equal(t, "A dummy item created for testing", item.Description)
		assert.Equal(t, "http://localhost/images/dummy.jpg", item.ImageUrl)
	})

	t.Run("fetch ids", func(t *testing.T) {
		page, err := repository.Fetch(context.TODO(), domain.ItemCriteria{
			Ids:         []string{id},
			Name:        "",
			Description: "",
			Page: domain.Page{
				Number: 1,
				Size:   10,
			},
		})

		assert.NoError(t, err)
		assert.NotNil(t, page)
		assert.NotEmpty(t, page.Items)
		assert.Equal(t, 1, len(page.Items))
	})

	t.Run("fetch name", func(t *testing.T) {
		page, err := repository.Fetch(context.TODO(), domain.ItemCriteria{
			Ids:         nil,
			Name:        "Dummy Item",
			Description: "",
			Page: domain.Page{
				Number: 1,
				Size:   10,
			},
		})

		assert.NoError(t, err)
		assert.NotNil(t, page)
		assert.NotEmpty(t, page.Items)
		assert.Equal(t, 1, len(page.Items))
	})

	t.Run("fetch description", func(t *testing.T) {
		page, err := repository.Fetch(context.TODO(), domain.ItemCriteria{
			Ids:         nil,
			Name:        "",
			Description: "testing",
			Page: domain.Page{
				Number: 1,
				Size:   10,
			},
		})

		assert.NoError(t, err)
		assert.NotNil(t, page)
		assert.NotEmpty(t, page.Items)
		assert.Equal(t, 1, len(page.Items))
	})

	t.Run("update", func(t *testing.T) {
		err := repository.Update(context.TODO(), &domain.Item{
			Id:          id,
			Name:        "Another name",
			Description: "A different description",
			ImageUrl:    "http://localhost/images/dummy.jpg",
		})

		assert.NoError(t, err)

		item, err := repository.Get(context.TODO(), id)

		assert.NoError(t, err)

		assert.Equal(t, "Another name", item.Name)
		assert.Equal(t, "A different description", item.Description)
		assert.Equal(t,  "http://localhost/images/dummy.jpg", item.ImageUrl)
	})

	t.Run("delete", func(t *testing.T) {
		err := repository.Delete(context.TODO(), id)

		assert.NoError(t, err)

		item, err := repository.Get(context.TODO(), id)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Equal(t, domain.ErrItemNotFound, err)
	})

}
