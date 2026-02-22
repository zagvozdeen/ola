package api

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Service) UploadFile(r *http.Request, user *models.User) core.Response {
	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, err)
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		return core.Err(http.StatusInternalServerError, err)
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return core.Err(http.StatusInternalServerError, err)
	}
	ct := http.DetectContentType(b)
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return core.Err(http.StatusInternalServerError, err)
	}
	//err = os.Mkdir(".data/files", os.ModePerm)
	//if err != nil && !errors.Is(err, os.ErrExist) {
	//	return core.Err(http.StatusInternalServerError, err)
	//}
	err = os.WriteFile(fmt.Sprintf(".data/files/%s", uid.String()), b, os.ModePerm)
	if err != nil {
		return core.Err(http.StatusInternalServerError, err)
	}
	f := &models.File{
		UUID:       uid,
		Content:    fmt.Sprintf("/files/%s", uid.String()),
		Size:       header.Size,
		MimeType:   mt,
		OriginName: header.Filename,
		UserID:     user.ID,
		CreatedAt:  time.Now(),
	}
	err = s.store.CreateFile(r.Context(), f)
	if err != nil {
		return core.Err(http.StatusInternalServerError, err)
	}
	return core.JSON(http.StatusCreated, f)
}
