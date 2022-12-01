package author

import (
	"context"
	"log"
	authorProto "mymachine707/protogen/blogpost"
)

type AuthorService struct {
	authorProto.UnimplementedAuthorServiceServer
}

func (s *AuthorService) Ping(ctx context.Context, req *authorProto.Empty) (*authorProto.Pong, error) {
	log.Println("Ping")
	return &authorProto.Pong{
		Message: "Ok",
	}, nil
}
