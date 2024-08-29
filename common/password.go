package common

import (
	"os"

	"github.com/urfave/cli/v2"
)

func NewPasswordFlag(output *[]byte) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "password",
		Aliases: []string{"P"},
		Usage:   "set archive `PASSWORD`",
		Action: func(ctx *cli.Context, s string) error {
			*output = []byte(s)
			return nil
		},
	}
}

func NewPasswordFileFlag(output *[]byte) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "password-file",
		Aliases: []string{"Pf"},
		Usage:   "set archive password from `FILE`",
		Action: func(ctx *cli.Context, s string) error {
			data, err := os.ReadFile(s)
			if err != nil {
				return NewContextError("cannot read password from file", err)
			}
			*output = data
			return nil
		},
	}
}

func PasswordTestTemplate() string {
	return "This is test string for test password for archive. X-bit type."
}
