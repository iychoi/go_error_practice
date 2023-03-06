package main

import (
	"errors"
	"fmt"
	_ "net/http/pprof"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

var (
	AppErr error = xerrors.New("app crashed")
)

func main() {
	log.SetOutput(os.Stdout)
	test_basic()
	fmt.Printf("\n\n")
	test_custom_error()
}

func test_basic() {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "test_basic",
	})

	// err must be an instance of error
	// calling layer1_2 --> layer1_1 --> AppErr
	err := layer1_2()

	// test logger
	logger.Errorf("error print1\n%v", err)
	logger.Errorf("")
	logger.Errorf("error print2\n%+v", err)
	logger.Errorf("")

	fmt.Printf("error print1\n%v\n", err)
	fmt.Printf("\n")
	fmt.Printf("error print2\n%+v\n", err)
	fmt.Printf("\n")

	// this doesn't work anymore
	if err == AppErr {
		fmt.Printf("err == AppErr\n")
	} else {
		fmt.Printf("err != AppErr\n")
	}

	// but this works
	if errors.Is(err, AppErr) {
		fmt.Printf("err is AppErr\n")
	} else {
		fmt.Printf("err is not AppErr\n")
	}
}

func test_custom_error() {
	// err must be an instance of error
	err := layer2_2()

	fmt.Printf("error print1\n%v\n", err)
	fmt.Printf("\n")
	fmt.Printf("error print2\n%+v\n", err)
	fmt.Printf("\n")

	// but this works
	if errors.Is(err, &CustomError{}) {
		fmt.Printf("err is CustomError\n")
	} else {
		fmt.Printf("err is not CustomError\n")
	}
}

func layer1_2() error {
	err := layer1_1()
	if err != nil {
		return xerrors.Errorf("layer2 failed: %w", err)
	}
	return nil
}

func layer1_1() error {
	err := raise_base_error()
	if err != nil {
		return xerrors.Errorf("layer1 failed: %w", err)
	}
	return nil
}

func raise_base_error() error {
	return AppErr
}

func layer2_2() error {
	err := layer2_1()
	if err != nil {
		return xerrors.Errorf("layer2 failed: %w", err)
	}
	return nil
}

func layer2_1() error {
	err := raise_custom_error()
	if err != nil {
		return xerrors.Errorf("layer1 failed: %w", err)
	}
	return nil
}

func raise_custom_error() error {
	return xerrors.Errorf("%w", NewCustomError("app crashed again!"))
}
