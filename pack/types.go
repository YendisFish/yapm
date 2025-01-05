package pack

type Config struct {
	Package      Package
	Dependencies map[string][]string
}

type Package struct {
	Name       string
	Author     string
	Version    string
	Repository []string
}
