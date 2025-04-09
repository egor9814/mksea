package tui

import (
	"log"
	tui_common "mksea/cmd/sea/internal/common"
	"mksea/common"
	"mksea/input"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

func Main() {
	var password []byte
	metaInfo := tui_common.MetaInfo()
	cliApp := cli.App{
		Usage: metaInfo.Name + " - Auto Unpacker",
		Flags: []cli.Flag{
			common.NewPasswordFlag(&password),
			common.NewPasswordFileFlag(&password),
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "print archive content",
			},
		},
		Action: func(ctx *cli.Context) error {
			passwordProvided := len(password) > 0
			if passwordProvided {
				if !tui_common.DecodeEncoderKey(password) {
					return tui_common.ErrInvPass
				}
			}

			isList := ctx.Bool("list")
			var m tea.Model
			if !passwordProvided && len(input.Env.PasswordTest) != 0 {
				m = newPasswordModel(isList)
			} else if isList {
				m = newListModel()
			} else {
				m = newUnpackModel()
			}

			if _, err := tea.NewProgram(m).Run(); err != nil {
				return err
			}
			return nil
		},
	}
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
