package src

type Flags struct {
	Help    bool `short:"h" long:"help" description:"Display help" global:"true"`
	Version bool `short:"v" long:"version" description:"Display version"`
	Init    struct {
		Settings bool `settings:"true" allow-unknown-arg:"true"`
		Print    bool `short:"p" long:"print" description:"Print generated content to stdout instead of writing .gitignore"`
	} `command:"init" description:"Initialize .gitignore - use ignoreinit init <template...> <location>" nonempty:"true"`
	Replace struct {
		Settings bool `settings:"true" allow-unknown-arg:"true"`
		Print    bool `short:"p" long:"print" description:"Print generated content to stdout instead of replacing .gitignore"`
	} `command:"replace" description:"Replace current .gitignore - use ignoreinit replace <template...> <location>" nonempty:"true"`
	Merge struct {
		Settings bool `settings:"true" allow-unknown-arg:"true"`
		Print    bool `short:"p" long:"print" description:"Print merged content to stdout instead of writing .gitignore"`
	} `command:"merge" description:"Merge gitignore templates into current .gitignore - use ignoreinit merge <template...> <location>" nonempty:"true"`
	List struct {
		Settings bool `settings:"true" allow-unknown-arg:"true"`
	} `command:"list" description:"List available .gitignore templates"`
	Search struct {
		Settings bool `settings:"true" allow-unknown-arg:"true"`
	} `command:"search" description:"Search available .gitignore templates - use ignoreinit search <term>" nonempty:"true"`
	Completion struct {
		Settings bool `settings:"true" allow-unknown-arg:"true"`
	} `command:"completion" description:"Generate shell completion - use ignoreinit completion <bash|zsh|fish|powershell>" nonempty:"true"`
}
