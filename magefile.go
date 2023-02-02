//go:build mage

package main

import (
	"os"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Clear() error {
	return sh.Run("rm", "fpc", "-f")
}

func Build() error {
	mg.Deps(Clear)

	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	env := map[string]string{
		"GOOS":   runtime.GOOS,
		"GOARCH": runtime.GOARCH,
	}
	_, err := sh.Exec(env, os.Stdout, os.Stderr, "go", "build", "-v", "-ldflags="+"-w -s", "-o", "fpc", "./cmd/cli")

	return err
}
