package filesystem

import (
	"fmt"
	"os"

	"github.com/sqweek/dialog"
)

func ReadFile(pathToFile string) (string, bool) {
	contents, err := os.ReadFile(pathToFile)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", false
	}

	return string(contents), true
}

func WriteFile(pathToFile string, contents string) bool {
	err := os.WriteFile(pathToFile, []byte(contents), 0644)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}

	return true
}

func DeleteFile(pathToFile string) {
	os.Remove(pathToFile)
}

func OpenDirectory() (string, bool) {
	path, err := dialog.Directory().Title("Select directory...").Browse()
	if err != nil {
		if err != dialog.ErrCancelled {
			fmt.Printf("Error: %s\n", err.Error())
		}

		return "", false
	}

	return path, true
}

func DoesFileExist(fullPath string) bool {
	_, err := os.Stat(fullPath)
	return err == nil
}
