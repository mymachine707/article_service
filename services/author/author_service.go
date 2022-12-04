package author

import (
	"context"
	"fmt"
	"log"
	"mymachine707/models"
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

	var updatedAt string
	if author.UpdatedAt != nil {
		updatedAt = author.UpdatedAt.String()
	}

	var deletedAt string
	if author.DeletedAt != nil {
		deletedAt = author.UpdatedAt.String()
	}

	return &blogpost.Author{
		Id:         author.ID,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt.String(),
		UpdatedAt:  updatedAt,
		DeletedAt:  deletedAt,
	}, nil
}

func (s *authorService) UpdateAuthor(ctx context.Context, req *blogpost.UpdateAuthorRequest) (*blogpost.Author, error) {
	fmt.Println("<<< ---- UpdateAuthor ---->>>")

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

	var updatedAt string
	if author.UpdatedAt != nil {
		updatedAt = author.UpdatedAt.String()
	}

	var deletedAt string
	if author.DeletedAt != nil {
		deletedAt = author.UpdatedAt.String()
	}

	return &blogpost.Author{
		Id:         author.ID,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt.String(),
		UpdatedAt:  updatedAt,
		DeletedAt:  deletedAt,
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
	var updatedAt string
	if author.UpdatedAt != nil {
		updatedAt = author.UpdatedAt.String()
	}

	var deletedAt string
	if author.DeletedAt != nil {
		deletedAt = author.UpdatedAt.String()
	}

	return &blogpost.Author{
		Id:         author.ID,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt.String(),
		UpdatedAt:  updatedAt,
		DeletedAt:  deletedAt,
	}, nil
}

func (s *authorService) GetAuthorList(ctx context.Context, req *blogpost.GetAuthorListRequest) (*blogpost.GetAuthorListResponse, error) {
	fmt.Println("<<< ---- GetAuthorList ---->>>")

	res := &blogpost.GetAuthorListResponse{
		Authors: make([]*blogpost.Author, 0),
	}

	authorList, err := s.stg.GetAuthorList(int(req.Offset), int(req.Limit), req.Search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorList: %s", err)
	}

	for _, v := range authorList {
		var updatedAt string

		if v.UpdatedAt != nil {
			updatedAt = v.UpdatedAt.String()
		}

		var deletedAt string
		if v.DeletedAt != nil {
			deletedAt = v.UpdatedAt.String()
		}

		res.Authors = append(res.Authors, &blogpost.Author{
			Id:         v.ID,
			Firstname:  v.Firstname,
			Lastname:   v.Lastname,
			Middlename: v.Middlename,
			Fullname:   v.Fullname,
			CreatedAt:  v.CreatedAt.String(),
			UpdatedAt:  updatedAt,
			DeletedAt:  deletedAt,
		})
	}

	return res, nil
}
func (s *authorService) GetAuthorById(ctx context.Context, req *blogpost.GetAuthorByIDRequest) (*blogpost.GetAuthorByIDResponse, error) {
	fmt.Println("<<< ---- GetAuthorById ---->>>")

	author, err := s.stg.GetAuthorByID(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorList: %s", err)
	}

	if author.DeletedAt != nil {
		return nil, status.Errorf(codes.NotFound, "Not found author with id: %s", req.Id)
	}

	var updatedAt string
	if author.UpdatedAt != nil {
		updatedAt = author.UpdatedAt.String()
	}

	var deletedAt string
	if author.DeletedAt != nil {
		deletedAt = author.UpdatedAt.String()
	}

	return &blogpost.GetAuthorByIDResponse{
		Id:         author.ID,
		Firstname:  author.Firstname,
		Lastname:   author.Lastname,
		Middlename: author.Middlename,
		Fullname:   author.Fullname,
		CreatedAt:  author.CreatedAt.String(),
		UpdatedAt:  updatedAt,
		DeletedAt:  deletedAt,
	}, nil
}
