package todo

import (
	"log"

	"github.com/kingstar718/gtr/internal/todo"
	"github.com/spf13/cobra"
)

func NewTodoCommand() *cobra.Command {
	storage, err := todo.NewStorage()
	if err != nil {
		log.Fatal(err)
	}

	service := todo.NewService(storage)

	cmd := &cobra.Command{
		Use:     "todo",
		Aliases: []string{"t", "-t", "--t", "-todo", "--todo"},
		Short: "\n-------------------------------------\n" +
			"| COMMAND: todo                   |\n" +
			"| TYPE: Task Management           |\n" +
			"| INPUT:                          |\n" +
			"|   1. add <title>                |\n" +
			"|   2. list [--status done]       |\n" +
			"|   3. delete <id>                |\n" +
			"|   4. (no args) - TUI Mode       |\n" +
			"| EXAMPLES:                       |\n" +
			"|   1. todo add \"Finish report\"   |\n" +
			"|   2. todo list --status done    |\n" +
			"|   3. todo delete <task-id>      |\n" +
			"|   4. todo                       |\n" +
			"-------------------------------------\n",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return RunTUI(service)
			}
			return cmd.Help()
		},
	}

	cmd.AddCommand(NewAddCommand(service))
	cmd.AddCommand(NewListCommand(service))
	cmd.AddCommand(NewDeleteCommand(service))

	return cmd
}
