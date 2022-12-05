package author

import (
	"context"
	"fmt"
	"log"

	"mymachine707/protogen/blogpost"
	"mymachine707/storage"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authorService struct {
	stg storage.Interfaces
	blogpost.UnimplementedAuthorServiceServer
}

func NewAuthorService(stg storage.Interfaces) *authorService {
	return &authorService{
		stg: stg,
	}
}
func (s *authorService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	fmt.Println("<<< ---- Ping ---->>>")
	log.Println("Ping")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}

func (s *authorService) CreateAuthor(ctx context.Context, req *blogpost.CreateAuthorRequest) (*blogpost.Author, error) {
	fmt.Println("<<< ---- CreateAuthor ---->>>")

	id := uuid.New()

	err := s.stg.AddAuthor(id.String(), req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.AddAuthor: %s", err)
	}

	author, err := s.stg.GetAuthorByID(id.String()) // maqsad tekshirish rostan  ham create bo'ldimi?
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorByID: %s", err)
	}

	return &blogpost.Author{
		Id:         author.Id,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt,
		UpdatedAt:  author.UpdatedAt,
		DeletedAt:  author.DeletedAt,
	}, nil
}

func (s *authorService) UpdateAuthor(ctx context.Context, req *blogpost.UpdateAuthorRequest) (*blogpost.Author, error) {
	fmt.Println("<<< ---- UpdateAuthor ---->>>")

	err := s.stg.UpdateAuthor(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateAuthor: %s", err)
	}

	author, err := s.stg.GetAuthorByID(req.Id) // maqsad tekshirish rostan  ham create bo'ldimi?
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorByID: %s", err)
	}

	return &blogpost.Author{
		Id:         author.Id,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt,
		UpdatedAt:  author.UpdatedAt,
		DeletedAt:  author.DeletedAt,
	}, nil
}

func (s *authorService) DeleteAuthor(ctx context.Context, req *blogpost.DeleteAuthorRequest) (*blogpost.Author, error) {
	fmt.Println("<<< ---- DeleteAuthor ---->>>")

	author, err := s.stg.GetAuthorByID(req.Id) // maqsad tekshirish rostan  ham create bo'ldimi?
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorByID: %s", err)
	}

	err = s.stg.DeleteAuthor(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteAuthor: %s", err)
	}

	return &blogpost.Author{
		Id:         author.Id,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt,
		UpdatedAt:  author.UpdatedAt,
		DeletedAt:  author.DeletedAt,
	}, nil
}

func (s *authorService) GetAuthorList(ctx context.Context, req *blogpost.GetAuthorListRequest) (*blogpost.GetAuthorListResponse, error) {
	fmt.Println("<<< ---- GetAuthorList ---->>>")

	res, err := s.stg.GetAuthorList(int(req.Offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorList: %s", err)
	}

	return res, nil
}

func (s *authorService) GetAuthorById(ctx context.Context, req *blogpost.GetAuthorByIDRequest) (*blogpost.GetAuthorByIDResponse, error) {
	fmt.Println("<<< ---- GetAuthorById ---->>>")

	author, err := s.stg.GetAuthorByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorByID: %s", err)
	}

	return author, nil
}
