package main

import "database/sql"

var globalImageStore ImageStore

const pageSize = 25

// ImageStore describes an interface for a type that can handle persistence and retrieval of image details
type ImageStore interface {
	Save(image *Image) error
	Find(id string) (*Image, error)
	FindAll(offset int) ([]Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
}

// DBImageStore handles storing and retrieving of image details
type DBImageStore struct {
	db *sql.DB
}

// NewDBImageStore creates a new instance of a database image store
func NewDBImageStore() ImageStore {
	return &DBImageStore{
		db: globalMySQLDB,
	}
}

// Save persists an image into the specified DBImageStore
func (store *DBImageStore) Save(image *Image) error {
	query := `REPLACE INTO images
        (id, user_id, name, location, description, size, created_at)
        VALUES
            (?, ?, ?, ?, ?, ?, ?)
    `
	_, err := store.db.Exec(query,
		image.ID,
		image.UserID,
		image.Name,
		image.Location,
		image.Description,
		image.Size,
		image.CreatedAt,
	)

	return err
}

// Find gets an image with the specified ID from the DBImageStore
func (store *DBImageStore) Find(id string) (*Image, error) {
	query := `SELECT id, user_id, name, location, description, size, created_at
        FROM images
        WHERE id = ?
    `
	row := store.db.QueryRow(query, id)

	image := Image{}

	err := row.Scan(
		&image.ID,
		&image.UserID,
		&image.Name,
		&image.Location,
		&image.Description,
		&image.Size,
		&image.CreatedAt,
	)

	return &image, err
}

// FindAll retrieves a paginated collection of images from the DBImageStore
func (store *DBImageStore) FindAll(offset int) ([]Image, error) {
	query := `SELECT id, user_id, name, location, description, size, created_at
        FROM images
        ORDER BY created_at DESC
        LIMIT ?
        OFFSET ?
    `

	rows, err := store.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}

	images := []Image{}
	for rows.Next() {
		image := Image{}
		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.Name,
			&image.Location,
			&image.Description,
			&image.Size,
			&image.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

// FindAllByUser retrieves a paginated collection of a user's images from the DBImageStore
func (store *DBImageStore) FindAllByUser(user *User, offset int) ([]Image, error) {
	query := `SELECT id, user_id, name, location, description, size, created_at
        FROM images
        WHERE user_id = ?
        ORDER BY created_at DESC
        LIMIT ?
        OFFSET ?
    `

	rows, err := store.db.Query(query, user.ID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	images := []Image{}
	for rows.Next() {
		image := Image{}
		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.Name,
			&image.Location,
			&image.Description,
			&image.Size,
			&image.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}
