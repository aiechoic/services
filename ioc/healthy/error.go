package healthy

import "fmt"

type Level string

const (
	LDebug Level = "debug"
	LInfo  Level = "info"
	LWarn  Level = "warn"
	LError Level = "error"
	LFatal Level = "fatal"
)

var levelMap = map[Level]int{
	LDebug: 1,
	LInfo:  2,
	LWarn:  3,
	LError: 4,
	LFatal: 5,
}

func (e Level) Number() int {
	if v, ok := levelMap[e]; ok {
		return v
	}
	return 5
}

type Error struct {
	Level Level
	Msg   string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Level, e.Msg)
}
