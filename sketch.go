package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/urfave/cli/v3"
)

const BOARD = `╭─╮
│%│
╰─╯`

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
		Name:  "go-sketch",
		Usage: "Create a sketch board",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("It works!")
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
