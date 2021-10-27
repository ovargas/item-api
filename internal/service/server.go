package service

import (
	"context"
	"github.com/ovargas/api-go/item/v1"
	"github.com/ovargas/api-go/storage/v1"
	"github.com/ovargas/item-api/internal/domain"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net/url"
)

type ItemService struct {
	item.UnimplementedItemServiceServer
	repository    Repository
	storageClient storage.StorageServiceClient
}

type Repository interface {
	Get(ctx context.Context, id string) (*domain.Item, error)
	Fetch(ctx context.Context, criteria domain.ItemCriteria) (*domain.ItemPage, error)
	Create(ctx context.Context, item *domain.Item) (string, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, id string) error
}

func New(repository Repository, storageClient storage.StorageServiceClient) *ItemService {
	return &ItemService{
		repository:    repository,
		storageClient: storageClient,
	}
}

func (s *ItemService) Get(ctx context.Context, request *item.GetRequest) (*item.Item, error) {
	item, err := s.repository.Get(ctx, request.GetId())
	if err != nil {
		return nil, err
	}
	return item.ToPB(), nil
}

func (s *ItemService) Fetch(ctx context.Context, request *item.FetchRequest) (*item.FetchResponse, error) {
	page, err := s.repository.Fetch(ctx, domain.ItemCriteria{
		Ids:         request.GetIds(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		Page: domain.Page{
			Number: request.GetPage().GetNumber(),
			Size:   request.GetPage().GetSize(),
		},
	})

	if err != nil {
		return nil, err
	}

	return &item.FetchResponse{
		Content:      page.Items.ToPB(),
		TotalRecords: page.TotalRecords,
	}, nil
}

func (s *ItemService) Create(ctx context.Context, request *item.CreateRequest) (*item.Item, error) {

	buffer, _ := proto.Marshal(request)

	log.Default().Printf(url.QueryEscape(string(buffer)))

	file, err := s.storageClient.Create(ctx, &storage.CreateRequest{
		Filename:  request.GetImage().GetName(),
		MediaType: request.GetImage().GetMediaType(),
		Bytes:     request.GetImage().GetBytes(),
	})

	if err != nil {
		return nil, err
	}

	item := &domain.Item{
		Name:        request.GetName(),
		Description: request.GetDescription(),
		ImageUrl:    file.GetUrl(),
	}

	id, err := s.repository.Create(ctx, item)

	if err != nil {
		return nil, err
	}

	item.Id = id
	return item.ToPB(), nil
}

func (s *ItemService) Update(ctx context.Context, request *item.UpdateRequest) (*emptypb.Empty, error) {
	item, err := s.repository.Get(ctx, request.GetId())

	if err != nil {
		return nil, err
	}

	if request.GetName() != "" {
		item.Name = request.GetName()
	}

	if request.GetDescription() != "" {
		item.Description = request.GetDescription()
	}

	err = s.repository.Update(ctx, item)

	return &emptypb.Empty{}, err
}
func (s *ItemService) Delete(ctx context.Context, request *item.DeleteRequest) (*emptypb.Empty, error) {
	err := s.repository.Delete(ctx, request.GetId())
	return &emptypb.Empty{}, err
}
