package command_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/tail"
)

// ==============================================================================
// Test Default Behavior (10 lines)
// ==============================================================================

func TestTail_DefaultTenLines(t *testing.T) {
	lines := make([]string, 15)
	for i := range lines {
		lines[i] = fmt.Sprintf("%d", i+1)
	}

	result := run.Command(command.Tail()).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"6", "7", "8", "9", "10", "11", "12", "13", "14", "15",
	})
}

func TestTail_LessThanDefault(t *testing.T) {
	result := run.Command(command.Tail()).
		WithStdinLines("1", "2", "3", "4", "5").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"1", "2", "3", "4", "5",
	})
}

func TestTail_ExactlyTenLines(t *testing.T) {
	lines := make([]string, 10)
	for i := range lines {
		lines[i] = fmt.Sprintf("%d", i+1)
	}

	result := run.Command(command.Tail()).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, lines)
}

func TestTail_EmptyInput(t *testing.T) {
	result := run.Quick(command.Tail())

	assertion.NoError(t, result.Err)
	assertion.Empty(t, result.Stdout)
}

func TestTail_SingleLine(t *testing.T) {
	result := run.Command(command.Tail()).
		WithStdinLines("only line").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"only line"})
}

// ==============================================================================
// Test Custom Line Counts
// ==============================================================================

func TestTail_ThreeLines(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines("1", "2", "3", "4", "5").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"3", "4", "5"})
}

func TestTail_OneLine(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(1))).
		WithStdinLines("first", "second", "third").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"third"})
}

func TestTail_FiveLines(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(5))).
		WithStdinLines("a", "b", "c", "d", "e", "f", "g").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"c", "d", "e", "f", "g"})
}

func TestTail_LargeCount(t *testing.T) {
	// Request 100 lines, but only provide 5
	result := run.Command(command.Tail(command.LineCount(100))).
		WithStdinLines("1", "2", "3", "4", "5").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"1", "2", "3", "4", "5"})
}

// ==============================================================================
// Test With Empty Lines
// ==============================================================================

func TestTail_EmptyLine(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines("a", "b", "", "", "c").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"", "", "c"})
}

func TestTail_AllEmptyLines(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines("", "", "", "", "").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"", "", ""})
}

func TestTail_EmptyLinesAtEnd(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines("content", "more", "", "").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"more", "", ""})
}

// ==============================================================================
// Test With Whitespace
// ==============================================================================

func TestTail_Whitespace(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines("normal", "  spaces", "\ttabs", "   both \t").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"  spaces",
		"\ttabs",
		"   both \t",
	})
}

func TestTail_WhitespaceOnly(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(2))).
		WithStdinLines("content", "   ", "\t\t").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"   ", "\t\t"})
}

func TestTail_LeadingSpaces(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(2))).
		WithStdinLines("line1", "    line2", "        line3").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"    line2",
		"        line3",
	})
}

func TestTail_TrailingSpaces(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(2))).
		WithStdinLines("line1", "line2    ", "line3        ").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"line2    ",
		"line3        ",
	})
}

// ==============================================================================
// Test With Unicode
// ==============================================================================

func TestTail_Unicode_Japanese(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(2))).
		WithStdinLines("こんにちは", "世界", "日本語").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"世界", "日本語"})
}

func TestTail_Unicode_Mixed(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines("Hello", "世界", "123", "مرحبا", "test").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"123", "مرحبا", "test"})
}

func TestTail_Unicode_Emoji(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(2))).
		WithStdinLines("😀", "👋", "🌍", "🚀").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"🌍", "🚀"})
}

func TestTail_Unicode_Arabic(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(2))).
		WithStdinLines("مرحبا", "سلام", "أهلا").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 2)
}

// ==============================================================================
// Test With Special Characters
// ==============================================================================

func TestTail_SpecialCharacters(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines("normal", "!@#$%", "^&*()", "{}[]", "<>?").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"^&*()", "{}[]", "<>?"})
}

func TestTail_Punctuation(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(2))).
		WithStdinLines("Hello!", "How are you?", "Goodbye.").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"How are you?", "Goodbye."})
}

func TestTail_Quotes(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(2))).
		WithStdinLines(`"double"`, `'single'`, "`backtick`").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{`'single'`, "`backtick`"})
}

// ==============================================================================
// Test Edge Cases
// ==============================================================================

func TestTail_VeryLongLine(t *testing.T) {
	longLine := strings.Repeat("a", 10000)
	result := run.Command(command.Tail(command.LineCount(1))).
		WithStdinLines("short", longLine).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 1)
	assertion.Equal(t, result.Stdout[0], longLine, "long line")
}

func TestTail_ManyLines(t *testing.T) {
	lines := make([]string, 1000)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i+1)
	}

	result := run.Command(command.Tail(command.LineCount(10))).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 10)
	assertion.Equal(t, result.Stdout[0], "line 991", "991st line")
	assertion.Equal(t, result.Stdout[9], "line 1000", "1000th line")
}

func TestTail_ExactLineCount(t *testing.T) {
	// Request exactly as many lines as provided
	result := run.Command(command.Tail(command.LineCount(5))).
		WithStdinLines("1", "2", "3", "4", "5").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"1", "2", "3", "4", "5"})
}

// ==============================================================================
// Test Error Handling
// ==============================================================================

func TestTail_InputError(t *testing.T) {
	result := run.Command(command.Tail()).
		WithStdinError(errors.New("read failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "read failed")
}

func TestTail_OutputError(t *testing.T) {
	result := run.Command(command.Tail()).
		WithStdinLines("test").
		WithStdoutError(errors.New("write failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "write failed")
}

// ==============================================================================
// Test Flags
// ==============================================================================

func TestTail_BytesFlag(t *testing.T) {
	// ByteCount flag is defined but not currently used in implementation
	result := run.Command(command.Tail(command.ByteCount(10))).
		WithStdinLines("a", "b", "c").
		Run()

	assertion.NoError(t, result.Err)
	// Current implementation ignores bytes flag, defaults to 10 lines
}

func TestTail_StartFromLineFlag(t *testing.T) {
	// StartFromLine flag is defined but not currently used in implementation
	result := run.Command(command.Tail(command.StartFromLine(5))).
		WithStdinLines("a", "b", "c").
		Run()

	assertion.NoError(t, result.Err)
	// Current implementation ignores this flag
}

func TestTail_FollowFlag(t *testing.T) {
	// Follow flag is defined but not currently used in implementation
	result := run.Command(command.Tail(command.Follow)).
		WithStdinLines("a", "b").
		Run()

	assertion.NoError(t, result.Err)
	// Current implementation ignores follow flag
}

func TestTail_QuietFlag(t *testing.T) {
	// Quiet flag is defined but not currently used in implementation
	result := run.Command(command.Tail(command.Quiet)).
		WithStdinLines("a", "b").
		Run()

	assertion.NoError(t, result.Err)
	// Current implementation ignores quiet flag
}

func TestTail_FollowRetryFlag(t *testing.T) {
	// FollowRetry flag is defined but not currently used in implementation
	result := run.Command(command.Tail(command.FollowRetry)).
		WithStdinLines("a", "b").
		Run()

	assertion.NoError(t, result.Err)
	// Current implementation ignores follow retry flag
}

func TestTail_VerboseFlag(t *testing.T) {
	// Verbose flag is defined but not currently used in implementation
	result := run.Command(command.Tail(command.Verbose)).
		WithStdinLines("a", "b").
		Run()

	assertion.NoError(t, result.Err)
	// Current implementation ignores verbose flag
}

func TestTail_SuppressHeadersFlag(t *testing.T) {
	// SuppressHeaders flag is defined but not currently used in implementation
	result := run.Command(command.Tail(command.SuppressHeaders)).
		WithStdinLines("a", "b").
		Run()

	assertion.NoError(t, result.Err)
	// Current implementation ignores suppress headers flag
}

func TestTail_AlwaysHeadersFlag(t *testing.T) {
	// AlwaysHeaders flag is defined but not currently used in implementation
	result := run.Command(command.Tail(command.AlwaysHeaders)).
		WithStdinLines("a", "b").
		Run()

	assertion.NoError(t, result.Err)
	// Current implementation ignores always headers flag
}

// ==============================================================================
// Table-Driven Tests
// ==============================================================================

func TestTail_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		count    command.LineCount
		input    []string
		expected []string
	}{
		{
			name:     "three from five",
			count:    3,
			input:    []string{"a", "b", "c", "d", "e"},
			expected: []string{"c", "d", "e"},
		},
		{
			name:     "one line",
			count:    1,
			input:    []string{"first", "second", "third"},
			expected: []string{"third"},
		},
		{
			name:     "all lines",
			count:    5,
			input:    []string{"a", "b"},
			expected: []string{"a", "b"},
		},
		{
			name:     "with empty lines",
			count:    3,
			input:    []string{"a", "b", "", "c"},
			expected: []string{"b", "", "c"},
		},
		{
			name:     "unicode",
			count:    2,
			input:    []string{"こんにちは", "世界", "日本語"},
			expected: []string{"世界", "日本語"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Command(command.Tail(tt.count)).
				WithStdinLines(tt.input...).
				Run()

			assertion.NoError(t, result.Err)
			assertion.Lines(t, result.Stdout, tt.expected)
		})
	}
}

// ==============================================================================
// Test Real-World Scenarios
// ==============================================================================

func TestTail_LogFileBottom(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines(
			"[2024-01-01] First entry",
			"[2024-01-02] Second entry",
			"[2024-01-03] Third entry",
			"[2024-01-04] Fourth entry",
			"[2024-01-05] Fifth entry",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"[2024-01-03] Third entry",
		"[2024-01-04] Fourth entry",
		"[2024-01-05] Fifth entry",
	})
}

func TestTail_CodeSnippet(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines(
			"package main",
			"",
			"import \"fmt\"",
			"",
			"func main() {",
			"    fmt.Println(\"Hello\")",
			"}",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"func main() {",
		"    fmt.Println(\"Hello\")",
		"}",
	})
}

func TestTail_CSVLastRows(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(2))).
		WithStdinLines(
			"Name,Age,City",
			"Alice,30,NYC",
			"Bob,25,LA",
			"Carol,35,SF",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"Bob,25,LA",
		"Carol,35,SF",
	})
}

func TestTail_DataLastRows(t *testing.T) {
	// Get last few rows of data
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines(
			"Row 1",
			"Row 2",
			"Row 3",
			"Row 4",
			"Row 5",
			"Row 6",
			"Row 7",
			"Row 8",
			"Row 9",
			"Row 10",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
	assertion.Equal(t, result.Stdout[0], "Row 8", "8th row")
	assertion.Equal(t, result.Stdout[2], "Row 10", "10th row")
}

// ==============================================================================
// Test Line Number Boundaries
// ==============================================================================

func TestTail_OneFromMany(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(1))).
		WithStdinLines("first", "second", "third", "fourth").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"fourth"})
}

func TestTail_TwoFromThree(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(2))).
		WithStdinLines("a", "b", "c").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"b", "c"})
}

func TestTail_NineFromTen(t *testing.T) {
	lines := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	expected := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10"}

	result := run.Command(command.Tail(command.LineCount(9))).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, expected)
}

// ==============================================================================
// Test Mixed Content
// ==============================================================================

func TestTail_MixedContent(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(5))).
		WithStdinLines(
			"normal text",
			"",
			"line with\ttabs",
			"  spaces  ",
			"unicode: 日本語",
			"special: !@#$",
			"more content",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 5)
	assertion.Equal(t, result.Stdout[0], "line with\ttabs", "line with tabs")
	assertion.Equal(t, result.Stdout[1], "  spaces  ", "spaces line")
	assertion.Equal(t, result.Stdout[2], "unicode: 日本語", "unicode line")
	assertion.Equal(t, result.Stdout[3], "special: !@#$", "special line")
	assertion.Equal(t, result.Stdout[4], "more content", "last line")
}

// ==============================================================================
// Test Head vs Tail Comparison
// ==============================================================================

func TestTail_Complementary(t *testing.T) {
	// tail gets the last N lines (opposite of head)
	input := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines(input...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"8", "9", "10"})
}

// ==============================================================================
// Test Empty Lines at Different Positions
// ==============================================================================

func TestTail_EmptyLinesAtStart(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(3))).
		WithStdinLines("", "", "content", "more", "last").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"content",
		"more",
		"last",
	})
}

func TestTail_EmptyLinesInMiddle(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(4))).
		WithStdinLines("before", "", "", "after").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"before",
		"",
		"",
		"after",
	})
}

// ==============================================================================
// Test Buffer Behavior
// ==============================================================================

func TestTail_BuffersAllLines(t *testing.T) {
	// tail must buffer all lines to know which are the last N
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i+1)
	}

	result := run.Command(command.Tail(command.LineCount(5))).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 5)
	// Verify it got the correct last 5 lines
	assertion.Equal(t, result.Stdout[0], "line 96", "96th line")
	assertion.Equal(t, result.Stdout[4], "line 100", "100th line")
}

// ==============================================================================
// Test Two Lines
// ==============================================================================

func TestTail_TwoLines(t *testing.T) {
	result := run.Command(command.Tail()).
		WithStdinLines("first", "second").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"first", "second"})
}

func TestTail_TwoLinesCustomCount(t *testing.T) {
	result := run.Command(command.Tail(command.LineCount(1))).
		WithStdinLines("first", "second").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"second"})
}

