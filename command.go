package command

import (
	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[yup.File, flags]

func Tail(parameters ...any) yup.Command {
	cmd := command(yup.Initialize[yup.File, flags](parameters...))
	if cmd.Flags.Lines == 0 && cmd.Flags.Bytes == 0 {
		cmd.Flags.Lines = 10
	}
	return cmd
}

func (p command) Executor() yup.CommandExecutor {
	lineCount := int(p.Flags.Lines)
	if lineCount == 0 {
		lineCount = 10
	}

	return yup.Inputs[yup.File, flags](p).Wrap(
		yup.AccumulateAndProcess(func(lines []string) []string {
			// Return last N lines
			if len(lines) <= lineCount {
				return lines
			}
			return lines[len(lines)-lineCount:]
		}).Executor(),
	)
}
