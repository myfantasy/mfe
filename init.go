package mfe

import log "github.com/sirupsen/logrus"

var actionLevel log.Level
var innerErrorLevel log.Level
var extErrorLevel log.Level

func init() {
	actionLevel = log.DebugLevel
	innerErrorLevel = log.WarnLevel
	extErrorLevel = log.WarnLevel
}

// ActionLevelSet set log level for action info (default DebugLevel)
func ActionLevelSet(level log.Level) {
	actionLevel = level
}

// InnerErrorLevelSet set log level for inner error (default DebugLevel)
func InnerErrorLevelSet(level log.Level) {
	innerErrorLevel = level
}

// ExtErrorLevelSet set log level for External error (default DebugLevel)
func ExtErrorLevelSet(level log.Level) {
	extErrorLevel = level
}

// LogInnerError send inner error (public for mfetl)
func LogInnerError(msg string) {
	LogWrapper(innerErrorLevel, msg, nil)
}

// LogExtError send External error info (public for mfetl)
func LogExtError(msg string) {
	LogWrapper(extErrorLevel, msg, nil)
}

// LogAction sent action info (public for mfetl)
func LogAction(msg string) {
	LogWrapper(actionLevel, msg, nil)
}

// LogInnerErrorF send inner error (public for mfetl)
func LogInnerErrorF(msg string, method string, action string) {
	LogWrapper(innerErrorLevel, msg, log.Fields{"method": method, "action": action})
}

// LogExtErrorF send External error info (public for mfetl)
func LogExtErrorF(msg string, method string, action string) {
	LogWrapper(extErrorLevel, msg, log.Fields{"method": method, "action": action})
}

// LogActionF sent action info (public for mfetl)
func LogActionF(msg string, method string, action string) {
	LogWrapper(actionLevel, msg, log.Fields{"method": method, "action": action})
}

// LogWrapper send log with custom level (public for mfetl)
func LogWrapper(level log.Level, msg string, fields log.Fields) {
	if fields == nil || len(fields) == 0 {
		if level == log.DebugLevel {
			log.Debug(msg)
			return
		}
		if level == log.InfoLevel {
			log.Info(msg)
			return
		}
		if level == log.WarnLevel {
			log.Warn(msg)
			return
		}
		if level == log.ErrorLevel {
			log.Error(msg)
			return
		}
		if level == log.FatalLevel {
			log.Fatal(msg)
			return
		}
		if level == log.PanicLevel {
			log.Panic(msg)
			return
		}
	}
	var entr = log.WithFields(fields)
	if level == log.DebugLevel {
		entr.Debug(msg)
		return
	}
	if level == log.InfoLevel {
		entr.Info(msg)
		return
	}
	if level == log.WarnLevel {
		entr.Warn(msg)
		return
	}
	if level == log.ErrorLevel {
		entr.Error(msg)
		return
	}
	if level == log.FatalLevel {
		entr.Fatal(msg)
		return
	}
	if level == log.PanicLevel {
		entr.Panic(msg)
		return
	}
}
