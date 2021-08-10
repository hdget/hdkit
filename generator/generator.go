package generator

// Generator represents a generator.
type Generator interface {
	PreGenerate() error                // pre actions for generator
	Generate(concrete Generator) error // main generate function
	PostGenerate() error               // post actions
	GetGenCodeFuncs() []func()         // get gen code functions
}
