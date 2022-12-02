package author

import (
	"context"
	"log"
	"mymachine707/protogen/blogpost"
)

type AuthorService struct {
	blogpost.UnimplementedAuthorServiceServer
}

func (s *AuthorService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Ping")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}
