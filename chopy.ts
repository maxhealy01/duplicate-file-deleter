package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var filesSeen = make(map[string][32]byte)
var duplicates []DuplicateFile
var totalFreedSize int64

var maxDuplicates = 100
var ErrMaxDuplicatesFound = errors.New("maximum number of duplicates found")

type DuplicateFile struct {
	Name    string
	Path    string
	Size    int64
	SizeStr string
	Kind    string
}

func calculateSHA256(filePath string) ([32]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return [32]byte{}, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return [32]byte{}, err
	}

	return sha256.Sum256(hash.Sum(nil)), nil
}

func visit(root string, path string, f os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("Error accessing a path %q: %v\n", path, err)
		return err
	}

	if !f.IsDir() {
		hash, err := calculateSHA256(path)
		if err != nil {
			fmt.Printf("Failed to calculate hash for file %q: %v\n", path, err)
			return nil
		}

		if prevHash, seen := filesSeen[f.Name()]; seen {
			if prevHash == hash {
				relativePath := strings.TrimPrefix(path, root)
				formattedSize := formatSize(f.Size())
				kind := filepath.Ext(path)
				duplicates = append(duplicates, DuplicateFile{Name: f.Name(), Path: relativePath, Size: f.Size(), SizeStr: formattedSize, Kind: kind})

				if len(duplicates) >= 10 {
					if shouldDelete(duplicates) {
						for _, duplicate := range duplicates {
							err := os.Remove(root + duplicate.Path)
							if err != nil {
								fmt.Printf("Failed to remove file %q: %v\n", root+duplicate.Path, err)
							} else {
								totalFreedSize += duplicate.Size
								fmt.Println("File deleted:", duplicate.Name, "at", duplicate.Path)
							}
						}
					}
					duplicates = nil
				}

				if len(duplicates) >= maxDuplicates {
					return ErrMaxDuplicatesFound
				}
			}
		} else {
			filesSeen[f.Name()] = hash
		}
	}
	return nil
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func shouldDelete(files []DuplicateFile) bool {
	totalSize := int64(0)
	for _, file := range files {
		totalSize += file.Size
		fmt.Printf("\nName: %s\nPath: %s\nSize: %s\nKind: %s\n", file.Name, file.Path, file.SizeStr, file.Kind)
	}

	fmt.Printf("\nDelete these files to save %s of space on your computer? [Y/n]: ", formatSize(totalSize))
	var answer string
	_, err := fmt.Scanln(&answer)
	if err != nil {
		fmt.Printf("Error reading your response: %v\n", err)
		return false
	}

	return strings.ToUpper(answer) == "Y"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a path as a command line argument in the format 'username/subdirectory'")
		os.Exit(1)
	}

	path := os.Args[1]
	root := fmt.Sprintf("/Users/%s", path)
	err := filepath.Walk(root, func(p string, f os.FileInfo, err error) error {
		return visit(root, p, f, err)
	})

	if err != nil {
		if errors.Is(err, ErrMaxDuplicatesFound) {
			fmt.Println("Stopped after finding 100 duplicates.")
		} else {
			fmt.Printf("error walking the path %q: %v\n", root, err)
		}
	}

	fmt.Printf("\n%s freed!\n", formatSize(totalFreedSize))
}
