package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/ASH-WIN-10/uniref-app-backend/internal/data"
)

func (app *application) CalculateFilesMetadata(form *multipart.Form, clientID int) []data.File {
	if form == nil || len(form.File) == 0 {
		return []data.File{}
	}

	dirPath := filepath.Join("assets", "files", fmt.Sprintf("%d", clientID))

	var filesMetadata []data.File
	fileCategories := []string{"purchase_order", "invoice", "handing_over_report", "pms_report"}
	for _, category := range fileCategories {
		if _, ok := form.File[category]; !ok {
			continue
		}

		for _, file := range form.File[category] {
			fileName := fmt.Sprintf("%v_%s_%s", time.Now().Format("2006-01-02_3:04_PM"), category, file.Filename)
			filePath := filepath.Join(dirPath, fileName)
			filesMetadata = append(filesMetadata, data.File{
				OriginalFileName: file.Filename,
				FileName:         fileName,
				FilePath:         filePath,
				Category:         category,
				ClientID:         clientID,
			})
		}
	}

	return filesMetadata
}

func (app *application) SaveFilesLocally(form *multipart.Form, filesMetadata []data.File) error {
	if len(filesMetadata) == 0 {
		return nil
	}

	dirPath := filepath.Dir(filesMetadata[0].FilePath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	for _, f := range filesMetadata {
		for _, file := range form.File[f.Category] {
			if file.Filename != f.OriginalFileName {
				continue
			}

			src, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer src.Close()

			dst, err := os.Create(f.FilePath)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				return fmt.Errorf("failed to copy file: %w", err)
			}
		}
	}

	return nil
}
