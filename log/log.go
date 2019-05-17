package log

type Log interface {
	Info(...interface{})
	Error(...interface{})
}

var DefaultLog Log = new(stdLog)
