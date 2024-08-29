package main

import (
	"fmt"
	"io"
	"log"
	"mksea/common"
	"mksea/input"
	"mksea/output"
	"os"
	"path/filepath"

	"github.com/ncruces/zenity"
	"github.com/urfave/cli/v2"
)

func init() {
	if targetPath, err := os.Getwd(); err != nil {
		log.Fatalf("cannot obtain work dir: %v", err)
	} else {
		output.Env.WorkDir = targetPath
	}
}

func new_cli_app() *cli.App {
	input.Env.MaxMem = 0
	var password []byte
	app := &cli.App{
		Usage:                  "Self-Extractable Archive",
		Version:                fmt.Sprintf("v%d.%d.%d%s", Version.Major, Version.Minor, Version.Patch, Version.Suffix),
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "target output `DIRECTORY` (current directory by default)",
				Action: func(ctx *cli.Context, s string) error {
					newTargetPath := s
					startPoint := newTargetPath
					for {
						if info, err := os.Stat(startPoint); err != nil {
							startPoint = filepath.Dir(startPoint)
						} else if !info.IsDir() {
							return fmt.Errorf("path \"%s\" is not a directory", startPoint)
						} else {
							break
						}
					}
					if err := os.MkdirAll(newTargetPath, 0755); err != nil {
						return common.NewContextError("cannot create output directory", err)
					}
					output.Env.WorkDir = newTargetPath
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "max",
				Aliases: []string{"m"},
				Usage:   "set memory `LIMIT` (available suffixes KMGT) (8G by default)",
				Action: func(ctx *cli.Context, s string) error {
					newMemoryLimit := uint64(0)
					r := []rune(s)
					i, l := 0, len(r)
					for ; i < l; i++ {
						if '0' <= r[i] && r[i] <= '9' {
							newMemoryLimit = newMemoryLimit*10 + uint64(r[i]-'0')
						} else {
							break
						}
					}
					if i == 0 {
						return fmt.Errorf("invalid memory limit value %q", s)
					}
					if i < l {
						valid := false
						for _, it := range "KMGT" {
							newMemoryLimit *= 1024
							if it == r[i] {
								valid = true
								break
							}
						}
						if !valid {
							return fmt.Errorf("unsupported memory limit suffix %q", r[i:])
						}
					}
					input.Env.MaxMem = newMemoryLimit
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "silent",
				Aliases: []string{"s"},
				Usage:   "disable logging while unpacking",
			},
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "print archive content",
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := testPassword(password); err == zenity.ErrCanceled {
				return nil
			} else if err != nil {
				return err
			}

			if input.Env.MaxMem == 0 {
				input.Env.MaxMem = 8 * 1024 * 1024 * 1024 // 8G by default
			}

			exe, exeOffset, err := archiveOffset()
			if err != nil {
				return common.NewContextError("cannot obtain info from sea file", err)
			}

			in, err := input.Open(exe, exeOffset)
			if err != nil {
				return common.NewContextError("cannot open sea file", err)
			}
			defer in.Close()

			if ctx.Bool("list") {
				for it, err := in.Next(); it != nil || err != nil; it, err = in.Next() {
					if err != nil {
						return common.NewContextError("cannot go to next file", err)
					}
					fmt.Println(filepath.FromSlash(it.Path))
					if _, err := io.Copy(io.Discard, it.Reader); err != nil {
						return common.NewContextError("cannot read file", err)
					}
				}
				return nil
			}

			progress := func() int {
				p := in.Progress()
				return int(100.0 * float64(p.Current()) / float64(p.All()))
			}

			verboseMode := !ctx.Bool("silent")
			logf := func(format string, args ...any) {
				if verboseMode {
					log.Printf(format, args...)
				}
			}
			logln := func(args ...any) {
				if verboseMode {
					log.Println(args...)
				}
			}

			logln("[  0%] preparing...")
			for it, err := in.Next(); it != nil || err != nil; it, err = in.Next() {
				if err != nil {
					return common.NewContextError("cannot go to next file", err)
				}
				logf(
					"[%3d%%] unpacking \"%s\"...\n",
					progress(),
					filepath.FromSlash(it.Path),
				)
				outFile, err := output.OpenRaw(it.Path, 0755)
				if err != nil {
					return common.NewContextError("cannot open file for write", err)
				}

				if _, err := io.Copy(outFile, it.Reader); err != nil {
					outFile.Close()
					return common.NewContextError("cannot unpack file", err)
				}

				outFile.Close()
			}
			logln("[100%] done!")
			return nil
		},
	}
	if len(input.Env.PasswordTest) > 0 {
		app.Flags = append(app.Flags,
			common.NewPasswordFlag(&password),
			common.NewPasswordFileFlag(&password),
		)
	}
	return app
}

func main() {
	if main_gui() {
		return
	}
	if err := new_cli_app().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
