package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetHomeDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	homedir = filepath.Join(homedir, ".yapm")

	return homedir
}
