package data

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"time"
)

type File struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	Category  string    `json:"category"`
	ClientID  int       `json:"client_id"`
}

type FileModel struct {
	DB *sql.DB
}

func (m *FileModel) GetFilesMetadata(form *multipart.Form) ([]File, error) {
	if form == nil || len(form.File) == 0 {
		return nil, errors.New("no files provided")
	}

	var filesMetadata []File

	fileCategories := []string{"purchase_order", "invoice", "handing_over_report", "pms_report"}
	for _, category := range fileCategories {
		if _, ok := form.File[category]; !ok {
			continue
		}

		for _, fileHeader := range form.File[category] {
			fileName := fmt.Sprintf("%v_%s", time.Now().Format("2006-01-02_3:04_PM"), fileHeader.Filename)
			filePath := fmt.Sprintf("/files/%s/%s", category, fileName)
			filesMetadata = append(filesMetadata, File{
				FileName: fileName,
				FilePath: filePath,
				Category: category,
			})
		}
	}

	return filesMetadata, nil
}
