package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hellflame/argparse"
)

var GlobalOptions struct {
	ManifestPath *string
}

func main() {

	/* DEFINE & PARSE ARGUMENTS */
	parser := argparse.NewParser("mg", "Utility to manage several git repositories simultaneously", &argparse.ParserConfig{
		WithColor: true,
		WithHint:  true,
		ContinueOnHelp: false,
	})

	GlobalOptions.ManifestPath = parser.String("m", "manifest", &argparse.Option{
		Default: ".mg.yml",
		Help: "Path to the yaml file with paths to repositories",
	})

	// GIT subcommand
	git_pass := parser.AddCommand("git", "Run an arbitrary git command in all repositories", &argparse.ParserConfig{
		WithColor:   true,
		DisableHelp: true,
		WithHint:    false,
	})
	git_args := git_pass.Strings("", "arg", &argparse.Option{Positional: true, Help: "Git arguments", Required: true})

	// SH subcommand
	sh_pass := parser.AddCommand("sh", "Run an arbitrary shell command in all repositories", &argparse.ParserConfig{
		WithColor:   true,
		DisableHelp: true,
	})
	sh_bin := sh_pass.String("", "bin", &argparse.Option{Positional: true, Help: "Binary to run in all repos", Required: true})
	sh_args := sh_pass.Strings("", "arg", &argparse.Option{Positional: true, Help: "Arguments passed to binary"})

	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}

	manifest, err := readManifest(*GlobalOptions.ManifestPath)

	if err != nil {
		fmt.Printf("Error reading manifest file %s: %s\n", *GlobalOptions.ManifestPath, err)
		os.Exit(1)
	}

	/* EXECUTE SUBCOMMANDS */
	switch {
	case git_pass.Invoked:
		println("Git args: ", strings.Join(*git_args, ", "))
		RunCmd("git", *git_args, manifest.paths())
	case sh_pass.Invoked:
		println("Shell args: ", *sh_bin, strings.Join(*sh_args, ", "))
	}
}
