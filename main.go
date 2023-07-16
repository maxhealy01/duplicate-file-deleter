package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

var (
	freedSpace     int64
	startTime      time.Time
	modulesDeleted int64
)

type File struct {
	Name string
	Size int64
	Kind string
	Hash string
	Path string
}

func main() {
	updates := make(chan struct{})
	startTime = time.Now()
	freedSpace = 0
	modulesDeleted = 0
	var wg sync.WaitGroup

	// Define a new "start" channel
	start := make(chan struct{})

	// Modify your goroutine as follows
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Block until we receive a signal on the "start" channel
		<-start

		ticker := time.NewTicker(5 * time.Second) // Ticker that ticks every five seconds

		for {
			select {
			case <-ticker.C: // Every five seconds the ticker ticks
				printProgress()
			case _, ok := <-updates: // Whenever there is a deletion happening
				if !ok {
					ticker.Stop() // Stop the ticker when the updates channel is closed
					printProgress()
					return
				}
				printProgress()
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the username and subdirectory (format username/subdirectory): ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	userDir := "/Users/" + text

	// check if the user would like to automatically delete node_modules folders
	fmt.Print("Would you like to delete node_modules folders? (yes/no): ")
	var answer string
	_, _ = fmt.Scan(&answer)
	deleteNodeModules := answer == "yes"

	// Ask if the user would like to delete duplicate files
	fmt.Print("Would you like to delete duplicate files? (yes/no): ")
	var duplicateAnswer string
	_, _ = fmt.Scan(&duplicateAnswer)
	deleteDuplicates := duplicateAnswer == "yes"

	var minSize int64

	// If the user does not want to delete duplicate files, terminate the program.
	if !deleteDuplicates {
		fmt.Println("No operation selected. Exiting the program.")
		os.Exit(0)
	}

	if deleteDuplicates {
		// Ask for the minimum file size
		fmt.Print("Enter the minimum file size for duplicate deletion in KB (files smaller than this will be ignored): ")
		var minSizeKB int64
		_, _ = fmt.Scan(&minSizeKB)
		minSize = minSizeKB * 1024
	}

	close(start)

	if deleteNodeModules {
		err := filepath.Walk(userDir, func(path string, f os.FileInfo, err error) error {
			return deleteNodeModulesFunc(path, f, err, updates)
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	if deleteNodeModules && modulesDeleted == 0 {
		fmt.Println("\nNo node_modules directories were found.")
	}

	fileMap := make(map[string]File)
	duplicates := make([]File, 0)
	deleteChan := make(chan []File)

	go func() {
		for files := range deleteChan {
			printAndDeleteBatch(files, updates)
		}
	}()

	_ = filepath.Walk(userDir, visit(fileMap, &duplicates, deleteChan, minSize))

	// Send remaining duplicates if there are any
	if len(duplicates) > 0 {
		deleteChan <- duplicates
	}
	close(updates)
	wg.Wait()
	close(deleteChan)
}

// Your printProgress function that prints the progress
func printProgress() {
	elapsed := int(time.Since(startTime).Seconds()) // Convert to integer seconds
	fmt.Printf("\033[2K\rFreed space: %s      Time elapsed: %ds", formatSize(freedSpace), elapsed)
}

func deleteNodeModulesFunc(path string, f os.FileInfo, err error, updates chan<- struct{}) error {
	if err != nil {
		return err
	}
	if f.IsDir() && f.Name() == "node_modules" {
		dirSize, err := getDirSize(path)
		if err != nil {
			return err
		}
		fmt.Printf("\033[32mnode_modules deleted, freeing %s of space\033[0m\n", formatSize(dirSize))
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
		freedSpace += dirSize
		updates <- struct{}{} // Trigger an update to the progress display
		return filepath.SkipDir
	}
	return nil
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func printAndDeleteBatch(files []File, updates chan<- struct{}) {
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
		freedSpace += totalSize // update the global variable
		updates <- struct{}{}   // trigger an update to the progress display
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
						*duplicates = make([]File, 0)
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
