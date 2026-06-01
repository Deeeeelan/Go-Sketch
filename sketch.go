package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/urfave/cli/v3"
)

var BOARD_BORDERS = [3][3]string{
	{"╭", "─", "╮"},
	{"│", " ", "│"},
	{"╰", "─", "╯"}}

func makeBoard(m model, inner_width int, inner_height int) string { // TODO: Add outer border
	board := BOARD_BORDERS[0][0]
	board += strings.Repeat(BOARD_BORDERS[0][1], inner_width)
	board += BOARD_BORDERS[0][2] + "\n"
	for i := 0; i < inner_height; i++ {
		board += BOARD_BORDERS[1][0]
		for j := 0; j < inner_width; j++ {
			if m.char_pos[0] == j && m.char_pos[1] == i {
				curr_rune := int(m.canvas_content[i][j])
				if curr_rune == ' ' {
					curr_rune = 0x2800
				}
				braille_cursor := curr_rune | int(math.Pow(2, float64(m.subchar_pos[1])))<<(3*m.subchar_pos[0])

				board += string(rune(braille_cursor))
			} else {
				board += string(m.canvas_content[i][j])
			}

		}
		board += BOARD_BORDERS[1][2] + "\n"
	}
	board += BOARD_BORDERS[2][0]
	board += strings.Repeat(BOARD_BORDERS[2][1], inner_width)
	board += BOARD_BORDERS[2][2] + "\n"
	return board
}

func clearCanvasContent(m *model) {
	for i := 0; i < m.canvas_height; i++ {
		m.canvas_content = append(m.canvas_content, []rune(strings.Repeat(" ", m.canvas_width)))
	}
}

func updateCanvas(m *model) {
	curr_rune := m.canvas_content[m.char_pos[1]][m.char_pos[0]]
	if curr_rune == ' ' {
		curr_rune = 0x2800
	}
	new_rune := curr_rune | rune(int(math.Pow(2, float64(m.subchar_pos[1])))<<(3*m.subchar_pos[0]))
	m.canvas_content[m.char_pos[1]][m.char_pos[0]] = new_rune
}

type model struct {
	char_pos       [2]int // (x, y)
	subchar_pos    [2]int // braille character dot
	canvas_width   int
	canvas_height  int
	canvas_content [][]rune
}

func initialModel(cw int, ch int) model {
	m := model{}
	m.canvas_width = cw
	m.canvas_height = ch

	clearCanvasContent(&m)
	return m
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
		case "esc":
			clearCanvasContent(&m)

		case "up", "k":
			if m.subchar_pos[1] > 0 {
				m.subchar_pos[1]--
			} else if m.char_pos[1] > 0 {
				m.subchar_pos[1] = 2
				m.char_pos[1]--
			}

		case "down", "j":
			if m.subchar_pos[1] < 2 {
				m.subchar_pos[1]++
			} else if m.char_pos[1] < m.canvas_height-1 {
				m.subchar_pos[1] = 0
				m.char_pos[1]++
			}

		case "left", "h":
			if m.subchar_pos[0] > 0 {
				m.subchar_pos[0]--
			} else if m.char_pos[0] > 0 {
				m.subchar_pos[0] = 1
				m.char_pos[0]--
			}
		case "right", "l":
			if m.subchar_pos[0] < 1 {
				m.subchar_pos[0]++
			} else if m.char_pos[0] < m.canvas_width-1 {
				m.subchar_pos[0] = 0
				m.char_pos[0]++
			}
		}
		updateCanvas(&m)
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
			new_model := initialModel(cmd.IntArg("width"), cmd.IntArg("height"))
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
