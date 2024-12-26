package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/log"
	"strings"
)

type saramaLogger struct{}

func (l *saramaLogger) Print(v ...interface{}) {
	log.
		Info().
		Str("logger", "sarama").
		Arr("v", v...).
		Msg("[Sarama] Print")
}

func (l *saramaLogger) Printf(format string, v ...interface{}) {
	message := strings.Builder{}
	_, _ = fmt.Fprintf(&message, format, v...)

	log.
		Info().
		Str("logger", "sarama").
		Str("message", message.String()).
		Msg("[Sarama] Printf")
}

func (l *saramaLogger) Println(v ...interface{}) {
	log.
		Info().
		Str("logger", "sarama").
		Arr("v", v...).
		Msg("[Sarama] Print")
}

// BuildLogger create custom logger for "sarama" library for debugging
func BuildLogger() sarama.StdLogger {
	return &saramaLogger{}
}
