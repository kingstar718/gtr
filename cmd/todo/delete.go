package todo

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kingstar718/gtr/internal/todo"
	"github.com/spf13/cobra"
)

func NewDeleteCommand(service *todo.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a task",
		Long: `Delete a task by ID.

Usage:
  gtr todo delete <task-id>
  gtr todo delete                 # Interactive mode`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return deleteTaskCLI(service, args[0])
			}
			return deleteTaskInteractive(service)
		},
	}

	return cmd
}

func deleteTaskCLI(service *todo.Service, id string) error {
	task, err := service.GetTaskByID(id)
	if err != nil {
		return err
	}

	if task == nil {
		return fmt.Errorf("task not found: %s", id)
	}

	confirm := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Delete task '%s'?", task.Title),
		Default: false,
	}
	survey.AskOne(prompt, &confirm)

	if !confirm {
		fmt.Println("❌ Delete cancelled.")
		return nil
	}

	if err := service.DeleteTask(id); err != nil {
		return err
	}

	fmt.Printf("✅ Task deleted successfully!\n\n")
	return nil
}

func deleteTaskInteractive(service *todo.Service) error {
	tasks, err := service.GetAllTasks()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks to delete.")
		return nil
	}

	taskOptions := make([]string, len(tasks))
	for i, task := range tasks {
		taskOptions[i] = fmt.Sprintf("[%s] %s (%s)", task.Priority, task.Title, task.ID[:12])
	}

	var selectedIndex int
	prompt := &survey.Select{
		Message: "Select task to delete:",
		Options: taskOptions,
	}
	if err := survey.AskOne(prompt, &selectedIndex); err != nil {
		return err
	}

	selectedTask := tasks[selectedIndex]

	confirm := false
	confirmPrompt := &survey.Confirm{
		Message: fmt.Sprintf("Delete task '%s'?", selectedTask.Title),
		Default: false,
	}
	survey.AskOne(confirmPrompt, &confirm)

	if !confirm {
		fmt.Println("❌ Delete cancelled.")
		return nil
	}

	if err := service.DeleteTask(selectedTask.ID); err != nil {
		return err
	}

	fmt.Printf("✅ Task deleted successfully!\n\n")
	return nil
}
