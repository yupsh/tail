package opt

// Custom types for parameters
type LineCount int
type ByteCount int
type StartFromLine int // For +N syntax (start from line N instead of last N)

// Boolean flag types with constants
type FollowFlag bool
const (
	Follow   FollowFlag = true
	NoFollow FollowFlag = false
)

type FollowRetryFlag bool
const (
	FollowRetry   FollowRetryFlag = true
	NoFollowRetry FollowRetryFlag = false
)

type QuietFlag bool
const (
	Quiet   QuietFlag = true
	NoQuiet QuietFlag = false
)

type VerboseFlag bool
const (
	Verbose   VerboseFlag = true
	NoVerbose VerboseFlag = false
)

type SuppressHeadersFlag bool
const (
	SuppressHeaders   SuppressHeadersFlag = true
	NoSuppressHeaders SuppressHeadersFlag = false
)

type AlwaysHeadersFlag bool
const (
	AlwaysHeaders   AlwaysHeadersFlag = true
	NoAlwaysHeaders AlwaysHeadersFlag = false
)

// Flags represents the configuration options for the tail command
type Flags struct {
	Lines           LineCount           // Number of lines to show (-n)
	Bytes           ByteCount           // Number of bytes to show (-c)
	StartFromLine   StartFromLine       // Start from line N instead of last N (+N)
	Follow          FollowFlag          // Follow file changes (-f)
	FollowRetry     FollowRetryFlag     // Follow with retry (-F)
	Quiet           QuietFlag           // Suppress headers when multiple files (-q)
	Verbose         VerboseFlag         // Always show headers (-v)
	SuppressHeaders SuppressHeadersFlag // Suppress headers (same as quiet)
	AlwaysHeaders   AlwaysHeadersFlag   // Always show headers (same as verbose)
}

// Configure methods for the opt system
func (l LineCount) Configure(flags *Flags)              { flags.Lines = l }
func (b ByteCount) Configure(flags *Flags)              { flags.Bytes = b }
func (s StartFromLine) Configure(flags *Flags)          { flags.StartFromLine = s }
func (f FollowFlag) Configure(flags *Flags)             { flags.Follow = f }
func (f FollowRetryFlag) Configure(flags *Flags)        { flags.FollowRetry = f }
func (q QuietFlag) Configure(flags *Flags)              { flags.Quiet = q }
func (v VerboseFlag) Configure(flags *Flags)            { flags.Verbose = v }
func (s SuppressHeadersFlag) Configure(flags *Flags)    { flags.SuppressHeaders = s }
func (a AlwaysHeadersFlag) Configure(flags *Flags)      { flags.AlwaysHeaders = a }
