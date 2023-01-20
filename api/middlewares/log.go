package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// StructuredLogger is a simple, but powerful implementation of a custom structured
// logger backed on logrus. I encourage users to copy it, adapt it and make it their
// own. Also take a look at https://github.com/pressly/lg for a dedicated pkg based
// on this work, designed for context-based http routers.

// NewStructuredLogger ...
// func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
// 	return middleware.RequestLogger(&StructuredLogger{logger})
// }

// StructuredLogger ...
type StructuredLogger struct {
	Logger *logrus.Logger
}

// NewLogEntry ...
// func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
// 	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
// 	logFields := logrus.Fields{}

// 	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

// 	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
// 		logFields["req_id"] = reqID
// 	}

// 	scheme := "http"
// 	if r.TLS != nil {
// 		scheme = "https"
// 	}
// 	logFields["http_scheme"] = scheme
// 	logFields["http_proto"] = r.Proto
// 	logFields["http_method"] = r.Method

// 	logFields["remote_addr"] = r.RemoteAddr
// 	logFields["user_agent"] = r.UserAgent()

// 	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

// 	entry.Logger = entry.Logger.WithFields(logFields)

// 	entry.Logger.Infoln("request started")

// 	return entry
// }

// StructuredLoggerEntry ...
type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Infoln("request complete")
}

// Panic ...
func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}
