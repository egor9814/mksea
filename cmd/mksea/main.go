package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"mksea/common"
	"mksea/output"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func new_cli_app() *cli.App {
	var excludedFiles cli.StringSlice
	packer := Packer{
		CompressFormat: CompressZstdHigh,
		Platforms:      make([]TargetPlatform, 0, 2),
	}
	return &cli.App{
		Name:                   "mksea",
		Usage:                  "MaKe Self-Extractable Archive",
		Version:                fmt.Sprintf("v%d.%d.%d%s", Version.Major, Version.Minor, Version.Patch, Version.Suffix),
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "base `NAME` of output executable",
				Action: func(ctx *cli.Context, s string) error {
					packer.Meta.Name = s
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "target output `DIRECTORY` (current directory by default)",
				Action: func(ctx *cli.Context, s string) error {
					packer.Output = s
					return nil
				},
			},
			&cli.StringSliceFlag{
				Name:        "exclude",
				Aliases:     []string{"E"},
				Usage:       "exclude `FILE` from packing",
				Destination: &excludedFiles,
			},
			&cli.StringSliceFlag{
				Name:    "platform",
				Aliases: []string{"p"},
				Usage:   "add target `PLATFROM`",
				Action: func(ctx *cli.Context, s []string) error {
					var platform TargetPlatform
					for _, it := range s {
						platform = TargetNone
						platform.FromString(it)
						if platform.IsValid() {
							packer.Platforms = append(packer.Platforms, platform)
						} else {
							return fmt.Errorf("unsupported platfrom %q", it)
						}
					}
					return nil
				},
			},
			common.NewPasswordFlag(&packer.Password),
			common.NewPasswordFileFlag(&packer.Password),
			&cli.BoolFlag{
				Name:    "encrypt",
				Aliases: []string{"e"},
				Usage:   "encrypt archive",
				Action: func(ctx *cli.Context, b bool) error {
					packer.Encrypt = b
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "compress-zstd",
				Aliases: []string{"z"},
				Usage:   "set compress `LEVEL` for zstd algorythm (low, mid, high)",
				Value:   "high",
				Action: func(ctx *cli.Context, s string) error {
					if !packer.CompressFormat.ZstdFromString(s) {
						return fmt.Errorf("unsupported compress level %q", s)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "compress-xz",
				Aliases: []string{"xz"},
				Usage:   "set compress `LEVEL` for zstd algorythm (low, mid, high)",
				Action: func(ctx *cli.Context, b bool) error {
					packer.CompressFormat = CompressXz
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "no-compress",
				Aliases: []string{"nc"},
				Usage:   "disable compressing",
				Action: func(ctx *cli.Context, b bool) error {
					packer.CompressFormat = CompressNone
					return nil
				},
			},
			&cli.Uint64Flag{
				Name:    "zstd-threads",
				Aliases: []string{"zt"},
				Usage:   "set threads `COUNT` for zstd algorythm (0 means count of cores)",
				Value:   0,
				Action: func(ctx *cli.Context, u uint64) error {
					packer.ZstdThreads = int(u)
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "gui",
				Aliases: []string{"g"},
				Usage:   "make installer gui (fyne-cross required)",
				Action: func(ctx *cli.Context, b bool) error {
					packer.Gui = b
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "silent",
				Aliases: []string{"s"},
				Usage:   "supress log messages",
				Action: func(ctx *cli.Context, b bool) error {
					packer.Silent = b
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			defer func() {
				if err := cleanup(); err != nil {
					packer.logf("cleanup failed: %v\n", err)
				}
			}()
			if err := init_wd(); err != nil {
				return common.NewContextError("init work dir failed", err)
			}

			excludedFileSet := NewFileSet()
			for _, it := range excludedFiles.Value() {
				excludedFileSet.Resolve(it, nil)
			}
			includedFileSet := NewFileSet()
			includeResolver := func(p string, _ fs.FileInfo) bool {
				return !strings.HasPrefix(p, workInstallerDir)
			}
			if ctx.Args().Len() == 0 {
				includedFileSet.Resolve("*", includeResolver)
			} else {
				for _, it := range ctx.Args().Slice() {
					includedFileSet.Resolve(it, includeResolver)
				}
			}
			includedFileSet.Remove(excludedFileSet)
			if includedFileSet.Len() == 0 {
				return errors.New("no input files")
			}
			packer.Input = includedFileSet.List()
			output.Env.WorkDir = workDir
			return packer.Pack()
		},
	}
}

func main() {
	if err := new_cli_app().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
