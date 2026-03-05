package todo

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kingstar718/gtr/internal/todo"
	"github.com/spf13/cobra"
)

func NewAddCommand(service *todo.Service) *cobra.Command {
	var priority string
	var dueDate string
	var tagsStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Long: `Add a new task to your todo list.
		
Usage:
  gtr todo add "Task title" --priority high --due 2025-03-10 --tags work,urgent
  gtr todo add                    # Interactive mode`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				title := strings.Join(args, " ")
				return addTaskCLI(service, title, priority, dueDate, tagsStr)
			}

			return addTaskInteractive(service)
		},
	}

	cmd.Flags().StringVarP(&priority, "priority", "p", "medium", "Task priority: high, medium, low")
	cmd.Flags().StringVarP(&dueDate, "due", "d", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&tagsStr, "tags", "t", "", "Comma-separated tags")

	return cmd
}

func addTaskCLI(service *todo.Service, title, priority, dueDate, tagsStr string) error {
	if title == "" {
		return fmt.Errorf("task title is required")
	}

	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	pri := todo.Priority(strings.ToLower(priority))
	if pri != todo.PriorityHigh && pri != todo.PriorityMedium && pri != todo.PriorityLow {
		return fmt.Errorf("invalid priority: %s", priority)
	}

	task := todo.NewTask(title, "", pri, dueDate, tags)
	if err := service.AddTask(task); err != nil {
		return err
	}

	fmt.Printf("\n✅ Task added successfully!\n")
	fmt.Printf("ID: %s\n", task.ID)
	fmt.Printf("Title: %s\n", task.Title)
	fmt.Printf("Priority: %s\n", task.Priority)
	if dueDate != "" {
		fmt.Printf("Due Date: %s\n", dueDate)
	}
	fmt.Println()

	return nil
}

func addTaskInteractive(service *todo.Service) error {
	var title string
	var description string
	var priority string
	var dueDate string
	var tags []string

	promptTitle := &survey.Input{
		Message: "Task title:",
	}
	if err := survey.AskOne(promptTitle, &title); err != nil {
		return err
	}

	if title == "" {
		return fmt.Errorf("task title is required")
	}

	promptDesc := &survey.Input{
		Message: "Description (optional):",
	}
	survey.AskOne(promptDesc, &description)

	priorityOptions := []string{"high", "medium", "low"}
	promptPriority := &survey.Select{
		Message: "Priority:",
		Options: priorityOptions,
		Default: "medium",
	}
	if err := survey.AskOne(promptPriority, &priority); err != nil {
		return err
	}

	promptDue := &survey.Input{
		Message: "Due date (YYYY-MM-DD, optional):",
	}
	survey.AskOne(promptDue, &dueDate)

	promptTags := &survey.Input{
		Message: "Tags (comma-separated, optional):",
	}
	var tagsStr string
	survey.AskOne(promptTags, &tagsStr)
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	task := todo.NewTask(title, description, todo.Priority(priority), dueDate, tags)
	if err := service.AddTask(task); err != nil {
		return err
	}

	fmt.Printf("\n✅ Task added successfully!\n")
	fmt.Printf("ID: %s\n", task.ID)
	fmt.Println()

	return nil
}
