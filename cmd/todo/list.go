package todo

import (
	"fmt"
	"strings"

	"github.com/kingstar718/gtr/internal/todo"
	"github.com/spf13/cobra"
)

func NewListCommand(service *todo.Service) *cobra.Command {
	var status string
	var priority string
	var tag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Long: `Display all tasks with optional filters.

Usage:
  gtr todo list
  gtr todo list --status done
  gtr todo list --priority high
  gtr todo list --tag work`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var filterStatus todo.Status
			var filterPriority todo.Priority

			if status != "" {
				filterStatus = todo.Status(status)
			}
			if priority != "" {
				filterPriority = todo.Priority(priority)
			}

			tasks, err := service.FilterTasks(filterStatus, filterPriority, tag)
			if err != nil {
				return err
			}

			if len(tasks) == 0 {
				fmt.Println("No tasks found.")
				return nil
			}

			displayTaskList(tasks)
			return nil
		},
	}

	cmd.Flags().StringVarP(&status, "status", "s", "", "Filter by status: pending, inprogress, done")
	cmd.Flags().StringVarP(&priority, "priority", "p", "", "Filter by priority: high, medium, low")
	cmd.Flags().StringVarP(&tag, "tag", "t", "", "Filter by tag")

	return cmd
}

func displayTaskList(tasks []todo.Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Printf("%-16s  %-30s  %-10s  %-12s  %-12s\n", "ID", "TITLE", "PRIORITY", "STATUS", "DUE DATE")
	fmt.Println(strings.Repeat("-", 100))

	for _, task := range tasks {
		fmt.Printf("%-16s  %-30s  %-10s  %-12s  %-12s\n",
			task.ID[:16],
			truncateStr(task.Title, 28),
			task.Priority,
			task.Status,
			task.DueDate,
		)
	}

	fmt.Println()
	fmt.Printf("Total: %d tasks\n\n", len(tasks))
}

func truncateStr(s string, width int) string {
	if len(s) > width {
		return s[:width-3] + "..."
	}
	return s
}
