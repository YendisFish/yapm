package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"yapm/install"
	make "yapm/make"
	"yapm/pack"
)

func main() {
	Initialize()

	args := os.Args[1:]
	err := ResolveArgs(args)

	if err != nil {
		fmt.Println(err)
	}
}

func ResolveArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("yapm expected an argument!")
	}

	switch args[0] {
	case "help":
		fmt.Println(help)
	case "install", "i":
		install.Install(args[1:])
	case "pack":
		pack.Pack(args[1:])
	case "make":
		make.Make(args[1:])
	}

	return nil
}

func Initialize() {
	homedir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println(err)
	}

	homedir = filepath.Join(homedir, ".yapm")

	inf, err := os.Stat(homedir)

	if err != nil || !inf.IsDir() {
		os.Mkdir(homedir, fs.ModeDir)
	}

	file, err := os.Create(filepath.Join(homedir, "repos.txt"))
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range initRepos {
		file.WriteString(v + "\n")
	}
}
