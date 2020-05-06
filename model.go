package main

// GoRepo is golang project
type GoRepo struct {
	Name    string
	URL     string
	Modules []GoModule
}

// GoModule is go pkg/lib
type GoModule struct {
	Name string
	URL  string
}

// GoModuleRank is imports count of go module
type GoModuleRank struct {
	URL   string
	Count int
}
