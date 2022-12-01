package article

import (
	"context"
	"log"
	blogpost "mymachine707/protogen/blogpost"
)

type ArticleService struct {
	blogpost.UnimplementedArticleServiceServer
}

func (s *ArticleService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Ping")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}
