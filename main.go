package main

import (
	"strings"

	"github.com/hellflame/argparse"
)

func main() {
    parser := argparse.NewParser("mg", "Utility to manage several git repositories simultaneously", &argparse.ParserConfig{
	WithColor: true,
	WithHint: true,
    })

    git_pass := parser.AddCommand("git", "Run an arbitrary git command in all repositories", &argparse.ParserConfig{})

    git_args := git_pass.Strings("a", "", &argparse.Option{Positional: true, Help: "Git arguments", Required: true})

    sh_pass := parser.AddCommand("sh", "Run an arbitrary command in all repositories", &argparse.ParserConfig{})

    sh_args := sh_pass.Strings("a", "", &argparse.Option{Positional: true, Help: "Shell command", Required: true})

    parser.Parse(nil)

    if git_args != nil {
	println("Git args: ", strings.Join(*git_args, ", "))
    }

    if sh_args != nil {
	println("Shell args: ", strings.Join(*sh_args, ", "))
    }
}
