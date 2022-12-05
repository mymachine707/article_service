package storage

import (
	"mymachine707/protogen/blogpost"
)

// Interfaces ...
type Interfaces interface {
	// For Article
	AddArticle(id string, entity *blogpost.CreateArticleRequest) error
	GetArticleByID(id string) (*blogpost.GetArticleByIDResponse, error)
	GetArticleList(offset, limit int, search string) (resp *blogpost.GetArticleListResponse, err error)
	UpdateArticle(article *blogpost.UpdateArticleRequest) error
	DeleteArticle(idStr string) error

	// For Author
	AddAuthor(id string, entity *blogpost.CreateAuthorRequest) error
	GetAuthorByID(id string) (*blogpost.Author, error)
	GetAuthorList(offset, limit int, serach string) (resp *blogpost.GetAuthorListResponse, err error)
	UpdateAuthor(author *blogpost.UpdateAuthorRequest) error
	DeleteAuthor(idStr string) error
}
