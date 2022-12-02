package author

import (
	"context"
	"log"
	"mymachine707/models"
	"mymachine707/protogen/blogpost"
	"mymachine707/storage"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthorService struct {
	stg storage.Interfaces
	blogpost.UnimplementedAuthorServiceServer
}

func (s *AuthorService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Ping")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}

func (s *AuthorService) CreateAuthor(ctx context.Context, req *blogpost.CreateAuthorRequest) (*blogpost.Author, error) {
	id := uuid.New()
	err := s.stg.AddAuthor(id.String(), models.CreateAuthorModul{
		Firstname:  req.Firstname,
		Lastname:   req.Lastname,
		Middlename: req.Middlename,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.AddAuthor: %s", err)
	}
	author, err := s.stg.GetAuthorByID(id.String()) // maqsad tekshirish rostan  ham create bo'ldimi?
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorByID: %s", err)
	}

	return &blogpost.Author{
		Id:         author.ID,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt.String(),
		UpdatedAt:  author.UpdatedAt.String(),
		DeletedAt:  author.DeletedAt.String(),
	}, nil
}

func (s *AuthorService) UpdateAuthor(ctx context.Context, req *blogpost.UpdateAuthorRequest) (*blogpost.Author, error) {
	err := s.stg.UpdateAuthor(models.UpdateAuthorModul{
		ID:         req.Id,
		Firstname:  req.Firstname,
		Lastname:   req.Lastname,
		Middlename: req.Middlename,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateAuthor: %s", err)
	}

	author, err := s.stg.GetAuthorByID(req.Id) // maqsad tekshirish rostan  ham create bo'ldimi?
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorByID: %s", err)
	}

	return &blogpost.Author{
		Id:         author.ID,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt.String(),
		UpdatedAt:  author.UpdatedAt.String(),
		DeletedAt:  author.DeletedAt.String(),
	}, nil
}

func (s *AuthorService) DeleteAuthor(ctx context.Context, req *blogpost.DeleteAuthorRequest) (*blogpost.Author, error) {

	author, err := s.stg.GetAuthorByID(req.Id) // maqsad tekshirish rostan  ham create bo'ldimi?
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorByID: %s", err)
	}

	err = s.stg.DeleteAuthor(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteAuthor: %s", err)
	}

	return &blogpost.Author{
		Id:         author.ID,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt.String(),
		UpdatedAt:  author.UpdatedAt.String(),
		DeletedAt:  author.DeletedAt.String(),
	}, nil
}

func (s *AuthorService) GetAuthorList(ctx context.Context, req *blogpost.GetAuthorListRequest) (*blogpost.GetAuthorListResponse, error) {
	authorList, err := s.stg.GetAuthorList(int(req.Offset), int(req.Limit), req.Search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorList: %s", err)
	}

	return authorList, nil
}
func (s *AuthorService) GetAuthorById(ctx context.Context, req *blogpost.GetAuthorByIDRequest) (*blogpost.GetAuthorByIDResponse, error) {

	author, err := s.stg.GetAuthorByID(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorList: %s", err)
	}

	return &blogpost.GetAuthorByIDResponse{
		Id:         author.ID,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt.String(),
		UpdatedAt:  author.UpdatedAt.String(),
		DeletedAt:  author.DeletedAt.String(),
	}, nil
}
