package article

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
	err := s.stg.AddArticle(id.String(), req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.AddArticle: %s", err.Error())
	}

	article, err := s.stg.GetArticleByID(id.String()) // maqsad tekshirish rostan  ham create bo'ldimi?

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleByID: %s", err.Error())
	}

	return &blogpost.Article{
		Id:        article.Id,
		Content:   article.Content,
		AuthorId:  article.Author.Id,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		DeletedAt: article.DeletedAt,
	}, nil
}

func (s *articleService) UpdateArticle(ctx context.Context, req *blogpost.UpdateArticleRequest) (*blogpost.Article, error) {
	fmt.Println("<<< ---- UpdateArticle ---->>>")
	err := s.stg.UpdateArticle(req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateArticle: %s", err.Error())
	}

	article, err := s.stg.GetArticleByID(req.Id) // maqsad tekshirish rostan  ham create bo'ldimi?

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleByID---!: %s", err.Error())
	}

	return &blogpost.Article{
		Id:        article.Id,
		Content:   article.Content,
		AuthorId:  article.Author.Id,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		DeletedAt: article.DeletedAt,
	}, nil
}

func (s *articleService) DeleteArticle(ctx context.Context, req *blogpost.DeleteArticleRequest) (*blogpost.Article, error) {
	fmt.Println("<<< ---- DeleteArticle ---- >>>")

	article, err := s.stg.GetArticleByID(req.Id) // maqsad tekshirish rostan  ham create bo'ldimi?

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleByID: %s", err.Error())
	}

	err = s.stg.DeleteArticle(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteArticle: %s", err.Error())
	}

	return &blogpost.Article{
		Id:        article.Id,
		Content:   article.Content,
		AuthorId:  article.Author.Id,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		DeletedAt: article.DeletedAt,
	}, nil
}

func (s *articleService) GetArticleList(ctx context.Context, req *blogpost.GetArticleListRequest) (*blogpost.GetArticleListResponse, error) {
	fmt.Println("<<< ---- GetArticleList ---->>>")

	res, err := s.stg.GetArticleList(int(req.Offset), int(req.Limit), req.Search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleList: %s", err.Error())
	}

	return res, nil
}
func (s *articleService) GetArticleById(ctx context.Context, req *blogpost.GetArticleByIDRequest) (*blogpost.GetArticleByIDResponse, error) {
	fmt.Println("<<< ---- GetArticleById ---->>>")

	article, err := s.stg.GetArticleByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleByID: %s", err.Error())
	}

	return article, nil
}
