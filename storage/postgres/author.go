package postgres

import (
	"errors"
	"mymachine707/protogen/blogpost"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var err error

// AddAuthor ...
func (stg Postgres) AddAuthor(id string, entity *blogpost.CreateAuthorRequest) error {
	if id == "" {
		return errors.New("id must exist")
	}
	fname := entity.Firstname + " " + entity.Lastname + " " + entity.Middlename

	_, err = stg.db.Exec(`INSERT INTO author (
		id,
		firstname,
		lastname,
		middlename,
		fullname
		) VALUES(
		$1,
		$2,
		$3,
		$4,
		$5
	)`,
		id,
		entity.Firstname,
		entity.Lastname,
		entity.Middlename,
		fname,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetAuthorByID ...
func (stg Postgres) GetAuthorByID(id string) (*blogpost.Author, error) {
	result := &blogpost.Author{}

	var createdAt, updatedAt, deletedAt *time.Time
	err := stg.db.QueryRow(`SELECT
		id,
		firstname,
		lastname,
		middlename,
		fullname,
		created_at,
		updated_at,
		deleted_at
	FROM author WHERE id=$1`, id).Scan(
		&result.Id,
		&result.Firstname,
		&result.Lastname,
		&result.Middlename,
		&result.Fullname,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)

	if err != nil {
		return result, err
	}

	if createdAt != nil {
		result.CreatedAt = timestamppb.New(*createdAt)
	}

	if updatedAt != nil {
		result.UpdatedAt = timestamppb.New(*updatedAt)
	}

	if deletedAt != nil {
		result.DeletedAt = timestamppb.New(*deletedAt)
	}

	return result, nil
}

// GetAuthorList ...
func (stg Postgres) GetAuthorList(offset, limit int, search string) (resp *blogpost.GetAuthorListResponse, err error) {

	resp = &blogpost.GetAuthorListResponse{
		Authors: make([]*blogpost.Author, 0),
	}

	rows, err := stg.db.Queryx(`
	
	Select * from author WHERE ((firstname ILIKE '%' || $1 || '%') OR 
		(lastname ILIKE '%' || $1 || '%') OR 
		(middlename ILIKE '%' || $1 || '%') OR 
		(fullname ILIKE '%' || $1 || '%'))
		AND deleted_at is null 
		LIMIT $2 
		OFFSET $3`, search, limit, offset)

	if err != nil {
		return resp, err
	}
	var updatedAt, deletedAt *time.Time

	for rows.Next() {
		var a *blogpost.Author

		//var updatedAt, deletedAt *string
		err = rows.Scan(
			&a.Id,
			&a.Firstname,
			&a.Lastname,
			&a.Middlename,
			&a.Fullname,
			&a.CreatedAt,
			&updatedAt,
			&deletedAt,
		)

		if err != nil {
			return resp, err
		}

		// if createdAt != nil {
		// 	a.CreatedAt = timestamppb.New(*createdAt)
		// }
		if updatedAt != nil {
			a.UpdatedAt = timestamppb.New(*updatedAt)
		}
		if deletedAt != nil {
			a.DeletedAt = timestamppb.New(*deletedAt)
		}

		resp.Authors = append(resp.Authors, a)

	}

	return resp, nil
}

// UpdateAuthor ...
func (stg Postgres) UpdateAuthor(author *blogpost.UpdateAuthorRequest) error {

	fname := author.Firstname + " " + author.Lastname + " " + author.Middlename

	rows, err := stg.db.NamedExec(`Update author set firstname=:f, lastname=:l, middlename=:m, fullname=:fn,updated_at=now() Where id=:id  and deleted_at is null`, map[string]interface{}{
		"id": author.Id,
		"f":  author.Firstname,
		"l":  author.Lastname,
		"m":  author.Middlename,
		"fn": fname,
	})

	if err != nil {
		return err
	}

	n, err := rows.RowsAffected()

	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("author not found")
}

// DeleteAuthor ...
func (stg Postgres) DeleteAuthor(idStr string) error {

	rows, err := stg.db.Exec(`UPDATE author SET deleted_at=now() Where id=$1 and deleted_at is null`, idStr)

	if err != nil {
		return err
	}

	n, err := rows.RowsAffected()

	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("Cannot delete Author becouse Author not found")
}

// hard delete uchun kod
// func (stg Postgres) removeAuthorDelete(slice []blogpost.Author, s int) []blogpost.Author {
// 	return append(slice[:s], slice[s+1:]...)
// }
