package zlog

// A Config is a zlog config.
type Config struct {
	// Name represents the service name.
	Name string `meta:",optional"`
	// Mode represents the logging mode, default is `console`.
	// console: log to console.
	// file: log to file.
	// volume: used in k8s, prepend the hostname to the log file name.
	Mode string `meta:",default=console,options=console|file|volume"`
	// Level represents the log level, default is `info`.
	Level string `meta:",default=info,options=debug|trace|info|warn|error|fatal|panic"`
	// Encoding represents the encoding type, default is `json`.
	// json: json encoding.
	// plain: plain text encoding, typically used in development.
	Encoding string `meta:",default=json,options=json|plain"`
	// TimeFormat represents the time format, default is `2006-01-02T15:04:05.000Z07:00`.
	TimeFormat string `meta:",optional"`
	TimeField  string `meta:",optional"`
	// Path represents the log file path, default is `logs`.
	Path string `meta:",default=logs"`
	// Compress represents whether to compress the log file, default is `false`.
	Compress bool `meta:",optional"`
	// Stat represents whether to log statistics, default is `true`.
	Stat bool `meta:",default=true"`
	// MaxBackups represents how many backup log files will be kept. 0 means all files will be kept forever.
	// Only take effect when RotationRuleType is `size`.
	// Even though `MaxBackups` sets 0, log files will still be removed
	// if the `KeepDays` limitation is reached.
	MaxBackups int `meta:",default=0"`
	// MaxSize represents how much space the writing log file takes up. 0 means no limit. The unit is `MB`.
	// Only take effect when RotationRuleType is `size`
	MaxSize int  `meta:",default=0"`
	Caller  int  `meta:",default=0"`
	Async   bool `meta:",default=false"`
}
