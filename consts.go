package main

const help string = `
yapm (Yet Another Package Manager)

help - Prints this message (usage, yapm help)
install/i [package, OPTIONAL: --build] - Installs a package
make [filename] - Builds a package from a tar file
update [OPTIONAL: package] - Updates a/all package/s
pack [OPTIONAL: dir] - Packages a directory
register [repository] - Adds a repository to the list of repositories

Created by YendisFish`

var initRepos []string = []string{"THIS WILL HAVE A LINK ONE DAY"}
