//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const cliBuildDst = "./build/cli/"
const lambdaBuildDst = "./build/lambdas/"

func clearCli() error {
	return sh.Run("rm", "-rf", cliBuildDst)
}

func clearLambdas() error {
	return sh.Run("rm", "-rf", lambdaBuildDst)
}

// Clears contents of ./build directory.
func Clear() error {
	mg.Deps(clearCli)
	mg.Deps(clearLambdas)

	return nil
}

// Builds the cli version of the app.
func Cli() error {
	mg.Deps(clearCli)

	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	env := map[string]string{
		"GOOS":        runtime.GOOS,
		"GOARCH":      runtime.GOARCH,
		"CGO_ENABLED": "0",
	}

	cliBuildDst := filepath.Join(cliBuildDst, "fpc")

	cliSrc := "./cmd/cli"

	_, err := sh.Exec(env, os.Stdout, os.Stderr, "go", "build", "-v", "-ldflags="+"-w -s", "-o", cliBuildDst, cliSrc)

	return err
}

func lambda(lambdaName string) error {
	env := map[string]string{
		"GOOS":        "linux",
		"GOARCH":      "amd64",
		"CGO_ENABLED": "0",
	}

	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	lambdaBinDst := filepath.Join(lambdaBuildDst, lambdaName)
	zipDst := filepath.Join(lambdaBuildDst, fmt.Sprintf("%s.zip", lambdaName))

	lambdaSrc := fmt.Sprintf("./cmd/lambdas/%s", lambdaName)

	if _, err := sh.Exec(env, os.Stdout, os.Stderr, "go", "build", "-v",
		"-o", lambdaBinDst, lambdaSrc); err != nil {
		return err
	}

	if _, err := sh.Exec(env, os.Stdout, os.Stderr, "zip", "-j",
		zipDst, lambdaBinDst); err != nil {
		return err
	}

	return nil
}

func lambdaFetch() error {
	return lambda("fetch")
}

func lambdaGenerate() error {
	return lambda("generate")
}

func lambdaGet() error {
	return lambda("get")
}

// Builds lambda-based version of the app.
func Lambdas() error {
	mg.Deps(clearLambdas)

	mg.Deps(lambdaFetch)
	mg.Deps(lambdaGenerate)
	mg.Deps(lambdaGet)

	return nil
}

// Runs all tests.
func Test() error {
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "go", "test", "./...", "-cover")
	return err
}
