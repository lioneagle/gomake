package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lioneagle/goutil/src/file"
)

type CoverageConfig struct {
	Help     bool
	Packages string
	Verbose  bool
	ShowHtml bool
	Regexp   string
	flagset  *flag.FlagSet
}

func (this *CoverageConfig) makeFlags() {
	this.flagset = flag.NewFlagSet("coverage", flag.ExitOnError)
	this.flagset.Usage = this.usage
	this.flagset.BoolVar(&this.Help, "help", false, "")
	this.flagset.StringVar(&this.Packages, "packages", ".", "Apply coverage analysis in each test to the given list of packages.")
	this.flagset.BoolVar(&this.Verbose, "v", false, "Verbose output.")
	this.flagset.BoolVar(&this.ShowHtml, "html", false, "Show html.")
	this.flagset.StringVar(&this.Regexp, "run", ".", "Run only those tests and examples matching the regular expression.")
}

func (this *CoverageConfig) usage() {
	fmt.Printf("%s",
		`Usage: coverage [flags]

    -packages pkg1,pkg2,pkg3
        Apply coverage analysis in each test to the given list of packages (default "." means all packages).

    -html
        Show html.

    -run regexp
        Run only those tests and examples matching the regular expression (default ".").

    -v
        Verbose output (default false).
`)
}

func (this *CoverageConfig) parse() {
	if len(os.Args) > 2 {
		this.flagset.Parse(os.Args[2:])
	}

	if this.Help {
		this.usage()
		os.Exit(2)
	}
}

type BenchmarkConfig struct {
	Help      bool
	Package   string
	BenchTime int
	Regexp    string
	GoTorch   bool
	flagset   *flag.FlagSet
}

func (this *BenchmarkConfig) makeFlags() {
	this.flagset = flag.NewFlagSet("benchmark", flag.ExitOnError)
	this.flagset.Usage = this.usage
	this.flagset.BoolVar(&this.Help, "help", false, "")
	this.flagset.IntVar(&this.BenchTime, "benchtime", 1, "Run enough iterations of each benchmark to take t, specified as seconds.")
	this.flagset.StringVar(&this.Regexp, "run", ".", "Run only those benchmarks matching a regular expression.")
	this.flagset.BoolVar(&this.GoTorch, "torch", false, "Run go-torch to get fireflame graph (default false).")
}

func (this *BenchmarkConfig) usage() {
	fmt.Printf("%s",
		`Usage: benchmark package [flags]

    -benchtime
        Run enough iterations of each benchmark to take t, specified as seconds (default 1 second).

    -run regexp
        Run only those tests and examples matching the regular expression (default ".").

    -torch
        Run go-torch to get fireflame graph (default false).

`)
}

func (this *BenchmarkConfig) parse() {
	if len(os.Args) < 3 {
		this.usage()
		os.Exit(2)
	}
	this.Package = os.Args[2]

	if os.Args[2][0] == '-' {
		this.usage()
		os.Exit(2)
	}

	if len(os.Args) > 3 {
		this.flagset.Parse(os.Args[3:])
	}

	if this.Help {
		this.usage()
		os.Exit(2)
	}
}

type InstallConfig struct {
	Help           bool
	Package        string
	OutputName     string
	Arch           string
	WithArchSuffix bool
	CopyToStdGoBin bool
	flagset        *flag.FlagSet
}

func (this *InstallConfig) usage() {
	fmt.Printf("%s",
		`Usage: install package [outputname] [flags]

    -output
        Output name.

    -arch  64 | 32 | all
        Output architure (default 64).

    -arch-suffix  
        Output filename with architecture suffix (default true).

    -copy2gobin  
        Copy file to std gobin (default false).
`)
}

func (this *InstallConfig) makeFlags() {
	this.flagset = flag.NewFlagSet("install", flag.ExitOnError)
	this.flagset.Usage = this.usage
	this.flagset.BoolVar(&this.Help, "help", false, "")
	this.flagset.StringVar(&this.OutputName, "output", "", "Output name")
	this.flagset.StringVar(&this.Arch, "arch", "64", "Output architure")
	this.flagset.BoolVar(&this.WithArchSuffix, "arch-suffix", true, "Output filename with architecture suffix")
	this.flagset.BoolVar(&this.CopyToStdGoBin, "copy2gobin", false, "Copy file to std gobin")
}

func (this *InstallConfig) parse() {
	if len(os.Args) < 3 {
		this.usage()
		os.Exit(2)
	}

	if os.Args[2][0] == '-' {
		this.usage()
		os.Exit(2)
	}

	this.Package = os.Args[2]

	if len(os.Args) > 3 {
		this.flagset.Parse(os.Args[3:])
	}

	if this.OutputName == "" {
		this.OutputName = this.Package
	}

	if this.Help {
		this.usage()
		os.Exit(2)
	}
}

type PProfConfig struct {
	Help      bool
	Package   string
	NodeCount int
	flagset   *flag.FlagSet
}

func (this *PProfConfig) makeFlags() {
	this.flagset = flag.NewFlagSet("pprof", flag.ExitOnError)
	this.flagset.Usage = this.usage
	this.flagset.BoolVar(&this.Help, "help", false, "")
	this.flagset.IntVar(&this.NodeCount, "nodecount", 30, "Max number of nodes to show")
}

func (this *PProfConfig) usage() {
	fmt.Printf("%s",
		`Usage: pprof package [flags]

    -nodecount
        Max number of nodes to show (default 30).
`)
}

func (this *PProfConfig) parse() {
	if len(os.Args) < 3 {
		this.usage()
		os.Exit(2)
	}

	this.Package = os.Args[2]

	if os.Args[2][0] == '-' {
		this.usage()
		os.Exit(2)
	}

	if len(os.Args) > 3 {
		this.flagset.Parse(os.Args[3:])
	}

	if this.Help {
		this.usage()
		os.Exit(2)
	}
}

type RunConfig struct {
	Command   string
	Coverage  CoverageConfig
	Benchmark BenchmarkConfig
	Install   InstallConfig
	Pprof     PProfConfig

	OldGobin string
}

func NewRunConfig() *RunConfig {
	config := &RunConfig{}
	config.Coverage.makeFlags()
	config.Benchmark.makeFlags()
	config.Install.makeFlags()
	config.Pprof.makeFlags()

	return config
}

func (this *RunConfig) usage() {
	name := file.RemoveFileSuffix(filepath.Base(os.Args[0]))
	fmt.Printf(
		`Usage: %s command [arguments]

The commands are:

    install     compile and install packages and dependencies
    coverage    run test and show coverage
    benchmark   run benchmark and show result
    pprof       run pprof for one package

Use "%s help [command]" for more information about a command.
`, name, name)
}

func (this *RunConfig) Parse() {
	name := file.RemoveFileSuffix(filepath.Base(os.Args[0]))
	if len(os.Args) < 2 {
		this.usage()
		os.Exit(2)
	}

	this.Command = os.Args[1]

	switch os.Args[1] {
	case "-help":
		this.usage()
		os.Exit(2)
	case "help":
		if len(os.Args) <= 2 {
			this.usage()
			os.Exit(2)
		}

		switch os.Args[2] {
		case "install":
			this.Install.usage()
		case "coverage":
			this.Coverage.usage()
		case "benchmark":
			this.Benchmark.usage()
		case "pprof":
			this.Pprof.usage()
		default:
			fmt.Printf("Unknown help topic \"%s\"\n", os.Args[2])
			fmt.Printf("Run \"%s help\" for usage.\n", name)
			os.Exit(2)
		}
	case "install":
		this.Install.parse()
	case "coverage":
		this.Coverage.parse()
	case "benchmark":
		this.Benchmark.parse()
	case "pprof":
		this.Pprof.parse()

	default:
		if os.Args[1][0] == '-' {
			fmt.Printf("flag provided but not defined: %s\n", os.Args[1])
			this.usage()
		} else {
			fmt.Printf("Unknown subcommand \"%s\"\n", os.Args[1])
			fmt.Printf("Run '%s help' for usage.\n", name)
		}
		os.Exit(2)
	}

	if os.Args[1] == "help" {

	} else {

	}

}
