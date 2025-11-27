package app

import (
	"fmt"
	"github.com/pachirode/docker-demo/pkg/flags"
	"github.com/pachirode/docker-demo/pkg/term"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Command struct {
	usage    string
	desc     string
	options  ClipOptions
	commands []*Command
	runFunc  RunCommandFunc
}

type CommandOption func(*Command)

func WithCommandOption(opts ClipOptions) CommandOption {
	return func(cmd *Command) {
		cmd.options = opts
	}
}

type RunCommandFunc func(args []string) error

func WithRunCommandFunc(run RunCommandFunc) CommandOption {
	return func(cmd *Command) {
		cmd.runFunc = run
	}
}

func NewCommand(usage string, desc string, opts ...CommandOption) *Command {
	cmd := &Command{
		usage: usage,
		desc:  desc,
	}

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

func (cmd *Command) AddCommand(otherCmd *Command) {
	cmd.commands = append(cmd.commands, otherCmd)
}

func (cmd *Command) AddCommands(otherCmds ...*Command) {
	cmd.commands = append(cmd.commands, otherCmds...)
}

func (cmd *Command) cobraCommand() *cobra.Command {
	cobraCommand := &cobra.Command{
		Use:  cmd.usage,
		Long: cmd.desc,
	}
	cobraCommand.SetOutput(os.Stdout)
	cobraCommand.Flags().SortFlags = false

	if len(cmd.commands) > 0 {
		for _, command := range cmd.commands {
			cobraCommand.AddCommand(command.cobraCommand())
		}
	}

	if cmd.runFunc != nil {
		cobraCommand.Run = cmd.runCommand
	}

	var namedFlagSets flags.NamedFlagSets
	if cmd.options != nil {
		namedFlagSets = cmd.options.Flags()
		for _, flagSet := range namedFlagSets.FlagSetMap {
			cobraCommand.Flags().AddFlagSet(flagSet)
		}
	}
	addHelpCommandFlag(cmd.usage, cobraCommand.Flags())

	usageFmt := "Usage:\n %s\n"
	cols, _, _ := term.TerminalSize(cobraCommand.OutOrStdout())
	cobraCommand.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		flags.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
	cobraCommand.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		flags.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})
	return cobraCommand
}

func (cmd *Command) runCommand(cobraCommand *cobra.Command, args []string) {
	if cmd.runFunc != nil {
		if err := cmd.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}
