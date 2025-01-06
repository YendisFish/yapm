package make

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"yapm/lib"
	"yapm/logger"
	"yapm/pack"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

func Make(args []string) {
	path := args[0]

	fmt.Println()
	logger.LogString(path, "Making", color.FgBlue, color.Bold)

	logger.Info(logger.CreateLogEntry("Extracting configs"))
	configs := readTarConfigs(path)

	logger.Info(logger.CreateLogEntry("Reading hashes.csv"))
	for k, v := range configs.Hashes {
		logger.LogRaw(strings.Trim(k, " "), color.FgMagenta, color.Bold)
		fmt.Print(": ")
		logger.LogRawln(strings.Trim(v, " "), color.FgBlue, color.Bold)
	}

	fmt.Println()

	logger.Info(logger.CreateLogEntry("Reading file destinations"))
	for k, v := range configs.Definitions {
		logger.LogRaw(strings.Trim(k, " "), color.FgMagenta, color.Bold)
		fmt.Print(": ")
		logger.LogRawln(strings.Trim(v, " "), color.FgBlue, color.Bold)
	}

	fmt.Println()

	logger.Info(logger.CreateLogEntry("Looking for cached files"))
	fles := getValidFiles(configs.Hashes, configs.Pkg.Package.Name)
	logger.Info(logger.CreateLogEntry("Extracting following files"))
	for k, v := range fles {
		logger.LogRaw(strings.Trim(k, " "), color.FgMagenta, color.Bold)
		fmt.Print(": ")
		logger.LogRawln(strings.Trim(v, " "), color.FgBlue, color.Bold)
	}

	fmt.Println()
}

func readTarConfigs(path string) ConfigBundle {
	var ret ConfigBundle = ConfigBundle{}

	fle, err := os.Open(path)
	if err != nil {
		logger.Error(logger.CreateLogEntry(err.Error()))
	}

	tread := tar.NewReader(fle)
	for {
		header, err := tread.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(logger.CreateLogEntry(err.Error()))
		}

		if header.Name == "definitions.json" {
			var buf bytes.Buffer
			_, err := io.Copy(&buf, tread)
			if err != nil {
				logger.Error(logger.CreateLogEntry(err.Error()))
			}

			json.Unmarshal(buf.Bytes(), &ret.Definitions)
		}

		if header.Name == "hashes.csv" {
			var buf bytes.Buffer
			_, err := io.Copy(&buf, tread)
			if err != nil {
				logger.Error(logger.CreateLogEntry(err.Error()))
			}

			ret.Hashes = readTarCsv(buf.String())
		}

		if header.Name == "pkg.toml" {
			var buf bytes.Buffer
			_, err := io.Copy(&buf, tread)
			if err != nil {
				logger.Error(logger.CreateLogEntry(err.Error()))
			}

			ret.Pkg = readConfig(buf.String())
		}
	}

	return ret
}

func readTarCsv(content string) map[string]string {
	ret := map[string]string{}
	lines := []string{}

	cur := ""
	for _, c := range content {
		if c == '\n' {
			lines = append(lines, cur)
			cur = ""

			continue
		}

		cur = cur + string(c)
	}

	for _, str := range lines {
		fields := strings.Split(str, ",")
		ret[fields[0]] = fields[1]
	}

	return ret
}

func readConfig(content string) pack.Config {
	tomlContent := content

	var conf pack.Config
	_, e := toml.Decode(tomlContent, &conf)
	if e != nil {
		fmt.Println("Could not parse pkg.toml!")
	}

	return conf
}

func getValidFiles(full map[string]string, pkgname string) map[string]string {
	homedir := lib.GetHomeDir()

	cache := filepath.Join(homedir, pkgname)
	if _, err := os.Stat(cache); !os.IsNotExist(err) {
		logger.Info(logger.CreateLogEntry("Found cached files!"))
		ret := map[string]string{}

		dir, err := os.ReadDir(cache)
		if err != nil {
			logger.Error(logger.CreateLogEntry(err.Error()))
		}

		cached := map[string]string{}
		for _, entry := range dir {
			if entry.Name() == "hashes.csv" {
				content, err := os.ReadFile(filepath.Join(cache, "hashes.csv"))
				if err != nil {
					logger.Error(logger.CreateLogEntry(err.Error()))
				}

				cached = readTarCsv(string(content))
			}
		}

		for k, v := range full {
			value, ok := cached[k]
			if !ok || v != value {
				ret[k] = v
			}

			return ret
		}
	}

	return full
}
