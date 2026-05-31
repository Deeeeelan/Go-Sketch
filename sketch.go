package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/urfave/cli/v3"
)

var BOARD_BORDERS = [3][3]string{
	{"╭", "─", "╮"},
	{"│", "%", "│"},
	{"╰", "─", "╯"}}

func makeBoard(inner_width int, inner_height int) string { // TODO: Add outer border
	board := BOARD_BORDERS[0][0]
	board += strings.Repeat(BOARD_BORDERS[0][1], inner_width)
	board += BOARD_BORDERS[0][2] + "\n"
	board += strings.Repeat(BOARD_BORDERS[1][0]+strings.Repeat(BOARD_BORDERS[1][1], inner_width)+BOARD_BORDERS[1][2]+"\n", inner_height)
	board += BOARD_BORDERS[2][0]
	board += strings.Repeat(BOARD_BORDERS[2][1], inner_width)
	board += BOARD_BORDERS[2][2] + "\n"
	return board
}

type model struct {
	cursor [2]int
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor[1] > 0 {
				m.cursor[1]++
			}

		}
	}

	return m, nil
}

func (m model) View() tea.View {
	s := "Sketch\n\n"

	s += "---\nBoard\n---\n"

	s += "\nPress q to quit.\n"

	return tea.NewView(s)
}

func main() {
	cmd := &cli.Command{
		Name:      "go-sketch",
		Usage:     "Create a sketch board",
		Suggest:   true,
		ArgsUsage: "width height",
		Arguments: []cli.Argument{
			&cli.IntArg{ // TODO: Add error for missing args
				Name: "width",
			},
			&cli.IntArg{
				Name: "height",
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Printf("Dims:  %d %d\n", cmd.IntArg("width"), cmd.IntArg("height"))
			fmt.Print(makeBoard(cmd.IntArg("width"), cmd.IntArg("height")))
			p := tea.NewProgram(initialModel())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
