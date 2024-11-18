package taskui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/webbben/task/internal/tasks"
	"github.com/webbben/task/internal/types"
	listcomponent "github.com/webbben/task/internal/ui/components/list-component"
	noteviewer "github.com/webbben/task/internal/ui/components/note-viewer"
)

type model struct {
	content    *types.Task
	noteViewer noteviewer.NoteViewerModel
	noteList   listcomponent.ListComponentModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// note viewer
	newNoteViewerModel, cmd := m.noteViewer.Update(msg)
	m.noteViewer = newNoteViewerModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.noteViewer.View()
}

func (m *model) setSelectedNote(noteTitle string) {
	noteContent, exists := m.content.Notes[noteTitle]
	if !exists {
		return
	}
	m.noteViewer.SetNoteContent(noteTitle, noteContent)
}

func RunUI(taskID string) error {
	task, err := tasks.GetTask(taskID)
	if err != nil {
		return fmt.Errorf("failed to run task UI: %d", err)
	}

	m := model{
		content:    task,
		noteViewer: *noteviewer.New("", ""),
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error occurred while viewing task: %d", err)
	}
	return nil
}
