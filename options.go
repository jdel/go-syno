package syno

import (
	"io"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// LogPanicLevel is a shorthand for logrus.PanicLevel
var LogPanicLevel = func() log.Level { return log.PanicLevel }()

// LogFatalLevel is a shorthand for logrus.FatalLevel
var LogFatalLevel = func() log.Level { return log.FatalLevel }()

// LogErrorLevel is a shorthand for logrus.PaniErrorLevelcLevel
var LogErrorLevel = func() log.Level { return log.ErrorLevel }()

// LogWarnLevel is a shorthand for logrus.WarnLevel
var LogWarnLevel = func() log.Level { return log.WarnLevel }()

// LogInfoLevel is a shorthand for logrus.InfoLevel
var LogInfoLevel = func() log.Level { return log.InfoLevel }()

// LogDebugLevel is a shorthand for logrus.DebugLevel
var LogDebugLevel = func() log.Level { return log.DebugLevel }()

// Options holds external context
type Options struct {
	PackagesDir string
	CacheDir    string
	ModelsFile  string
	Language    string
	MD5         bool
}

var o Options

func init() {
	// Default options
	o = Options{
		PackagesDir: filepath.Join(executablePath(), "packages"),
		CacheDir:    filepath.Join(executablePath(), "cache"),
		ModelsFile:  filepath.Join(executablePath(), "models.yml"),
		Language:    "enu",
	}
	log.SetLevel(log.ErrorLevel)
}

// SetLogLevel sets the logrus log level
func SetLogLevel(l log.Level) {
	log.SetLevel(l)
}

// SetLogOutput sets the logrus log output
func SetLogOutput(i io.Writer) {
	log.SetOutput(i)
}

// SetOptions sets global options
func SetOptions(opt Options) {
	log.Debug("Overriding default options")
	o = opt
}

// GetOptions returns global options
func GetOptions() *Options {
	return &o
}
