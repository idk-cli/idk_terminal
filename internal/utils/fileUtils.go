package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func BackupFile(filePath string) error {
	// Open the original file
	originalFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	// Create a backup file path by appending .idk.backup to the original file path
	backupFilePath := filePath + ".idk.backup"

	// Create/open the backup file
	backupFile, err := os.Create(backupFilePath)
	if err != nil {
		return err
	}
	defer backupFile.Close()

	// Copy the contents of the original file to the backup file
	_, err = io.Copy(backupFile, originalFile)
	if err != nil {
		return err
	}

	// Ensure all writes are flushed to the backup file
	err = backupFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func GetAbsoluteHomeDirectoryPath(paths []string) string {
	// Attempt to get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If there's an error, return an empty string
		return ""
	}

	// Start with the home directory as the base for the absolute path
	absolutePath := homeDir

	// Range over the paths slice and join each segment to absolutePath
	for _, v := range paths {
		absolutePath = filepath.Join(absolutePath, v)
	}

	// Return the constructed absolute path
	return absolutePath
}

// ListFilesAndDirs lists all files and directories in the current working directory.
func ListFilesAndDirs() ([]string, error) {
	// Get the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Read the directory contents
	entries, err := os.ReadDir(pwd)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name()) // entry.Name() gives just the name, not the path
	}

	return names, nil
}

func FindReadmeData() (string, error) {
	// Get the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// Read the directory contents
	files, err := os.ReadDir(pwd)
	if err != nil {
		return "", err
	}

	// Look for a README file
	for _, file := range files {
		if !file.IsDir() && isReadmeFile(file.Name()) {
			content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", pwd, file.Name()))
			if err != nil {
				return "", err
			}
			return string(content), nil
		}
	}

	// No README file found
	return "", nil
}

// isReadmeFile checks if the filename matches common README patterns.
func isReadmeFile(filename string) bool {
	lowerFilename := strings.ToLower(filename)
	return strings.HasPrefix(lowerFilename, "readme") && (strings.HasSuffix(lowerFilename, ".md") || strings.HasSuffix(lowerFilename, ".txt") || lowerFilename == "readme")
}

func FindMakefileData() (string, error) {
	// Get the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// Read the directory contents
	files, err := os.ReadDir(pwd)
	if err != nil {
		return "", err
	}

	// Look for a Makefile
	for _, file := range files {
		if !file.IsDir() && isMakefile(file.Name()) {
			content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", pwd, file.Name()))
			if err != nil {
				return "", err
			}
			return string(content), nil
		}
	}

	// No Makefile found
	return "", nil
}

// isMakefile checks if the filename matches common Makefile names.
func isMakefile(filename string) bool {
	return filename == "Makefile" || filename == "makefile"
}

func GetCurrentDirName() (string, error) {
	// Get the full path of the current working directory.
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Extract the directory name from the full path.
	dirName := filepath.Base(pwd)
	return dirName, nil
}
