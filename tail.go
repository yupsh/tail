package tail

import (
	"bufio"
	"context"
	"fmt"
	"io"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"

	localopt "github.com/yupsh/tail/opt"
)

// Flags represents the configuration options for the tail command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Tail creates a new tail command with the given parameters
func Tail(parameters ...any) yup.Command {
	cmd := command(opt.Args[string, Flags](parameters...))
	// Set default if no lines/bytes specified
	if cmd.Flags.Lines == 0 && cmd.Flags.Bytes == 0 {
		cmd.Flags.Lines = 10
	}
	return cmd
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	return yup.ProcessFilesWithContext(
		ctx, c.Positional, stdin, stdout, stderr,
		yup.FileProcessorOptions{
			CommandName:     "tail",
			ShowHeaders:     !bool(c.Flags.Quiet),
			BlankBetween:    true,
			ContinueOnError: true,
		},
		func(ctx context.Context, source yup.InputSource, output io.Writer) error {
			return c.processReader(ctx, source.Reader, output, source.Filename, false)
		},
	)
}

func (c command) processReader(ctx context.Context, reader io.Reader, output io.Writer, filename string, showHeader bool) error {
	if showHeader {
		fmt.Fprintf(output, "==> %s <==\n", filename)
	}

	if c.Flags.Bytes > 0 {
		return c.processBytes(ctx, reader, output)
	}

	return c.processLines(ctx, reader, output)
}

func (c command) processLines(ctx context.Context, reader io.Reader, output io.Writer) error {
	var lines []string
	scanner := bufio.NewScanner(reader)

	// Read all lines into memory (simple implementation)
	for yup.ScanWithContext(ctx, scanner) {
		lines = append(lines, scanner.Text())
	}

	// Check if context was cancelled
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Output the last N lines
	start := len(lines) - int(c.Flags.Lines)
	if start < 0 {
		start = 0
	}

	for i := start; i < len(lines); i++ {
		fmt.Fprintln(output, lines[i])
	}

	return nil
}

func (c command) processBytes(ctx context.Context, reader io.Reader, output io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Read all data into memory (simple implementation)
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// Output the last N bytes
	start := len(data) - int(c.Flags.Bytes)
	if start < 0 {
		start = 0
	}

	_, writeErr := output.Write(data[start:])
	return writeErr
}
