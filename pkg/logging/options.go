package logging

func (l *logger) Info(msg string) {
	l.logStd.Info(msg)
	l.logFile.Info(msg)
}

func (l *logger) Error(msg string) {
	l.logStd.Error(msg)
	l.logFile.Error(msg)
}
