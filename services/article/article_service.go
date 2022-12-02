package article

import (
	"context"
	"log"
	"mymachine707/models"
	blogpost "mymachine707/protogen/blogpost"
	"mymachine707/storage"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ArticleService struct {
	stg storage.Interfaces
	blogpost.UnimplementedArticleServiceServer
}

func NewArticleService(stg storage.Interfaces) *ArticleService {
	return &ArticleService{
		stg: stg,
	}
}

func (s *ArticleService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Ping")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}

func (s *ArticleService) CreateArticle(ctx context.Context, req *blogpost.CreateArticleRequest) (*blogpost.Article, error) {

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
		return nil, status.Errorf(codes.Internal, "s.stg.AddArticle: %s", err)
	}

	article, err := s.stg.GetArticleByID(id.String()) // maqsad tekshirish rostan  ham create bo'ldimi?

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleByID: %s", err)
	}

	return &blogpost.Article{
		Id: article.ID,
		Content: &blogpost.Content{
			Title: req.Content.Title,
			Body:  req.Content.Body,
		},
		AuthorId:  article.Author.ID,
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: article.UpdatedAt.String(),
	}, nil
}
func (s *ArticleService) UpdateArticle(ctx context.Context, req *blogpost.UpdateArticleRequest) (*blogpost.Article, error) {

	err := s.stg.UpdateArticle(models.UpdateArticleModul{
		ID: req.Id,
		Content: models.Content{
			Title: req.Content.Title,
			Body:  req.Content.Body,
		},
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateArticle: %s", err)
	}

	article, err := s.stg.GetArticleByID(req.Id) // maqsad tekshirish rostan  ham create bo'ldimi?

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleByID: %s", err)
	}

	return &blogpost.Article{
		Id: article.ID,
		Content: &blogpost.Content{
			Title: req.Content.Title,
			Body:  req.Content.Body,
		},
		AuthorId:  article.Author.ID,
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: article.UpdatedAt.String(),
	}, nil
}
func (s *ArticleService) DeleteArticle(ctx context.Context, req *blogpost.DeleteArticleRequest) (*blogpost.Article, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteArticle not implemented")
}
func (s *ArticleService) GetArticleList(ctx context.Context, req *blogpost.GetArticleListRequest) (*blogpost.GetArticleListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleList not implemented")
}
func (s *ArticleService) GetArticleById(ctx context.Context, req *blogpost.GetArticleByIDRequest) (*blogpost.GetArticleByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleById not implemented")
}
