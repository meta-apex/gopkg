package zlog

type GeneralLogger interface {
	Trace(msg string)
	Tracef(format string, args ...any)
	Debug(msg string)
	Debugf(format string, args ...any)
	Info(msg string)
	Infof(format string, args ...any)
	Warn(msg string)
	Warnf(format string, args ...any)
	Error(msg string)
	Errorf(format string, args ...any)
}

type defaultGeneralLogger struct {
	log *Logger
}

func NewGeneralLogger(log *Logger, name string) GeneralLogger {
	return &defaultGeneralLogger{
		log: log.WithName(name),
	}
}

func (d *defaultGeneralLogger) Trace(msg string) {
	d.log.Trace().Msg(msg)
}

func (d *defaultGeneralLogger) Tracef(format string, args ...any) {
	d.log.Trace().Msgf(format, args...)
}

func (d *defaultGeneralLogger) Debug(msg string) {
	d.log.Debug().Msg(msg)
}

func (d *defaultGeneralLogger) Debugf(format string, args ...any) {
	d.log.Debug().Msgf(format, args...)
}

func (d *defaultGeneralLogger) Info(msg string) {
	d.log.Info().Msg(msg)
}

func (d *defaultGeneralLogger) Infof(format string, args ...any) {
	d.log.Info().Msgf(format, args...)
}

func (d *defaultGeneralLogger) Warn(msg string) {
	d.log.Warn().Msg(msg)
}

func (d *defaultGeneralLogger) Warnf(format string, args ...any) {
	d.log.Warn().Msgf(format, args...)
}

func (d *defaultGeneralLogger) Error(msg string) {
	d.log.Error().Msg(msg)
}

func (d *defaultGeneralLogger) Errorf(format string, args ...any) {
	d.log.Error().Msgf(format, args...)
}
