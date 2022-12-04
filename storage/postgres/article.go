package postgres

import (
	"errors"

	// "mymachine707/blogpost"
	"mymachine707/protogen/blogpost"
)

// AddArticle ...
func (stg Postgres) AddArticle(id string, entity *blogpost.CreateArticleRequest) error {
	if id == "" {
		return errors.New("id must exist")
	}

	if entity.Content == nil {
		entity.Content = &blogpost.Content{}
	}

	_, err := stg.GetAuthorByID(entity.AuthorId)

	if err != nil {
		return err
	}

	_, err = stg.db.Exec(`INSERT INTO article (
		id,
		title,
		body,
		author_id		
		) VALUES (
		$1,
		$2,
		$3,
		$4
)`,
		id,
		entity.Content.Title,
		entity.Content.Body,
		entity.AuthorId,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetArticleByID ...  //  ????
func (stg Postgres) GetArticleByID(id string) (*blogpost.GetArticleByIDResponse, error) {
	// var res blogpost.GetArticleByIDResponse
	res := &blogpost.GetArticleByIDResponse{
		Content: &blogpost.Content{},
		Author:  &blogpost.GetArticleByIDResponse_Author{},
	}

	var tempMiddlename *string

	if id == "" {
		return res, errors.New("id must exist")
	}

	err := stg.db.QueryRow(`SELECT 
    ar.id,
    ar.title,
    ar.body,
    ar.created_at,
    ar.updated_at,
    ar.deleted_at,
    au.id,
    au.firstname,
    au.lastname,
	au.middlename,
	au.fullname,
    au.created_at,
    au.updated_at,
    au.deleted_at FROM article AS ar JOIN author AS au ON ar.author_id = au.id WHERE ar.id = $1`, id).Scan(
		&res.Id,
		&res.Content.Title,
		&res.Content.Body,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
		&res.Author.Id,
		&res.Author.Firstname,
		&res.Author.Lastname,
		&tempMiddlename,
		&res.Author.Fullname,
		&res.Author.CreatedAt,
		&res.Author.UpdatedAt,
		&res.Author.DeletedAt,
	)

	if tempMiddlename != nil {
		res.Author.Middlename = *tempMiddlename
	}

	if err != nil {
		return res, err
	}

	if res.DeletedAt != "" {
		return res, errors.New("article not found")
	}

	return res, err
}

// GetArticleList ...
func (stg Postgres) GetArticleList(offset, limit int, search string) (*blogpost.GetArticleListResponse, error) {

	resp := &blogpost.GetArticleListResponse{
		Articles: make([]*blogpost.Article, 0),
	}

	rows, err := stg.db.Queryx(`
	
	SELECT * FROM article WHERE 
	
	(title || ' ' || body ILIKE '%' || $1 || '%') AND deleted_at is null
	LIMIT $2
	OFFSET $3
	
	`, search, limit, offset)

	if err != nil {
		return resp, err
	}

	for rows.Next() {
		a := &blogpost.Article{
			Content: &blogpost.Content{},
		}

		err := rows.Scan(
			&a.Id,
			&a.Content.Title,
			&a.Content.Body,
			&a.AuthorId,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.DeletedAt,
		)

		//err := rows.StructScan(&a)
		if err != nil {
			return resp, err
		}

		//fmt.Printf("%d a---> %#v\n", i, a)
		resp.Articles = append(resp.Articles, a) // ??? &  ?
	}

	return resp, err
}

// UpdateArticle ...
func (stg Postgres) UpdateArticle(article *blogpost.UpdateArticleRequest) error {

	if article.Content == nil {
		article.Content = &blogpost.Content{}
	}

	res, err := stg.db.NamedExec("UPDATE article SET title=:t, body=:b, updated_at=now() WHERE id=:id AND deleted_at is null", map[string]interface{}{
		"id": article.Id,
		"t":  article.Content.Title,
		"b":  article.Content.Body,
	})

	if err != nil {
		return err
	}

	n, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("article not found")
}

// DeleteArticle ...
func (stg Postgres) DeleteArticle(idStr string) error {

	res, err := stg.db.Exec("UPDATE article Set deleted_at=now() WHERE id=$1 AND deleted_at is null", idStr)

	if err != nil {
		return err
	}

	n, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("article not found")
}
