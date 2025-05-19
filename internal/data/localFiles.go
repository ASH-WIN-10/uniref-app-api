package data

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveFileLocally(fileHeader *multipart.FileHeader, file *File) error {
	dirPath := filepath.Join("assets", "files", fmt.Sprintf("%d", file.ClientID))
	fileName := fmt.Sprintf("%v_%s_%s", time.Now().Format("2006-01-02_3:04_PM"), file.Category, fileHeader.Filename)
	filePath := filepath.Join(dirPath, fileName)

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	file.FileName = fileName
	file.FilePath = filePath
	file.OriginalFileName = fileHeader.Filename

	return nil
}

func CalculateFilesMetadata(form *multipart.Form, clientID int) []File {
	if form == nil || len(form.File) == 0 {
		return []File{}
	}

	dirPath := filepath.Join("assets", "files", fmt.Sprintf("%d", clientID))

	var filesMetadata []File
	fileCategories := []string{"purchase_order", "invoice", "handing_over_report", "pms_report"}
	for _, category := range fileCategories {
		if _, ok := form.File[category]; !ok {
			continue
		}

		for _, file := range form.File[category] {
			fileName := fmt.Sprintf("%v_%s_%s", time.Now().Format("2006-01-02_3:04_PM"), category, file.Filename)
			filePath := filepath.Join(dirPath, fileName)
			filesMetadata = append(filesMetadata, File{
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

func SaveFilesLocally(form *multipart.Form, filesMetadata []File) error {
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

func DeleteFiles(newFiles, oldFiles []File) ([]File, error) {
	if len(oldFiles) == 0 {
		return nil, nil
	}

	deletedFiles := []File{}
	for _, oldFile := range oldFiles {
		if oldFile.FilePath == "" {
			continue
		}

		found := false
		for _, newFile := range newFiles {
			if oldFile.FilePath == newFile.FilePath {
				found = true
				break
			}
		}

		if !found {
			err := os.Remove(oldFile.FilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to delete file: %w", err)
			}
			deletedFiles = append(deletedFiles, oldFile)
		}
	}

	return deletedFiles, nil
}

func GetNewlyAddedFiles(oldFiles, newFiles []File) []File {
	if len(newFiles) == 0 {
		return nil
	}

	newlyAddedFiles := []File{}
	for _, newFile := range newFiles {
		found := false
		for _, oldFile := range oldFiles {
			if newFile.FilePath == oldFile.FilePath {
				found = true
				break
			}
		}

		if !found {
			newlyAddedFiles = append(newlyAddedFiles, newFile)
		}
	}

	return newlyAddedFiles
}

func DeleteFileLocally(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
