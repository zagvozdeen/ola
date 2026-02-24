package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetFileByUUID(ctx context.Context, uuid uuid.UUID) (*models.File, error) {
	file := &models.File{}
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT id, uuid, content, size, mime_type, origin_name, user_id, created_at FROM files WHERE uuid = $1",
		uuid,
	).Scan(&file.ID, &file.UUID, &file.Content, &file.Size, &file.MimeType, &file.OriginName, &file.UserID, &file.CreatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return file, nil
}

func (s *Store) CreateFile(ctx context.Context, file *models.File) error {
	err := s.querier(ctx).QueryRow(
		ctx,
		"INSERT INTO files (uuid, content, size, mime_type, origin_name, user_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		file.UUID, file.Content, file.Size, file.MimeType, file.OriginName, file.UserID, file.CreatedAt,
	).Scan(&file.ID)
	return wrapDBError(err)
}
