package pack

import (
	"archive/tar"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"yapm/logger"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

func Pack(args []string) {
	dir := "."
	if len(args) != 0 {
		dir = args[0]
	}

	var conf Config = ReadConfig(dir)

	fmt.Println()
	logger.LogRaw("[", color.FgWhite, color.Bold)
	logger.LogRaw("Package", color.FgBlue, color.Bold)
	logger.LogRawln("]", color.FgWhite, color.Bold)

	logger.LogString(conf.Package.Name, "Name", color.FgMagenta, color.Bold)
	logger.LogString(conf.Package.Author, "Author", color.FgMagenta, color.Bold)
	logger.LogString(conf.Package.Version, "Version", color.FgMagenta, color.Bold)

	logger.LogRaw("Repository: ", color.FgMagenta, color.Bold)
	for _, v := range conf.Package.Repository {
		fmt.Print("(" + v + ")")
	}
	fmt.Println()

	logDependencies(conf.Dependencies)

	hashes := RetrieveHashes(dir)
	writePkg(hashes, &conf, dir)
}

func ReadConfig(dir string) Config {
	content, err := os.ReadFile(filepath.Join(dir, "pkg.toml"))
	if err != nil {
		fmt.Println("Could not find pkg.toml!")
	}

	tomlContent := string(content)

	var conf Config
	_, e := toml.Decode(tomlContent, &conf)
	if e != nil {
		fmt.Println("Could not parse pkg.toml!")
	}

	return conf
}

func logDependencies(dependenices map[string][]string) {
	//check for dependencies
	//install ones that do not already exists

	fmt.Println()
	logger.LogRaw("[", color.FgWhite, color.Bold)
	logger.LogRaw("Dependencies", color.FgGreen, color.Bold)
	logger.LogRawln("]", color.FgWhite, color.Bold)

	for k, v := range dependenices {
		logger.LogRaw(k+": ", color.FgBlue, color.Bold)
		for _, val := range v {
			logger.LogRaw("(" + val + ")")
		}
		fmt.Println()
	}
}

func readDirRecurse(dir string) map[string]string {
	fles := map[string]string{}

	og, e1 := os.Getwd()
	if e1 != nil {
		logger.Error(logger.CreateLogEntry(e1.Error()))
	}

	e2 := os.Chdir(dir)
	if e2 != nil {
		logger.Error(logger.CreateLogEntry(e2.Error()))
	}

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			abs, err := filepath.Abs(path)
			if err != nil {
				logger.Error(logger.CreateLogEntry(err.Error()))
			}

			fles[abs] = path
		}

		return nil
	})

	if err != nil {
		logger.Error(logger.CreateLogEntry("Could not read directory!"))
	}

	e3 := os.Chdir(og)
	if e3 != nil {
		logger.Error(logger.CreateLogEntry(e3.Error()))
	}

	return fles
}

func RetrieveHashes(dir string) map[string]string {
	fles := readDirRecurse(dir)

	fmt.Println()
	logger.LogRawln("Hashing Files...")

	ret := map[string]string{}
	for k, v := range fles {

		fle, err := os.ReadFile(k)
		if err != nil {
			logger.Error(logger.CreateLogEntry("Could not read "+k+"!", []string{"Error", err.Error()}))
		}

		dat := sha256.Sum256(fle)
		str := fmt.Sprintf("%x", dat)

		ret[v] = str

		logger.LogString(str, k, color.FgGreen, color.Bold)
	}

	return ret
}

func writePkg(hashes map[string]string, config *Config, dir string) {
	fmt.Println()

	logger.LogRawln("Creating tar file...")
	fle, err := os.Create("./" + config.Package.Name + "." + config.Package.Version + ".tar")
	if err != nil {
		logger.Error(logger.CreateLogEntry("Could not create tar package!", []string{"Error", err.Error()}))
	}

	fles := readDirRecurse(dir)

	logger.LogRawln("Writing tar contents...")

	tarWriter := tar.NewWriter(fle)
	defer tarWriter.Close()

	for k, v := range fles {
		logger.Log(logger.CreateLogEntry("Writing "+k+" to tar file"), "Tar Writer", color.FgBlue, color.Bold)

		fle, err := os.ReadFile(k)
		if err != nil {
			logger.Error(logger.CreateLogEntry("Could not read "+v+"!", []string{"Error", err.Error()}))
		}

		tarWriter.WriteHeader(&tar.Header{
			Name: v,
			Size: int64(len(fle)),
		})

		_, err2 := tarWriter.Write(fle)
		if err2 != nil {
			logger.Error(logger.CreateLogEntry("Failed to write to tar pkg!"))
		}
	}

	file, err := os.Create("./hashes.csv")
	if err != nil {
		logger.Error(logger.CreateLogEntry("Could not create hashes file!", []string{"Error", err.Error()}))
	}

	for k, v := range hashes {
		file.WriteString(k + ", " + v + "\n")
	}

	csvInf, err := file.Stat()
	if err != nil {
		logger.Error(logger.CreateLogEntry("Could not find size of hashes.csv!", []string{"Error", err.Error()}))
	}

	tarWriter.WriteHeader(&tar.Header{
		Name: "hashes.csv",
		Size: csvInf.Size(),
	})

	file.Close()

	csvBytes, err := os.ReadFile("./hashes.csv")
	if err != nil {
		logger.Error(logger.CreateLogEntry("Could not embed hashes file!", []string{"Error", err.Error()}))
	}

	tarWriter.Write(csvBytes)
	os.Remove("./hashes.csv")

	fmt.Println()
	logger.LogRawln("Created pkg!")
}
