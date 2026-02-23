package store

import (
	"context"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) CreateFile(ctx context.Context, file *models.File) error {
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO files (uuid, content, size, mime_type, origin_name, user_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		file.UUID, file.Content, file.Size, file.MimeType, file.OriginName, file.UserID, file.CreatedAt,
	).Scan(&file.ID)
	return wrapDBError(err)
}
