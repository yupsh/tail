package command

type LineCount int
type ByteCount int
type StartFromLine int

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

type flags struct {
	Lines           LineCount
	Bytes           ByteCount
	StartFromLine   StartFromLine
	Follow          FollowFlag
	FollowRetry     FollowRetryFlag
	Quiet           QuietFlag
	Verbose         VerboseFlag
	SuppressHeaders SuppressHeadersFlag
	AlwaysHeaders   AlwaysHeadersFlag
}

func (l LineCount) Configure(flags *flags)           { flags.Lines = l }
func (b ByteCount) Configure(flags *flags)           { flags.Bytes = b }
func (s StartFromLine) Configure(flags *flags)       { flags.StartFromLine = s }
func (f FollowFlag) Configure(flags *flags)          { flags.Follow = f }
func (f FollowRetryFlag) Configure(flags *flags)     { flags.FollowRetry = f }
func (q QuietFlag) Configure(flags *flags)           { flags.Quiet = q }
func (v VerboseFlag) Configure(flags *flags)         { flags.Verbose = v }
func (s SuppressHeadersFlag) Configure(flags *flags) { flags.SuppressHeaders = s }
func (a AlwaysHeadersFlag) Configure(flags *flags)   { flags.AlwaysHeaders = a }
