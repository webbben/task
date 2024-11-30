package taskui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/webbben/task/internal/tasks"
	"github.com/webbben/task/internal/types"
	listcomponent "github.com/webbben/task/internal/ui/components/list-component"
	noteviewer "github.com/webbben/task/internal/ui/components/note-viewer"
)

type model struct {
	content        *types.Task
	noteViewer     *noteviewer.NoteViewerModel
	noteViewerOpen bool
	noteList       listcomponent.ListComponentModel
}

type noteListItem struct {
	title   string
	content string
}

func (item noteListItem) FilterValue() string {
	return item.title
}

func (item noteListItem) Title() string {
	return item.title
}

func (item noteListItem) Description() string {
	if len(item.content) < 15 {
		return item.content
	}
	return item.content[:15] + "..."
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// note viewer
	newNoteViewerModel, cmd := m.noteViewer.Update(msg)
	m.noteViewer = newNoteViewerModel
	cmds = append(cmds, cmd)

	// note list
	newNoteListModel, cmd := m.noteList.Update(msg)
	m.noteList = newNoteListModel
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.noteViewerOpen {
		return m.noteViewer.View()
	}

	return m.noteList.View()
}

func (m *model) setSelectedNote(noteTitle string) {
	noteContent, exists := m.content.Notes[noteTitle]
	if !exists {
		log.Println("note not found")
		return
	}
	if noteContent == "" {
		log.Println("selected note is empty?")
		return
	}
	m.noteViewer.SetNoteContent(noteTitle, noteContent)
	m.noteViewerOpen = true
}

func (m *model) onSelectNote(item list.Item) {
	noteItem, ok := item.(noteListItem)
	if !ok {
		log.Println("note list item can't be converted to list.Item")
		return
	}
	m.setSelectedNote(noteItem.title)
}

func (m *model) onCloseNote() {
	m.noteViewerOpen = false
}

func RunUI(taskID string) error {
	task, err := tasks.GetTask(taskID)
	if err != nil {
		return fmt.Errorf("failed to run task UI: %d", err)
	}

	m := &model{
		content: task,
	}

	// set up note list
	noteList := make([]list.Item, 0)
	for title, text := range task.Notes {
		noteList = append(noteList, noteListItem{
			title:   title,
			content: text,
		})
	}
	m.noteList = listcomponent.New(noteList, task.Title, 0, 0, m.onSelectNote)

	m.noteViewer = noteviewer.New("", "", m.onCloseNote)

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error occurred while viewing task: %d", err)
	}
	return nil
}
