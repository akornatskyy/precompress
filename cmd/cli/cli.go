package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

type Options struct {
	Paths    []string
	MinSize  int64
	MaxDepth int
}

func ParseArgs() (Options, error) {
	var opts Options
	var showHelp bool
	var showVersion bool

	flag.Int64Var(&opts.MinSize, "min-size", 1024, "Minimum content size")
	flag.IntVar(&opts.MaxDepth, "max-depth", 0, "Descend at most levels of directories")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Show help message (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "Show version info and exit")
	flag.BoolVar(&showVersion, "v", false, "Show version info and exit (shorthand)")
	flag.Parse()

	if showHelp {
		printUsage()
		os.Exit(0)
	}

	if showVersion {
		printVersion()
		os.Exit(0)
	}

	opts.Paths = flag.Args()
	return opts, nil
}

func printUsage() {
	fmt.Println(`Usage:
  precompress [options] <files,dirs,...>

Options:
  --min-size <bytes>     Minimum content size (default 1024)
  --max-depth <levels>   Descend at most levels of directories
  -h, --help             Show this help message
  -v, --version          Show version`)
}

func printVersion() {
	fmt.Printf(`Version:           %s
Go version:        %s
Git commit:        %s
Built:             %s
OS/Arch:           %s/%s
`,
		Version,
		runtime.Version(),
		GitCommit,
		BuildDate,
		runtime.GOOS,
		runtime.GOARCH,
	)
}
