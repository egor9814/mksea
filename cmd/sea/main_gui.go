//go:build fyne_gui

package main

import (
	"fmt"
	"io"
	"log"
	"mksea/common"
	"mksea/input"
	"mksea/output"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/urfave/cli/v2"
)

type PagerControlButtons struct {
	NextButton *widget.Button
	PrevButton *widget.Button
}

type Page struct {
	Title   string
	Content func(buttons *PagerControlButtons) fyne.CanvasObject
}

type Pages []Page

type PagerLayout struct {
}

func (l *PagerLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSquareSize(0)
}

func (l *PagerLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for _, it := range objects {
		it.Resize(containerSize)
		it.Move(pos)

		pos = pos.AddXY(containerSize.Width, 0)
	}
}

type Pager struct {
	widget.BaseWidget
	PagerControlButtons
	PageChanged func(oldIndex, newIndex int)
	pages       Pages
	index       int
	l           PagerLayout
}

func NewPager(pages Pages) *Pager {
	// TODO: fyne i18n!
	pager := &Pager{
		pages: pages,
	}
	pager.ExtendBaseWidget(pager)
	pager.NextButton = widget.NewButton("Next", func() {
		pager.Next()
	})
	pager.PrevButton = widget.NewButton("Previous", func() {
		pager.Prev()
	})
	pager.PrevButton.Disable()
	if len(pages) == 0 {
		pager.NextButton.Disable()
	}
	return pager
}

func (p *Pager) CreateRenderer() fyne.WidgetRenderer {
	objects := make([]fyne.CanvasObject, 0, len(p.pages))
	for _, it := range p.pages {
		if it.Content == nil {
			continue
		}
		title := widget.NewLabel(it.Title)
		title.TextStyle.Bold = true
		objects = append(objects, container.New(
			layout.NewVBoxLayout(),
			title,
			widget.NewSeparator(),
			container.NewPadded(
				it.Content(&p.PagerControlButtons),
			),
			layout.NewSpacer(),
		))
	}
	c := container.New(&p.l, objects...)
	return widget.NewSimpleRenderer(c)
}

func (p *Pager) ScrollTo(index int) {
	if p.index == index {
		return
	}
	oldIndex := p.index
	size := p.Size()
	move := canvas.NewPositionAnimation(
		fyne.NewPos(-size.Width*float32(oldIndex), 0),
		fyne.NewPos(-size.Width*float32(index), 0),
		canvas.DurationStandard,
		func(pos fyne.Position) {
			p.Move(pos)
		},
	)
	move.Curve = fyne.AnimationEaseInOut
	p.index = index
	move.Start()
	if p.index == len(p.pages)-1 {
		p.NextButton.Disable()
	} else {
		p.NextButton.Enable()
	}
	if p.index == 0 {
		p.PrevButton.Disable()
	} else {
		p.PrevButton.Enable()
	}
	if p.PageChanged != nil {
		p.PageChanged(oldIndex, index)
	}
}

func (p *Pager) Next() {
	p.ScrollTo(min(p.index+1, len(p.pages)-1))
}

func (p *Pager) Prev() {
	p.ScrollTo(max(p.index-1, 0))
}

func UnpackWindow(a fyne.App) fyne.Window {
	baseName := metaInfo.Name

	w := a.NewWindow(fmt.Sprintf("%s Auto Unpacker", baseName))
	w.SetMaster()

	welcomePage := Page{
		Title: fmt.Sprintf("Welcome to the %s Auto Unpacker!", baseName),
		Content: func(buttons *PagerControlButtons) fyne.CanvasObject {
			return container.NewVBox(
				layout.NewSpacer(),
				func() fyne.CanvasObject {
					l := widget.NewRichText(
						&widget.TextSegment{
							Text: fmt.Sprintf(
								"The Auto Unpacker will unpack %s on your computer. Click ",
								baseName,
							),
							Style: widget.RichTextStyle{
								Inline: true,
							},
						},
						&widget.TextSegment{
							Text: "Next",
							Style: widget.RichTextStyle{
								TextStyle: fyne.TextStyle{
									Bold: true,
								},
								Inline: true,
							},
						},
						&widget.TextSegment{
							Text: " to continue or ",
							Style: widget.RichTextStyle{
								Inline: true,
							},
						},
						&widget.TextSegment{
							Text: "Cancel",
							Style: widget.RichTextStyle{
								TextStyle: fyne.TextStyle{
									Bold: true,
								},
								Inline: true,
							},
						},
						&widget.TextSegment{
							Text: " to exit the Installer.",
							Style: widget.RichTextStyle{
								Inline: true,
							},
						},
					)
					l.Wrapping = fyne.TextWrapWord
					return l
				}(),
				layout.NewSpacer(),
			)
		},
	}

	availableMemoryLimit := []string{
		"1G",
		"2G",
		"4G",
		"8G",
	}
	currentMemoryLimit := availableMemoryLimit[len(availableMemoryLimit)-1]
	memoryLimitChooser := widget.NewRadioGroup(
		availableMemoryLimit,
		func(s string) {
			currentMemoryLimit = s
		},
	)
	memoryLimitChooser.Selected = currentMemoryLimit
	memoryLimitChooser.Required = true
	memoryConfigurationPage := Page{
		Title: "Configuration",
		Content: func(buttons *PagerControlButtons) fyne.CanvasObject {
			return container.NewVBox(
				container.NewPadded(
					container.NewVBox(
						widget.NewLabel("Choose memory limit:"),
						memoryLimitChooser,
					),
				),
			)
		},
	}

	targetPath := filepath.Join(output.Env.WorkDir, metaInfo.Name)
	unpackTargetPathSegment := &widget.TextSegment{
		Text: targetPath,
		Style: widget.RichTextStyle{
			TextStyle: fyne.TextStyle{
				Bold: true,
			},
			Inline: true,
		},
	}
	summaryPageText := widget.NewRichText(
		&widget.TextSegment{
			Text: "Click ",
			Style: widget.RichTextStyle{
				Inline: true,
			},
		},
		&widget.TextSegment{
			Text: "Next",
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{
					Bold: true,
				},
				Inline: true,
			},
		},
		&widget.TextSegment{
			Text: fmt.Sprintf(
				" to unpack %s to ",
				baseName,
			),
			Style: widget.RichTextStyle{
				Inline: true,
			},
		},
		unpackTargetPathSegment,
		&widget.TextSegment{
			Text: ".",
			Style: widget.RichTextStyle{
				Inline: true,
			},
		},
	)
	summaryPageText.Wrapping = fyne.TextWrapWord
	unpackPathPage := Page{
		Title: "Unpack path",
		Content: func(buttons *PagerControlButtons) fyne.CanvasObject {
			path := widget.NewEntry()
			path.Text = filepath.FromSlash(targetPath)
			return container.NewVBox(
				widget.NewLabel("Choose installation path:"),
				path,
				container.NewHBox(
					layout.NewSpacer(),
					widget.NewButton("Open", func() {
						fo := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
							if err != nil {
								log.Println(err)
								return
							}
							if lu == nil {
								return
							}
							path.Text = filepath.FromSlash(lu.String()[7:])
							targetPath = path.Text
							unpackTargetPathSegment.Text = targetPath
							summaryPageText.Refresh()
							path.Refresh()
						}, w)
						startPoint := filepath.ToSlash(path.Text)
						for len(startPoint) > 0 {
							if info, err := os.Stat(startPoint); err == nil && info.IsDir() {
								break
							} else {
								startPoint = filepath.Dir(startPoint)
							}
						}
						if len(startPoint) == 0 {
							startPoint, _ = os.UserHomeDir()
						}
						uri, err := storage.ParseURI("file://" + filepath.ToSlash(startPoint))
						if err != nil {
							log.Fatal(err)
						}
						uril, err := storage.ListerForURI(uri)
						if err != nil {
							log.Fatal(err)
						}
						fo.SetLocation(uril)
						fo.Show()
					}),
				),
			)
		},
	}

	summaryPage := Page{
		Title: "Summary",
		Content: func(buttons *PagerControlButtons) fyne.CanvasObject {
			return container.NewVBox(
				summaryPageText,
			)
		},
	}

	extractingItem := binding.NewString()
	extractingProgress := binding.NewFloat()
	extractingProgressBar := widget.NewProgressBarWithData(extractingProgress)
	installPage := Page{
		Title: "Unpacking...",
		Content: func(buttons *PagerControlButtons) fyne.CanvasObject {
			return container.NewVBox(
				extractingProgressBar,
				widget.NewLabelWithData(binding.NewSprintf("Unpacking %s...", extractingItem)),
			)
		},
	}
	setupMtx := sync.Mutex{}
	setupPause := false
	setupCancel := false
	setupErrChan := make(chan error, 1)
	startSetup := func() error {
		func() {
			input.Env.MaxMem = 1
			for _, it := range availableMemoryLimit {
				if it == currentMemoryLimit {
					input.Env.MaxMem *= 1024 * 1024 * 1024
					return
				}
				input.Env.MaxMem *= 2
			}
			panic("unreachable")
		}()

		output.Env.WorkDir = targetPath

		exe, exeOffset, err := archiveOffset()
		if err != nil {
			return fmt.Errorf("cannot obtain info from sea file: %v", err)
		}

		in, err := input.Open(exe, exeOffset)
		if err != nil {
			return err
		}
		extractingProgressBar.Max = float64(len(metaInfo.Files))

		go func() {
			i := 0
			for {
				setupMtx.Lock()
				c := setupCancel
				p := setupPause
				setupMtx.Unlock()
				if c {
					return
				}
				if p {
					time.Sleep(time.Second)
					continue
				}
				i++
				extractingProgress.Set(float64(i))
				it, err := in.Next()
				if err != nil {
					in.Close()
					setupErrChan <- err
					return
				}
				if it == nil {
					break
				}
				extractingItem.Set(it.Path)

				outFile, err := output.OpenRaw(it.Path, 0755)
				if err != nil {
					in.Close()
					setupErrChan <- err
					return
				}

				if _, err := io.Copy(outFile, it.Reader); err != nil {
					outFile.Close()
					in.Close()
					setupErrChan <- err
					return
				}

				outFile.Close()
			}
			setupErrChan <- nil
		}()

		if err := <-setupErrChan; err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				return err
			}
		}

		in.Close()

		return nil
	}

	installError := binding.NewString()
	finishPage := Page{
		Title: "Finish",
		Content: func(buttons *PagerControlButtons) fyne.CanvasObject {
			return widget.NewLabel("Unpacking successful!")
		},
	}
	finishErrorPage := Page{
		Title: "Unpacking error",
		Content: func(buttons *PagerControlButtons) fyne.CanvasObject {
			return widget.NewLabelWithData(installError)
		},
	}

	installDone := false
	confirmCancel := func() {
		if installDone {
			w.Close()
			return
		}
		setupMtx.Lock()
		setupPause = true
		setupMtx.Unlock()
		dialog.ShowConfirm("Cancel unpacking", "Are you sure?", func(b bool) {
			if b {
				setupMtx.Lock()
				setupCancel = true
				setupMtx.Unlock()
				setupErrChan <- &exec.ExitError{}
				w.Close()
			} else {
				setupMtx.Lock()
				setupPause = false
				setupMtx.Unlock()
			}
		}, w)
	}
	cancelBtn := widget.NewButton("Cancel", func() {
		confirmCancel()
	})
	exitBtn := widget.NewButton("Quit", func() {
		w.Close()
	})
	exitBtn.Hide()
	p := NewPager(Pages{
		welcomePage,
		memoryConfigurationPage,
		unpackPathPage,
		summaryPage,
		installPage,
		finishPage,
		finishErrorPage,
	})
	p.PageChanged = func(oldIndex, newIndex int) {
		if newIndex == 1 {
			p.PrevButton.Disable()
		} else if newIndex >= 4 {
			p.PrevButton.Disable()
			p.NextButton.Disable()
			if newIndex == 4 {
				err := startSetup()
				p.Next()
				cancelBtn.Disable()
				if err != nil {
					installError.Set(err.Error())
					p.Next()
				}
			} else {
				cancelBtn.Hide()
				exitBtn.Show()
				installDone = true
			}
		}
	}
	buttonsDivider := container.NewCenter()
	buttonsDivider.Resize(fyne.NewSquareSize(32))
	buttons := container.NewHBox(layout.NewSpacer(), p.PrevButton, p.NextButton, buttonsDivider, cancelBtn, exitBtn)

	content := container.NewVBox(p, layout.NewSpacer(), buttons)

	w.SetContent(content)
	w.Resize(fyne.NewSize(480, 320))
	w.SetFixedSize(true)
	w.SetCloseIntercept(func() {
		confirmCancel()
	})

	return w
}

func PasswordWindow(a fyne.App, onAccept func()) {
	baseTitle := "Type Archive Password"
	attemtsTitle := baseTitle + ". Attempts left: "

	w := a.NewWindow(baseTitle)
	w.Resize(fyne.NewSize(480, 60))

	attempts := passwordAttempts
	check := func(text string) {
		if attempts == 0 {
			a.Quit()
		}
		password := []byte(text)
		if decodeEncoderKey(password) {
			onAccept()
			w.Close()
			return
		}
		attempts--
		w.SetTitle(attemtsTitle + strconv.Itoa(attempts))
	}

	entry := widget.NewPasswordEntry()
	entry.OnSubmitted = func(s string) {
		check(s)
	}

	w.SetContent(container.NewVBox(
		widget.NewForm(
			widget.NewFormItem(
				"Password:",
				entry,
			),
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Ok", func() {
				check(entry.Text)
			}),
			widget.NewButton("Cancel", a.Quit),
		),
	))

	w.Show()
}

func run_gui() {
	a := app.New()
	PasswordWindow(a, func() {
		UnpackWindow(a).Show()
	})
	a.Run()
}

func main_gui() bool {
	if len(os.Args) > 1 {
		var password []byte
		passwordApp := &cli.App{
			Usage: metaInfo.Name + " - Auto Unpacker",
			Flags: []cli.Flag{
				common.NewPasswordFlag(&password),
				common.NewPasswordFileFlag(&password),
			},
		}
		if err := passwordApp.Run(os.Args); err != nil {
			log.Fatal(err)
		}
		if !decodeEncoderKey(password) {
			log.Fatal(errInvPass)
		}
	}
	run_gui()
	return true
}
