package install

import (
	"fmt"
	"yapm/logger"

	"github.com/fatih/color"
)

func Install(splice []string) {
	var pkgs []string
	var builds []string = nil

	for _, v := range splice {
		if v == "--build" {
			builds = []string{}
			continue
		}

		if builds != nil {
			builds = append(builds, v)
		} else {
			pkgs = append(pkgs, v)
		}
	}

	logInstallerMessage(pkgs, builds)

	install(pkgs)
	build(builds)
}

func install(pkgs []string) {
	for _, v := range pkgs {
		_ = v
		// ideally this will actually have a full on LogEntry that lists the info of the package and all that! ONE DAY
	}
}

func build(builds []string) {
	for _, v := range builds {
		_ = v
		// ideally this will actually have a full on LogEntry that lists the info of the package and all that! ONE DAY
	}
}

func logInstallerMessage(inst []string, bld []string) {
	fmt.Println("Installing Packages...")

	if inst != nil {
		logger.LogRaw("Installing: ", color.FgGreen, color.Bold)
		logger.PrintIndented(inst, "            ", 10)
	}

	if bld != nil {
		logger.LogRaw("Building: ", color.FgBlue, color.Bold)
		logger.PrintIndented(bld, "          ", 10)
	}
}

func checkExists() {

}

func checkExistsLinux() {

}

func checkExistsWindows() {

}
