package command

import (
	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[gloo.File, flags]

func Tail(parameters ...any) gloo.Command {
	cmd := command(gloo.Initialize[gloo.File, flags](parameters...))
	if cmd.Flags.Lines == 0 && cmd.Flags.Bytes == 0 {
		cmd.Flags.Lines = 10
	}
	return cmd
}

func (p command) Executor() gloo.CommandExecutor {
	lineCount := int(p.Flags.Lines)
	if lineCount == 0 {
		lineCount = 10
	}

	return gloo.Inputs[gloo.File, flags](p).Wrap(
		gloo.AccumulateAndProcess(func(lines []string) []string {
			// Return last N lines
			if len(lines) <= lineCount {
				return lines
			}
			return lines[len(lines)-lineCount:]
		}).Executor(),
	)
}
