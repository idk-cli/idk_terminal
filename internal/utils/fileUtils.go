package utils

import (
	"io"
	"os"
	"path/filepath"
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
