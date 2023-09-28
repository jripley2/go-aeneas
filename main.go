package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

var (
	logLevel = 0
	batch    = ""
	version  = "dev"
	commit   = "none"
	date     = "unknown"
)

func processArguments() {
	var (
		showHelp          *bool
		showVersion       bool
		showVersionNumber bool
	)

	// Parse flags
	// see: https://pkg.go.dev/github.com/spf13/pflag
	flag.IntVar(&logLevel, "verbose", 0, "verbose level")
	flag.StringVar(&batch, "batch", "", "batch JSON filename")
	flag.Lookup("verbose").NoOptDefVal = "1"
	flag.Lookup("verbose").Shorthand = "v"
	flag.BoolVar(&showVersion, "version", false, "display full version information")
	flag.BoolVar(&showVersionNumber, "version-number", false, "display version number")
	// Note: if we use BoolVar for help, we still see "pflag: help requested"
	showHelp = flag.BoolP("help", "h", false, "display help")
	flag.Parse()

	if logLevel > 0 {
		fmt.Fprintln(os.Stderr, "Logging level:", logLevel)
	}

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersionNumber {
		fmt.Println(version)
		os.Exit(0)
	}

	if showVersion {
		// GoReleaser automatically sets the version, commit and date
		// see: https://goreleaser.com/cookbooks/using-main.version/
		fmt.Printf("go-aeneas version %s (commit %s, built at %s)\n", version, commit, date)
		os.Exit(0)
	}
}

func main() {
	processArguments()

	tasks := []*Task{}
	if len(batch) > 0 {
		fmt.Println("Batch file:", batch)
		content, err := os.ReadFile(batch)
		if err != nil {
			log.Fatal("Error while reading batch file", err)
		}

		err = json.Unmarshal(content, &tasks)
		if err != nil {
			log.Fatal("Error parsing batch json file", err)
		}
	} else if len(os.Args) >= 5 {
		task := &Task{"", os.Args[1], os.Args[2], os.Args[3], os.Args[4]}
		tasks = append(tasks, task)
	}

	for _, task := range tasks {
		if len(task.Description) > 0 {
			fmt.Println("")
			fmt.Println("*** ", task.Description, " ***")
			fmt.Println("")
		}
		fmt.Println("Audio   : ", task.AudioFilename)
		fmt.Println("Phrase  : ", task.PhraseFilename)
		fmt.Println("Output  : ", task.OutputFilename)
		fmt.Println("Params  : ", task.Parameters)
	}
}
