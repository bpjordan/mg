package main

import (
	"log"
	"strings"

	"github.com/hellflame/argparse"
)

func main() {

	/* DEFINE & PARSE ARGUMENTS */
	parser := argparse.NewParser("mg", "Utility to manage several git repositories simultaneously", &argparse.ParserConfig{
		WithColor: true,
		WithHint:  true,
	})

	manifestPath := parser.String("m", "manifest", &argparse.Option{
		Default: ".mg.yaml",
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

	parser.Parse(nil)

	_, err := readManifest(*manifestPath)

	if err != nil {
		log.Fatal("Failed to read manifest: ", err)
	}

	/* EXECUTE SUBCOMMANDS */
	switch {
	case git_pass.Invoked:
		println("Git args: ", strings.Join(*git_args, ", "))
	case sh_pass.Invoked:
		println("Shell args: ", *sh_bin, strings.Join(*sh_args, ", "))
	}
}
