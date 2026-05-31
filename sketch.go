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
	{"│", " ", "│"},
	{"╰", "─", "╯"}}

func makeBoard(m model, inner_width int, inner_height int, char_pos [2]int) string { // TODO: Add outer border
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
	char_pos       [2]int
	subchar_pos    [2]int // braille character dot
	canvas_width   int
	canvas_height  int
	canvas_content string
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
			if m.char_pos[1] > 0 {
				m.char_pos[1]--
			}
		case "down", "j":
			if m.char_pos[1] < m.canvas_height {
				m.char_pos[1]++
			}
		case "left", "h":
			if m.char_pos[0] > 0 {
				m.char_pos[0]--
			}
		case "right", "l":
			if m.char_pos[0] < m.canvas_width {
				m.char_pos[0]++
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	s := "Sketch\n\n"

	current_board := makeBoard(m, m.canvas_width, m.canvas_height)
	s += current_board

	s += fmt.Sprintf("cords: (%d, %d), [%d, %d]", m.char_pos[0], m.char_pos[1], m.subchar_pos[0], m.subchar_pos[1])
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
			new_model := initialModel()
			new_model.canvas_width = cmd.IntArg("width")
			new_model.canvas_height = cmd.IntArg("height")
			p := tea.NewProgram(new_model)

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
