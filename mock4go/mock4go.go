package main

import (
	"fmt"
	"github.com/jvshahid/mock4go"
	"os"
	"os/exec"
	"path"
	"strings"
)

func createTempDir(args *Args) (string, error) {
	name := path.Join(args.Destination, "src")
	return name, os.MkdirAll(name, os.ModePerm)
}

type Args struct {
	Verbose        bool
	Keep           bool
	InstrumentOnly bool
	Destination    string
	cmd            []string // the command to run and its arguments
	cmdArgs        []string // the command to run and its arguments
	packages       []string // the list of packages
	testArgs       []string // the arguments to the test binary
}

func NewArgs() *Args {
	return &Args{
		Keep:    false,
		Verbose: false,
	}
}

func readMock4goArgs(args *Args) (int, error) {
	i := 1
	for ; i < len(os.Args) && strings.HasPrefix(os.Args[i], "-"); i++ {
		switch strings.ToLower(os.Args[i]) {
		case "-k", "--keep":
			args.Keep = true
		case "-d", "--destination":
			i++
			args.Destination = os.Args[i]
		case "-v", "--verbose":
			args.Verbose = true
		case "-i", "--instrument-only":
			args.InstrumentOnly = true
		}
	}

	if args.Destination == "" {
		args.Destination = path.Join(os.TempDir(), "mock4go")
	}

	if args.InstrumentOnly && !args.Keep {
		return -1, fmt.Errorf("Error: cannot use -i without -k")
	}
	return i, nil
}

func parseNamesAndArgs(lastIndex int) ([]string, []string, int) {
	cmd := make([]string, 0)
	cmdArgs := make([]string, 0)

	parsingArgs := false
	i := lastIndex

	for ; i < len(os.Args); i++ {
		fmt.Printf("argument: %s\n", os.Args[i])
		// do we have
		if strings.HasPrefix(os.Args[i], "-") {
			parsingArgs = true
			if os.Args[i] != "--" {
				cmdArgs = append(cmdArgs, os.Args[i])
			}
		} else if parsingArgs {
			// we are here right now: cmd -arg1 -arg2 package
			//                                        ^
			break
		} else {
			cmd = append(cmd, os.Args[i])
		}
	}
	return cmd, cmdArgs, i
}

func fixArgs(args *Args) {
	if len(args.cmd) == 0 {
		args.cmd = []string{"go"}
		args.cmdArgs = []string{"test"}
	}
}

func parseArgs() (*Args, error) {
	args := NewArgs()
	lastIdx, err := readMock4goArgs(args)
	if err != nil {
		return nil, err
	}
	// asssume that we have the command to run
	cmd, cmdArgs, lastIdx := parseNamesAndArgs(lastIdx)
	// if we reached the end of the arguments then we must have parsed the packages not the command
	if lastIdx == len(os.Args) {
		args.packages = cmd
		args.testArgs = cmdArgs
		fixArgs(args)
		return args, nil
	}
	args.cmd = cmd
	args.cmdArgs = cmdArgs

	packages, testArgs, _ := parseNamesAndArgs(lastIdx)
	args.packages = packages
	args.testArgs = testArgs

	return args, nil
}

func printUsage() {
	usage := `
Usage: mock4go [mock4go arguments] [test command] [test command arguments] [--] [package names] [test binary args]

Where:
  [mock4go arguments]:
    -v: enable debug output
    -d|--destination: destination directory where instrumented code will be created
    -k|--keep: don't delete instrumented code after running the tests
    -i|--instrument-only: don't run the tests, only instrument the code (error if used without -k)
  [test command]:
    The command to use to run the tests, e.g. mock4go go test ...., or mock4go gocov ....
    If not specified, it will default to 'go test'
  [test command arguments]:
    The test command arguments, see the usage help for the command you use for testing
  [--]:
    You should use -- if you didn't pass any arguments to the test command, otherwise
    there is no way to tell when the test command end and the package names start
  [package names]:
    A space delimited list of package names.
  [test binary args]:
    The arguments to pass the test binary created.

Examples:
  mock4go go test -v db -database=localhost:8080 (use the go test command with -v argument to test the db package)
  mock4go gocov -v db -database=localhost:8080   (use gocov instead)
  mock4go db -database=localhost:8080            (default to go test if the test command wasn't specified)
`
	fmt.Printf(usage)
}

func run() int {
	args, err := parseArgs()
	fmt.Printf("args: %#v\n", args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		printUsage()
		return 2
	}

	if len(args.packages) == 0 {
		fmt.Fprintf(os.Stderr, "No packages was specified on the command line", err)
		printUsage()
		return 2
	}

	tmpDir, err := createTempDir(args)

	defer func() {
		if !args.Keep {
			err := os.RemoveAll(args.Destination)
			fmt.Printf("removing directory %s\n", args.Destination)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot remove directory %s. Error: %s\n", args.Destination, err)
			}
		}
	}()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create temporary directory. Error: %s\n", err)
		return 2
	}

	pkgs := make([]string, 0, len(args.packages))

	for _, packageName := range args.packages {
		pkg, err := api.InstrumentPackage(packageName, tmpDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
			return 2
		}
		pkgs = append(pkgs, pkg.Name)
	}
	api.InstrumentPackage(api.Mock4goImport, tmpDir)

	// run the tests
	cmd := args.cmd
	cmd = append(cmd, args.cmdArgs...)
	cmd = append(cmd, pkgs...)
	cmd = append(cmd, args.testArgs...)
	goBinPath, err := exec.LookPath(args.cmd[0])
	os.Setenv("GOPATH", strings.Replace(tmpDir, "/src", "", -1))

	if !args.InstrumentOnly {
		fmt.Printf("command: %v\n", cmd)

		proc, err := os.StartProcess(goBinPath, cmd, &os.ProcAttr{
			Env: os.Environ(),
			Files: []*os.File{
				os.Stdin,
				os.Stdout,
				os.Stderr,
			},
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return 2
		}
		status, err := proc.Wait()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return 2
		}
		if status.Success() {
			return 0
		}
	} else {
		return 0
	}

	return 1
}

func main() {
	os.Exit(run())
}
