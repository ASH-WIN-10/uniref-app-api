package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/ASH-WIN-10/Himwan-Refrigerations-Backend/internal/data"
)

func (app *application) SaveFilesLocally(form *multipart.Form, clientID int) ([]data.File, error) {
	if form == nil || len(form.File) == 0 {
		app.logger.Info("No files provided")
		return []data.File{}, nil
	}

	dirPath := filepath.Join("assets", "files", fmt.Sprintf("%d", clientID))

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	var filesMetadata []data.File

	fileCategories := []string{"purchase_order", "invoice", "handing_over_report", "pms_report"}
	for _, category := range fileCategories {
		if _, ok := form.File[category]; !ok {
			continue
		}

		for _, file := range form.File[category] {
			src, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open file: %w", err)
			}
			defer src.Close()

			fileName := fmt.Sprintf("%v_%s_%s", time.Now().Format("2006-01-02_3:04_PM"), category, file.Filename)
			filePath := filepath.Join(dirPath, fileName)
			dst, err := os.Create(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to create file: %w", err)
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				return nil, fmt.Errorf("failed to copy file: %w", err)
			}

			filesMetadata = append(filesMetadata, data.File{
				FileName: fileName,
				FilePath: filePath,
				Category: category,
				ClientID: clientID,
			})
		}
	}

	return filesMetadata, nil
}
