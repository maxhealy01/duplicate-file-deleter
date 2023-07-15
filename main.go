package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

type File struct {
	Name string
	Size int64
	Kind string
	Hash string
	Path string
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter the username and subdirectory (format username/subdirectory): ")
    text, _ := reader.ReadString('\n')
    text = strings.TrimSpace(text)
    userDir := "/Users/" + text

    fmt.Print("Enter the minimum file size to consider for deletion (in kilobytes): ")
    sizeInput, _ := reader.ReadString('\n')
    sizeInput = strings.TrimSpace(sizeInput)
    minSizeKB, err := strconv.ParseInt(sizeInput, 10, 64)
    if err != nil {
        fmt.Println("Invalid size input, defaulting to 0")
        minSizeKB = 0
    }
    minSize := minSizeKB * 1024 // convert KB to bytes

    fileMap := make(map[string]File)
    duplicates := make([]File, 0)
    deleteChan := make(chan []File)

    go func() {
        for files := range deleteChan {
            printAndDeleteBatch(files)
        }
    }()

    _ = filepath.Walk(userDir, visit(fileMap, &duplicates, deleteChan, minSize))

    // Send remaining duplicates if there are any
    if len(duplicates) > 0 {
        deleteChan <- duplicates
    }
    close(deleteChan)
}

func printAndDeleteBatch(files []File) {
	var totalSize int64 = 0
	for _, file := range files {
		totalSize += file.Size
	}

	fmt.Println("Found duplicates:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Name\tSize (KB)\tKind\tPath\t")
	for _, file := range files {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t\n", file.Name, file.Size/1024, file.Kind, file.Path)
	}
	w.Flush()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Delete these files to save", formatSize(totalSize), "of space on your computer")
	fmt.Print("Are you sure you want to delete these files? (Y/N): ")
	text, _ := reader.ReadString('\n')
	if strings.TrimSpace(text) == "Y" {
		for _, file := range files {
			os.Remove(file.Path)
		}
		fmt.Printf("\033[32mDeleted %d files, freed %s of space!\033[0m\n", len(files), formatSize(totalSize))
	}
}


func visit(fileMap map[string]File, duplicates *[]File, deleteChan chan []File, minSize int64) filepath.WalkFunc {
    return func(path string, f os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !f.IsDir() {
            size := f.Size()
            if size < minSize {
                return nil
            }
            kind := getKind(path)
            name := f.Name()
            hash, err := calculateHash(path)
            if err != nil {
                return nil
            }
            file := File{Name: name, Size: size, Kind: kind, Hash: hash, Path: path}

            if val, ok := fileMap[hash]; ok {
                if file.Size == val.Size {
                    *duplicates = append(*duplicates, file)
                    if len(*duplicates) >= 15 {
                        deleteChan <- *duplicates
                        *duplicates = nil
                    }
                }
            } else {
                fileMap[hash] = file
            }
        }
        return nil
    }
}

func calculateHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	hash := hex.EncodeToString(h.Sum(nil))
	return hash, nil
}

func formatSize(size int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB"}
	i := 0
	fSize := float64(size)
	for fSize >= 1024 && i < len(sizes)-1 {
		fSize /= 1024
		i++
	}
	return fmt.Sprintf("%.2f %s", fSize, sizes[i])
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
