package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/iamvkosarev/wishlist/back/internal/model"
	"github.com/iamvkosarev/wishlist/back/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

const prepareStorageStatement = `
CREATE TABLE IF NOT EXISTS wishlists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL CHECK (length(name) <= 100),
    description TEXT CHECK (length(description) <= 500),
    owner_id INTEGER NOT NULL,
    display_type INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_wishlists_owner ON wishlists(owner_id);

CREATE TABLE IF NOT EXISTS wishes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wishlist_id INTEGER NOT NULL,
    name TEXT NOT NULL CHECK (length(name) <= 100),
    description TEXT CHECK (length(description) <= 500),
    wish_url TEXT,
    image_url TEXT,
    assigned_to_id INTEGER,
    FOREIGN KEY(wishlist_id) REFERENCES wishlists(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_wishes_wishlist_id ON wishes(wishlist_id);;
`

func New(address string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", address)
	if err != nil {
		return nil, fmt.Errorf("%s address: %v error: %w", op, address, err)
	}
	stmt, err := db.Prepare(prepareStorageStatement)
	if err != nil {
		return nil, fmt.Errorf("%s address: %v error: %w", op, address, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) GetWishlist(wishlistID int64) (model.Wishlist, error) {
	const op = "storage.sqlite.GetWishlist"

	stmt, err := s.db.Prepare("SELECT * FROM wishlists WHERE id = ?")
	if err != nil {
		return model.Wishlist{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	var wishlist model.Wishlist
	err = stmt.QueryRow(wishlistID).Scan(
		&wishlist.ID, &wishlist.Name, &wishlist.Description, &wishlist.OwnerID, &wishlist.DisplayType,
	)
	if err != nil {
		return model.Wishlist{}, fmt.Errorf("%s: %w", op, err)
	}
	return wishlist, nil
}

func (s *Storage) SaveWishlist(ownerID int64, name string, description string, displayType model.DisplayType) (
	int64,
	error,
) {
	const op = "storage.sqlite.SaveWishlist"

	stmt, err := s.db.Prepare("INSERT INTO wishlists(owner_id, name, description, display_type) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(ownerID, name, description, displayType)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrorWishlistExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetWishlists(owner int64) ([]model.Wishlist, error) {
	const op = "storage.sqlite.GetAllWishlists"

	rows, err := s.db.Query(
		"SELECT id, name, description, owner_id, display_type FROM wishlists WHERE owner_id = ?",
		owner,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var wishlists []model.Wishlist
	for rows.Next() {
		var w model.Wishlist
		if err := rows.Scan(&w.ID, &w.Name, &w.Description, &w.OwnerID, &w.DisplayType); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		wishlists = append(wishlists, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return wishlists, nil
}
