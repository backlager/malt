package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/backlager/base-malt/commands"
)

const (
	// AppName is the application name
	AppName = "backlager"

	// AppVersion is the application version
	AppVersion = "0.1.0"
)

var (
	// Subcommands lists the available commands and help topics.
	// The order here is the order in which they are printed by help.
	subcommands = []*commands.Command{
		commands.Add,
		commands.Edit,
		commands.Set,
		commands.Delete,
		commands.Sync,

		//helpConfigure,
		//helpSend,
	}

	// exitStatus
	exitStatus = 0

	// exitMu
	exitMu sync.Mutex

	// usageTemplate
	usageTemplate = `backlager is a distributed project management tool.

	Usage:

		backlager command [arguments]

	The commands are:
	{{range .}}{{if .Runnable}}
		{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

	Use "backlager help [command]" for more information about a command.

	Additional help topics:
	{{range .}}{{if not .Runnable}}
		{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

	Use "backlager help [topic]" for more information about that topic.

	`

	// helpTemplate
	helpTemplate = `{{if .Runnable}}usage: backlager {{.UsageLine}}

	{{end}}{{.Long | trim}}
	`
	// atexitFuncs
	atexitFuncs []func()
)

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range subcommands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				_ = cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}

			cmd.Run(cmd, args)
			exit()
			return
		}
	}

	fmt.Fprintf(os.Stderr, "%s: unknown subcommand %q\nRun 'help' for usage.\n", AppName, args[0])
	setExitStatus(2)
	exit()
}

func setExitStatus(n int) {
	exitMu.Lock()

	if exitStatus < n {
		exitStatus = n
	}

	exitMu.Unlock()
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, subcommands)
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

// help implements the 'help' command.
func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'backlager help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: backlager help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'backlager help'
	}

	arg := args[0]

	for _, cmd := range subcommands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'backlager help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q. Run 'backlager help'.\n", arg)
	os.Exit(2) // failed at 'backlager help cmd'
}

func atexit(f func()) {
	atexitFuncs = append(atexitFuncs, f)
}

func exit() {
	for _, f := range atexitFuncs {
		f()
	}
	os.Exit(exitStatus)
}

func fatalf(format string, args ...interface{}) {
	errorf(format, args...)
	exit()
}

func errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
	setExitStatus(1)
}

func exitIfErrors() {
	if exitStatus != 0 {
		exit()
	}
}

func run(cmdargs ...interface{}) {
	cmdline := stringList(cmdargs...)
	//if buildN || buildX {
	//	fmt.Printf("%s\n", strings.Join(cmdline, " "))
	//	if buildN {
	//		return
	//	}
	//}

	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		errorf("%v", err)
	}
}

// stringList's arguments should be a sequence of string or []string values.
// stringList flattens them into a single []string.
func stringList(args ...interface{}) []string {
	var x []string
	for _, arg := range args {
		switch arg := arg.(type) {
		case []string:
			x = append(x, arg...)
		case string:
			x = append(x, arg)
		default:
			panic("stringList: invalid argument")
		}
	}
	return x
}
