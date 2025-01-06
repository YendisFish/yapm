package make

import "yapm/pack"

type ConfigBundle struct {
	Hashes      map[string]string
	Definitions map[string]string
	Pkg         pack.Config
}
