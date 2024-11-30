package listcomponent

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ListComponentModel struct {
	list            list.Model
	onEnterCallback func(item list.Item)
}

func (m ListComponentModel) Init() tea.Cmd {
	return nil
}

func (m ListComponentModel) onEnter() {
	i := m.list.Index()
	listItems := m.list.Items()
	if i < 0 || i > len(listItems)-1 {
		log.Println("invalid index or list of items")
		return
	}
	m.onEnterCallback(listItems[i])
}

func (m ListComponentModel) Update(msg tea.Msg) (ListComponentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			m.onEnter()
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListComponentModel) View() string {
	return docStyle.Render(m.list.View())
}

func New(items []list.Item, title string, width, height int, onEnterFn func(item list.Item)) ListComponentModel {

	m := ListComponentModel{list: list.New(items, list.NewDefaultDelegate(), width, height)}
	m.list.Title = title
	m.list.DisableQuitKeybindings()
	m.onEnterCallback = onEnterFn

	return m
}
