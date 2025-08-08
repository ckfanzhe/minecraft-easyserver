package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractZip extracts zip file to target directory
func ExtractZip(src, dest string) error {
	// Open zip file for reading
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()

	// Create destination directory
	os.MkdirAll(dest, 0755)

	// Extract files
	for _, f := range r.File {
		// Skip problematic system files that may cause permission issues
		fileName := filepath.Base(f.Name)
		if strings.ToLower(fileName) == "desktop.ini" || strings.ToLower(fileName) == "thumbs.db" {
			continue
		}

		// Construct file path
		path := filepath.Join(dest, f.Name)

		// Check path security to prevent directory traversal attacks
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			// Create directory with safe permissions
			os.MkdirAll(path, 0755)
			continue
		}

		// Create file directory
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// Extract file
		fileReader, err := f.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		// Use safe file permissions instead of preserving original permissions
		// This prevents permission issues with files from different operating systems
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}

	return nil
}