package item_api

import (
	"context"
	"github.com/google/uuid"

	"github.com/ovargas/api-go/commons/v1"
	"github.com/ovargas/api-go/item/v1"
	"github.com/ovargas/api-go/storage/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestAll(t *testing.T) {

	conn, err := grpc.Dial("localhost:10001", grpc.WithInsecure())
	defer func(storageConnection *grpc.ClientConn) {
		_ = storageConnection.Close()
	}(conn)

	if err != nil {
		log.Fatalln(err)
	}

	itemClient := item.NewItemServiceClient(conn)

	itemFullTest(t, itemClient)
}

func itemFullTest(t *testing.T, itemClient item.ItemServiceClient) {
	var itemId string

	fileName := uuid.New().String() + ".txt"

	t.Run("create", func(t *testing.T) {
		itemCreated, err := itemClient.Create(context.TODO(), &item.CreateRequest{
			Name:        "The Item",
			Description: "A Description",
			Image: &storage.File{
				Name:      fileName,
				MediaType: "plain/text",
				Content: &storage.File_Bytes{
					Bytes: []byte("Hello World!!!"),
				},
			},
		},
		)

		assert.NoError(t, err)
		assert.NotNil(t, itemCreated)
		assert.Equal(t, "The Item", itemCreated.Name)
		assert.Equal(t, "A Description", itemCreated.Description)
		assert.Equal(t, fileName, itemCreated.ImageUrl)

		itemId = itemCreated.Id
	})

	t.Run("get", func(t *testing.T) {
		itemReturned, err := itemClient.Get(context.TODO(), &item.GetRequest{Id: itemId})
		assert.NoError(t, err)
		assert.NotNil(t, itemReturned)
		assert.Equal(t, "The Item", itemReturned.Name)
		assert.Equal(t, "A Description", itemReturned.Description)
		assert.Equal(t, fileName, itemReturned.ImageUrl)
	})

	t.Run("fetch", func(t *testing.T) {
		itemReturned, err := itemClient.Fetch(context.TODO(), &item.FetchRequest{
			Description: "Description",
			Page: &commons.Page{
				Number: 1,
				Size:   10,
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, itemReturned.Content)
	})
}

func BenchmarkItem(b *testing.B) {

	conn, err := grpc.Dial("localhost:10001", grpc.WithInsecure())

	defer func(storageConnection *grpc.ClientConn) {
		_ = storageConnection.Close()
	}(conn)

	if err != nil {
		log.Fatalln(err)
	}

	itemClient := item.NewItemServiceClient(conn)

	for i := 0; i < b.N; i++ {
		fileName := uuid.New().String() + ".txt"
		_, _ = itemClient.Create(context.TODO(), &item.CreateRequest{
			Name:        "The Item",
			Description: "A Description",
			Image: &storage.File{
				Name:      fileName,
				MediaType: "plain/text",
				Content: &storage.File_Bytes{
					Bytes: []byte("Hello World!!!"),
				},
			},
		})
	}
}
