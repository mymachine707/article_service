package article

import (
	"context"
	"fmt"
	"log"
	"mymachine707/models"
	blogpost "mymachine707/protogen/blogpost"
	"mymachine707/storage"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type articleService struct {
	stg storage.Interfaces
	blogpost.UnimplementedArticleServiceServer
}

func NewArticleService(stg storage.Interfaces) *articleService {
	return &articleService{
		stg: stg,
	}
}

func (s *articleService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Ping")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}

func (s *articleService) CreateArticle(ctx context.Context, req *blogpost.CreateArticleRequest) (*blogpost.Article, error) {
	fmt.Println("<<< ---- CreateArticle ---->>>")
	// create new article
	id := uuid.New()
	err := s.stg.AddArticle(id.String(), models.CreateArticleModul{
		Content: models.Content{
			Title: req.Content.Title,
			Body:  req.Content.Body,
		},
		AuthorID: req.AuthorId,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.AddArticle: %s", err.Error())
	}

	article, err := s.stg.GetArticleByID(id.String()) // maqsad tekshirish rostan  ham create bo'ldimi?

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleByID: %s", err.Error())
	}

	var updatedAt string
	if article.UpdatedAt != nil {
		updatedAt = article.UpdatedAt.String()
	}

	var deletedAt string
	if article.DeletedAt != nil {
		deletedAt = article.DeletedAt.String()
	}
	return &blogpost.Article{
		Id: article.ID,
		Content: &blogpost.Content{
			Title: req.Content.Title,
			Body:  req.Content.Body,
		},
		AuthorId:  article.Author.ID,
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}, nil
}
func (s *articleService) UpdateArticle(ctx context.Context, req *blogpost.UpdateArticleRequest) (*blogpost.Article, error) {
	fmt.Println("<<< ---- UpdateArticle ---->>>")
	err := s.stg.UpdateArticle(models.UpdateArticleModul{
		ID: req.Id,
		Content: models.Content{
			Title: req.Content.Title,
			Body:  req.Content.Body,
		},
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateArticle: %s", err.Error())
	}

	article, err := s.stg.GetArticleByID(req.Id) // maqsad tekshirish rostan  ham create bo'ldimi?

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleByID---!: %s", err.Error())
	}

	var updatedAt string
	if article.UpdatedAt != nil {
		updatedAt = article.UpdatedAt.String()
	}

	var deletedAt string
	if article.DeletedAt != nil {
		deletedAt = article.DeletedAt.String()
	}

	return &blogpost.Article{
		Id: article.ID,
		Content: &blogpost.Content{
			Title: req.Content.Title,
			Body:  req.Content.Body,
		},
		AuthorId:  article.Author.ID,
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func (s *articleService) DeleteArticle(ctx context.Context, req *blogpost.DeleteArticleRequest) (*blogpost.Article, error) {
	fmt.Println("<<< ---- DeleteArticle ---->>>")
	err := s.stg.DeleteArticle(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteArticle: %s", err.Error())
	}

	article, err := s.stg.GetArticleByID(req.Id) // maqsad tekshirish rostan  ham create bo'ldimi?

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleByID: %s", err.Error())
	}

	var updatedAt string
	if article.UpdatedAt != nil {
		updatedAt = article.UpdatedAt.String()
	}

	var deletedAt string
	if article.DeletedAt != nil {
		deletedAt = article.DeletedAt.String()
	}

	return &blogpost.Article{
		Id: article.ID,
		Content: &blogpost.Content{
			Title: article.Content.Title,
			Body:  article.Content.Body,
		},
		AuthorId:  article.Author.ID,
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func (s *articleService) GetArticleList(ctx context.Context, req *blogpost.GetArticleListRequest) (*blogpost.GetArticleListResponse, error) {
	fmt.Println("<<< ---- GetArticleList ---->>>")
	res := &blogpost.GetArticleListResponse{
		Articles: make([]*blogpost.Article, 0), // Article list nil bo'masligi uchun
	}

	articleList, err := s.stg.GetArticleList(int(req.Offset), int(req.Limit), req.Search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleList: %s", err.Error())
	}

	for _, v := range articleList {
		var updatedAt string
		if v.UpdatedAt != nil {
			updatedAt = v.UpdatedAt.String()
		}

		var deletedAt string
		if v.DeletedAt != nil {
			deletedAt = v.DeletedAt.String()
		}
		res.Articles = append(res.Articles, &blogpost.Article{
			Id: v.ID,
			Content: &blogpost.Content{
				Title: v.Content.Title,
				Body:  v.Content.Body,
			},
			AuthorId:  v.AuthorID,
			CreatedAt: v.CreatedAt.String(),
			UpdatedAt: updatedAt,
			DeletedAt: deletedAt,
		})

	}

	return res, nil
}
func (s *articleService) GetArticleById(ctx context.Context, req *blogpost.GetArticleByIDRequest) (*blogpost.GetArticleByIDResponse, error) {
	fmt.Println("<<< ---- GetArticleById ---->>>")

	article, err := s.stg.GetArticleByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleByID: %s", err.Error())
	}

	if article.DeletedAt != nil {
		return nil, status.Errorf(codes.NotFound, "Not found article with id: %s", req.Id)
	}

	var authorDeletedAt string
	if article.Author.DeletedAt != nil {
		authorDeletedAt = article.Author.DeletedAt.String()
	}

	var authorUpdatedAt string
	if article.Author.UpdatedAt != nil {
		authorUpdatedAt = article.Author.UpdatedAt.String()
	}
	var updatedAt string
	if article.UpdatedAt != nil {
		updatedAt = article.UpdatedAt.String()
	}

	var deletedAt string
	if article.DeletedAt != nil {
		deletedAt = article.DeletedAt.String()
	}

	// !!!
	return &blogpost.GetArticleByIDResponse{
		Id: article.ID,
		Content: &blogpost.Content{
			Title: article.Title,
			Body:  article.Body,
		},
		Author: &blogpost.GetArticleByIDResponse_Author{
			Id:         article.Author.ID,
			Firstname:  article.Author.Firstname,
			Lastname:   article.Author.Lastname,
			Middlename: article.Author.Middlename,
			Fullname:   article.Author.Fullname,
			CreatedAt:  article.Author.CreatedAt.String(),
			UpdatedAt:  authorUpdatedAt,
			DeletedAt:  authorDeletedAt,
		},
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}, nil
}
