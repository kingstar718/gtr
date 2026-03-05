package todo

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kingstar718/gtr/internal/todo"
)

type view string

const (
	tableView  view = "table"
	editView   view = "edit"
	deleteView view = "delete"
)

type Model struct {
	service        *todo.Service
	tasks          []todo.Task
	table          table.Model
	ready          bool
	quitting       bool
	filterStatus   todo.Status
	filterPriority todo.Priority
	selectedIndex  int
	currentView    view
	selectedTask   *todo.Task
	editField      int
	titleInput     textinput.Model
	descInput      textinput.Model
	duedateInput   textinput.Model
}

func NewModel(service *todo.Service) (*Model, error) {
	tasks, err := service.GetAllTasks()
	if err != nil {
		return nil, err
	}

	columns := []table.Column{
		{Title: "TITLE", Width: 35},
		{Title: "PRIORITY", Width: 10},
		{Title: "STATUS", Width: 12},
		{Title: "DUE DATE", Width: 12},
	}

	rows := tasksToRows(tasks)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("4"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("4")).
		Bold(true)
	t.SetStyles(s)

	// 初始化输入框
	titleInput := textinput.New()
	titleInput.Placeholder = "Task title"
	titleInput.CharLimit = 100
	titleInput.Width = 50

	descInput := textinput.New()
	descInput.Placeholder = "Description (optional)"
	descInput.CharLimit = 200
	descInput.Width = 50

	duedateInput := textinput.New()
	duedateInput.Placeholder = "Due date (YYYY-MM-DD)"
	duedateInput.CharLimit = 10
	duedateInput.Width = 15

	return &Model{
		service:       service,
		tasks:         tasks,
		table:         t,
		selectedIndex: 0,
		currentView:   tableView,
		titleInput:    titleInput,
		descInput:     descInput,
		duedateInput:  duedateInput,
	}, nil
}

func tasksToRows(tasks []todo.Task) []table.Row {
	rows := make([]table.Row, len(tasks))
	for i, task := range tasks {
		title := task.Title
		if len(title) > 33 {
			title = title[:30] + "..."
		}

		rows[i] = table.Row{
			title,
			string(task.Priority),
			string(task.Status),
			task.DueDate,
		}
	}
	return rows
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.table.SetHeight(msg.Height - 8)
			m.ready = true
		} else {
			m.table.SetHeight(msg.Height - 8)
		}

	case tea.KeyMsg:
		switch m.currentView {
		case tableView:
			return m.handleTableView(msg)
		case editView:
			return m.handleEditView(msg)
		case deleteView:
			return m.handleDeleteView(msg)
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	m.selectedIndex = m.table.Cursor()
	return m, cmd
}

func (m *Model) handleTableView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit

	case "up", "k":
		if m.table.Cursor() > 0 {
			m.table.SetCursor(m.table.Cursor() - 1)
			m.selectedIndex = m.table.Cursor()
		}

	case "down", "j":
		if m.table.Cursor() < len(m.tasks)-1 {
			m.table.SetCursor(m.table.Cursor() + 1)
			m.selectedIndex = m.table.Cursor()
		}

	case "e":
		if len(m.tasks) > 0 && m.selectedIndex < len(m.tasks) {
			m.currentView = editView
			m.selectedTask = &m.tasks[m.selectedIndex]
			m.editField = 0
			m.initEditInputs()
		}

	case "d":
		if len(m.tasks) > 0 && m.selectedIndex < len(m.tasks) {
			m.currentView = deleteView
			m.selectedTask = &m.tasks[m.selectedIndex]
		}

	case "1":
		m.filterStatus = todo.StatusPending
		m.refreshTasks()

	case "2":
		m.filterStatus = todo.StatusInProgress
		m.refreshTasks()

	case "3":
		m.filterStatus = todo.StatusDone
		m.refreshTasks()

	case "0":
		m.filterStatus = ""
		m.filterPriority = ""
		m.refreshTasks()

	case "h":
		m.filterPriority = todo.PriorityHigh
		m.refreshTasks()

	case "m":
		m.filterPriority = todo.PriorityMedium
		m.refreshTasks()

	case "n":
		m.filterPriority = todo.PriorityLow
		m.refreshTasks()

	case "c":
		if len(m.tasks) > 0 && m.selectedIndex < len(m.tasks) {
			m.cycleStatus(m.selectedIndex)
		}

	case "p":
		if len(m.tasks) > 0 && m.selectedIndex < len(m.tasks) {
			m.cyclePriority(m.selectedIndex)
		}

	case "?":
		m.showHelp()
	}

	return m, nil
}

func (m *Model) handleEditView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		oldField := m.editField
		m.editField = (m.editField + 1) % 4
		m.updateFieldFocus(oldField, m.editField)
		return m, nil

	case "shift+tab":
		oldField := m.editField
		m.editField = (m.editField - 1 + 4) % 4
		m.updateFieldFocus(oldField, m.editField)
		return m, nil

	case "enter":
		m.saveEdit()
		m.currentView = tableView
		return m, nil

	case "esc":
		m.currentView = tableView
		return m, nil
	}

	var cmd tea.Cmd
	switch m.editField {
	case 0:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case 1:
		m.descInput, cmd = m.descInput.Update(msg)
	case 2:
		m.duedateInput, cmd = m.duedateInput.Update(msg)
	case 3:
		switch msg.String() {
		case "up", "k":
			m.cyclePriorityInEdit()
		case "down", "j":
			m.cyclePriorityInEdit()
		}
	}

	return m, cmd
}

func (m *Model) updateFieldFocus(oldField, newField int) {
	switch oldField {
	case 0:
		m.titleInput.Blur()
	case 1:
		m.descInput.Blur()
	case 2:
		m.duedateInput.Blur()
	}

	switch newField {
	case 0:
		m.titleInput.Focus()
	case 1:
		m.descInput.Focus()
	case 2:
		m.duedateInput.Focus()
	}
}

func (m *Model) handleDeleteView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		if m.selectedTask != nil {
			m.service.DeleteTask(m.selectedTask.ID)
			m.refreshTasks()
			m.currentView = tableView
		}

	case "n", "esc":
		m.currentView = tableView
	}

	return m, nil
}

func (m *Model) initEditInputs() {
	if m.selectedTask == nil {
		return
	}

	m.titleInput.SetValue(m.selectedTask.Title)
	m.titleInput.Focus()
	m.titleInput.Blur()

	m.descInput.SetValue(m.selectedTask.Description)
	m.descInput.Focus()
	m.descInput.Blur()

	m.duedateInput.SetValue(m.selectedTask.DueDate)
	m.duedateInput.Focus()
	m.duedateInput.Blur()

	if m.editField == 0 {
		m.titleInput.Focus()
	}
}

func (m *Model) saveEdit() {
	if m.selectedTask == nil {
		return
	}

	m.selectedTask.Title = m.titleInput.Value()
	m.selectedTask.Description = m.descInput.Value()
	m.selectedTask.DueDate = m.duedateInput.Value()

	m.service.UpdateTask(m.selectedTask.ID, m.selectedTask)
	m.refreshTasks()
}

func (m *Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	if m.quitting {
		return ""
	}

	switch m.currentView {
	case tableView:
		return m.viewTable()
	case editView:
		return m.viewEdit()
	case deleteView:
		return m.viewDelete()
	}

	return ""
}

func (m *Model) viewTable() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("4")).
		MarginBottom(1)
	b.WriteString(titleStyle.Render("📋 My Tasks"))
	b.WriteString("\n")

	b.WriteString(m.table.View())
	b.WriteString("\n")

	taskCount := len(m.tasks)
	if m.filterStatus != "" || m.filterPriority != "" {
		b.WriteString(fmt.Sprintf("Filtered: status=%s priority=%s | ", m.filterStatus, m.filterPriority))
	}

	b.WriteString(fmt.Sprintf("Total: %d tasks\n", taskCount))
	b.WriteString("\n")
	b.WriteString("Commands: (e)dit | (d)elete | (c)ycle-status | (p)riority | (k/j)↑↓ | ")
	b.WriteString("Filter: (0)all (1)pending (2)inprogress (3)done (h)igh (m)edium (n)low | (?)help | (q)uit\n")

	return b.String()
}

func (m *Model) viewEdit() string {
	if m.selectedTask == nil {
		return ""
	}

	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("4"))
	b.WriteString(titleStyle.Render("✏️  Edit Task"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 50))
	b.WriteString("\n\n")

	focusStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("4")).
		Foreground(lipgloss.Color("15")).
		Bold(true)
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Width(12).
		Align(lipgloss.Right)

	titleLabel := "Title"
	if m.editField == 0 {
		titleLabel = focusStyle.Render("▶ " + titleLabel)
	} else {
		titleLabel = labelStyle.Render(titleLabel)
	}
	b.WriteString(titleLabel + " ")
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\n")

	descLabel := "Description"
	if m.editField == 1 {
		descLabel = focusStyle.Render("▶ " + descLabel)
	} else {
		descLabel = labelStyle.Render(descLabel)
	}
	b.WriteString(descLabel + " ")
	b.WriteString(m.descInput.View())
	b.WriteString("\n\n")

	duedateLabel := "Due Date"
	if m.editField == 2 {
		duedateLabel = focusStyle.Render("▶ " + duedateLabel)
	} else {
		duedateLabel = labelStyle.Render(duedateLabel)
	}
	b.WriteString(duedateLabel + " ")
	b.WriteString(m.duedateInput.View())
	b.WriteString("\n\n")

	priorityLabel := "Priority"
	priorityValue := fmt.Sprintf("[%s]", m.selectedTask.Priority)
	if m.editField == 3 {
		priorityLabel = focusStyle.Render("▶ " + priorityLabel)
		priorityValue = focusStyle.Render(priorityValue)
	} else {
		priorityLabel = labelStyle.Render(priorityLabel)
		priorityValue = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render(priorityValue)
	}
	b.WriteString(priorityLabel + " ")
	b.WriteString(priorityValue)
	b.WriteString("\n\n")

	b.WriteString(strings.Repeat("─", 50))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(
		"[TAB/S-TAB] Navigate | [↑/↓] Priority | [Enter] Save | [Esc] Cancel"))
	b.WriteString("\n")

	return b.String()
}

func (m *Model) viewDelete() string {
	if m.selectedTask == nil {
		return ""
	}

	var b strings.Builder

	warningStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("1")).
		MarginBottom(1)
	b.WriteString(warningStyle.Render("⚠️  Delete Task\n\n"))

	b.WriteString(fmt.Sprintf("Are you sure you want to delete:\n\n"))
	b.WriteString(fmt.Sprintf("  %s\n\n", m.selectedTask.Title))

	b.WriteString("Press (y) to confirm or (n) to cancel\n")

	return b.String()
}

func (m *Model) refreshTasks() {
	tasks, _ := m.service.GetAllTasks()
	if m.filterStatus != "" || m.filterPriority != "" {
		m.tasks, _ = m.service.FilterTasks(m.filterStatus, m.filterPriority, "")
	} else {
		m.tasks = tasks
	}

	m.table.SetRows(tasksToRows(m.tasks))
	if m.selectedIndex >= len(m.tasks) && m.selectedIndex > 0 {
		m.selectedIndex = len(m.tasks) - 1
		m.table.SetCursor(m.selectedIndex)
	}
}

func (m *Model) cycleStatus(idx int) {
	if idx >= 0 && idx < len(m.tasks) {
		switch m.tasks[idx].Status {
		case todo.StatusPending:
			m.tasks[idx].Status = todo.StatusInProgress
		case todo.StatusInProgress:
			m.tasks[idx].Status = todo.StatusDone
		case todo.StatusDone:
			m.tasks[idx].Status = todo.StatusPending
		}
		m.service.UpdateTask(m.tasks[idx].ID, &m.tasks[idx])
		m.refreshTasks()
	}
}

func (m *Model) cyclePriority(idx int) {
	if idx >= 0 && idx < len(m.tasks) {
		switch m.tasks[idx].Priority {
		case todo.PriorityHigh:
			m.tasks[idx].Priority = todo.PriorityMedium
		case todo.PriorityMedium:
			m.tasks[idx].Priority = todo.PriorityLow
		case todo.PriorityLow:
			m.tasks[idx].Priority = todo.PriorityHigh
		}
		m.service.UpdateTask(m.tasks[idx].ID, &m.tasks[idx])
		m.refreshTasks()
	}
}

func (m *Model) cyclePriorityInEdit() {
	if m.selectedTask == nil {
		return
	}

	switch m.selectedTask.Priority {
	case todo.PriorityHigh:
		m.selectedTask.Priority = todo.PriorityMedium
	case todo.PriorityMedium:
		m.selectedTask.Priority = todo.PriorityLow
	case todo.PriorityLow:
		m.selectedTask.Priority = todo.PriorityHigh
	}
}

func (m *Model) showHelp() {
	help := `
╔════════════════════════════════════════════════════════════╗
║           TODO Application Help                            ║
╚════════════════════════════════════════════════════════════╝

Navigation:
  ↑/↓ or k/j  - Move cursor up/down in table
  Home/End    - Jump to first/last

Actions:
  e           - Edit selected task
  d           - Delete selected task (with confirmation)
  c           - Cycle status (pending→inprogress→done→pending)
  p           - Cycle priority (high→medium→low→high)

Edit Mode (after pressing 'e'):
  tab         - Move to next field
  shift+tab   - Move to previous field
  ↑/↓         - Cycle priority in Priority field
  enter       - Save changes
  esc         - Cancel editing

Filters:
  0           - Clear all filters
  1           - Show pending tasks
  2           - Show in-progress tasks
  3           - Show done tasks
  h           - Show high priority
  m           - Show medium priority
  n           - Show low priority

General:
  q           - Quit application
  ctrl+c      - Quit application
  ?           - Show this help
`
	fmt.Println(help)
}

func RunTUI(service *todo.Service) error {
	model, err := NewModel(service)
	if err != nil {
		return err
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
