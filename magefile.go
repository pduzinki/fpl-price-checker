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

func ClearCli() error {
	return sh.Run("rm", "-rf", cliBuildDst)
}

func ClearLambdas() error {
	return sh.Run("rm", "-rf", lambdaBuildDst)
}

func Cli() error {
	mg.Deps(ClearCli)

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

func LambdaFetch() error {
	return lambda("fetch")
}

func LambdaGenerate() error {
	return lambda("generate")
}

func LambdaGet() error {
	return lambda("get")
}

func Lambdas() error {
	mg.Deps(ClearLambdas)

	mg.Deps(LambdaFetch)
	mg.Deps(LambdaGenerate)
	mg.Deps(LambdaGet)

	return nil
}

func Test() error {
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "go", "test", "./...", "-cover")
	return err
}
