package src

type Flags struct {
	Help    bool `short:"h" long:"help" description:"Display help" global:"true"`
	Version bool `short:"v" long:"version" description:"Display version"`
	Init    struct {
		Settings bool `settings:"true" allow-unknown-arg:"true"`
	} `command:"init" description:"Initialize .gitignore - use ignoreinit init <language> <location>" nonempty:"true"`
	Replace struct {
		Settings bool `settings:"true" allow-unknown-arg:"true"`
	} `command:"replace" description:"Replace current .gitignore - use ignoreinit replace <language> <location>" nonempty:"true"`
}
