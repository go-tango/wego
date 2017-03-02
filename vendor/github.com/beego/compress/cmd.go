package compress

import (
	"flag"
	"fmt"
	"os"
)

type commander interface {
	Parse([]string)
	Run(*Settings) error
}

var (
	commands = make(map[string]commander)
)

func (s *Settings) RunCommand(params ...string) {
	cmdName := "compress"
	if len(params) > 0 {
		cmdName = params[0]
	}

	if len(os.Args) < 2 || os.Args[1] != cmdName {
		return
	}

	args := argString(os.Args[2:])
	name := args.Get(0)

	if name == "help" {
		printHelp()
	}

	if cmd, ok := commands[name]; ok {
		cmd.Parse(os.Args[3:])
		cmd.Run(s)
		os.Exit(0)
	} else {
		if name == "" {
			printHelp()
		} else {
			printHelp(fmt.Sprintf("unknown command %s", name))
		}
	}
}

func (s *Settings) RunCompress(force, skip, verbose bool) {
	compressJsFiles(s, force, skip, verbose)
	compressCssFiles(s, force, skip, verbose)
}

func printHelp(errs ...string) {
	content := `compress command usage:

    js     - compress all js files
    css    - compress all css files
    all    - compress all files
`

	if len(errs) > 0 {
		fmt.Println(errs[0])
	}
	fmt.Println(content)
	os.Exit(2)
}

type compressAll struct {
	js, css bool
	force   bool
	skip    bool
	verbose bool
}

func (d *compressAll) Parse(args []string) {
	var name string
	if d.js && d.css {
		name = "all"
	} else if d.js {
		name = "js"
	} else {
		name = "css"
	}
	flagSet := flag.NewFlagSet("compress command: "+name, flag.ExitOnError)
	flagSet.BoolVar(&d.force, "force", false, "force compress file")
	flagSet.BoolVar(&d.skip, "skip", false, "skip all cached file")
	flagSet.BoolVar(&d.verbose, "v", false, "verbose info")
	flagSet.Parse(args)
}

func (d *compressAll) Run(s *Settings) error {
	if d.js {
		compressJsFiles(s, d.force, d.skip, d.verbose)
	}
	if d.css {
		compressCssFiles(s, d.force, d.skip, d.verbose)
	}
	return nil
}

func init() {
	commands["js"] = &compressAll{js: true}
	commands["css"] = &compressAll{css: true}
	commands["all"] = &compressAll{js: true, css: true}
}
