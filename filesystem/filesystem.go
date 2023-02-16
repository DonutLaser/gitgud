package filesystem

import (
	"fmt"
	"os"
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

func DoesFileExist(fullPath string) bool {
	_, err := os.Stat(fullPath)
	return err == nil
}
