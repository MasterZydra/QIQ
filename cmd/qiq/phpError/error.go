package phpError

import "fmt"

const (
	E_ERROR             int64 = 1
	E_WARNING           int64 = 2
	E_PARSE             int64 = 4
	E_NOTICE            int64 = 8
	E_CORE_ERROR        int64 = 16
	E_CORE_WARNING      int64 = 32
	E_COMPILE_ERROR     int64 = 64
	E_COMPILE_WARNING   int64 = 128
	E_USER_ERROR        int64 = 256
	E_USER_WARNING      int64 = 512
	E_USER_NOTICE       int64 = 1024
	E_STRICT            int64 = 2048
	E_RECOVERABLE_ERROR int64 = 4096
	E_DEPRECATED        int64 = 8192
	E_USER_DEPRECATED   int64 = 16384
	E_ALL               int64 = 32767
)

type ErrorType string

const (
	ErrorPhpError          ErrorType = "Error"
	WarningPhpError        ErrorType = "Warning"
	ParsePhpError          ErrorType = "ParserError"
	NoticePhpError         ErrorType = "Notice"
	CorePhpError           ErrorType = "CoreError"
	CoreWarningPhpError    ErrorType = "CoreWarning"
	CompilePhpError        ErrorType = "CompilerError"
	CompileWarningPhpError ErrorType = "CompilerWarning"
	UserPhpError           ErrorType = "UserError"
	UserWarningPhpError    ErrorType = "UserWarning"
	UserNoticePhpError     ErrorType = "UserNotice"
	StrictPhpError         ErrorType = "Strict"
	RecoverablePhpError    ErrorType = "RecoverableError"
	DeprecatedPhpError     ErrorType = "Deprecated"
	UserDeprecatedPhpError ErrorType = "UserDeprecated"
	// Non-PHP error types
	EventError ErrorType = "Event"
)

const (
	ExitEvent     string = "exit"
	ReturnEvent   string = "return"
	ContinueEvent string = "continue"
	BreakEvent    string = "break"
)

type Error interface {
	GetErrorType() ErrorType
	GetMessage() string
	GetRawMessage() string
	Error() string
}

// MARK: PhpError

type PhpError struct {
	errorType ErrorType
	message   string
}

func (err *PhpError) GetErrorType() ErrorType { return err.errorType }

func (err *PhpError) GetMessage() string { return err.String() }

func (err *PhpError) GetRawMessage() string { return err.message }

func (err *PhpError) Error() string { return err.String() }

func (err *PhpError) String() string {
	switch err.errorType {
	case NoticePhpError:
		return "Notice: " + err.message
	case WarningPhpError:
		return "Warning: " + err.message
	case ErrorPhpError:
		return "Fatal error: " + err.message
	case ParsePhpError:
		return "Parse error: " + err.message
	case DeprecatedPhpError:
		return "Deprecated: " + err.message
	default:
		return err.message
	}
}

func NewParseError(format string, a ...any) Error {
	return &PhpError{errorType: ParsePhpError, message: fmt.Sprintf(format, a...)}
}

func NewError(format string, a ...any) Error {
	return &PhpError{errorType: ErrorPhpError, message: fmt.Sprintf(format, a...)}
}

func NewWarning(format string, a ...any) Error {
	return &PhpError{errorType: WarningPhpError, message: fmt.Sprintf(format, a...)}
}

func NewNotice(format string, a ...any) Error {
	return &PhpError{errorType: NoticePhpError, message: fmt.Sprintf(format, a...)}
}

func NewEvent(event string) Error {
	return &PhpError{errorType: EventError, message: event}
}

func NewDeprecatedError(format string, a ...any) Error {
	return &PhpError{errorType: DeprecatedPhpError, message: fmt.Sprintf(format, a...)}
}

// MARK: ContinueEventError

type ContinueEventError struct {
	*PhpError
	breakoutLevel int64
}

func (err *ContinueEventError) GetBreakoutLevel() int64 { return err.breakoutLevel }

func NewBreakEvent(breakoutLevel int64) Error {
	return &ContinueEventError{PhpError: &PhpError{errorType: EventError, message: BreakEvent}, breakoutLevel: breakoutLevel}
}

func NewContinueEvent(breakoutLevel int64) Error {
	return &ContinueEventError{PhpError: &PhpError{errorType: EventError, message: ContinueEvent}, breakoutLevel: breakoutLevel}
}
