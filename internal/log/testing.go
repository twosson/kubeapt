package log

import "testing"

type testingLogger struct {
	t *testing.T
}

// TestLogger returns a logger for tests
func TestLogger(t *testing.T) Logger {
	return &testingLogger{t: t}
}

func (t *testingLogger) Debugf(format string, args ...interface{}) {
	t.t.Logf(format, args...)
}
func (t *testingLogger) Infof(format string, args ...interface{}) {
	t.t.Logf(format, args...)
}
func (t *testingLogger) Warnf(format string, args ...interface{}) {
	t.t.Logf(format, args...)
}
func (t *testingLogger) Errorf(format string, args ...interface{}) {
	t.t.Errorf(format, args...)
}
func (t *testingLogger) With(args ...interface{}) Logger {
	return t
}
func (t *testingLogger) Named(string) Logger {
	return t
}
