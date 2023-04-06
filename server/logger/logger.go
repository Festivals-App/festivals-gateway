package logger

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Middleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			requestStart := time.Now()
			requestID := r.Header.Get("X-Request-ID")
			log := logger.With().Timestamp().Str("type", "access").Str("request_id", requestID).Logger()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {

				requestEnd := time.Now()
				status := ww.Status()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.Error().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("Recovered from panicking routine")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				if status < 300 {
					// log successfull requests at trace lvl
					log.Trace().
						Fields(map[string]interface{}{
							"url":        r.Host + r.URL.Path,
							"method":     r.Method,
							"status":     status,
							"latency_ms": float64(requestEnd.Sub(requestStart).Nanoseconds()) / 1000000.0,
							"bytes_out":  ww.BytesWritten(),
						}).
						Msg("Incoming request")
				} else {
					// log failed requests at debug lvl
					log.Debug().
						Fields(map[string]interface{}{
							"url":        r.Host + r.URL.Path,
							"method":     r.Method,
							"status":     status,
							"latency_ms": float64(requestEnd.Sub(requestStart).Nanoseconds()) / 1000000.0,
							"bytes_out":  ww.BytesWritten(),
						}).
						Msg("Incoming request")
				}
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func Initialize(logfile string, console bool) {

	logFile, err := NewRollingFile(logfile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to allocate new rolling log file")
	}

	var writers []io.Writer
	if console {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	writers = append(writers, logFile)

	multiWriter := io.MultiWriter(writers...)
	logger := zerolog.New(multiWriter).With().Timestamp().Logger()
	log.Logger = logger
}

func NewRollingFile(file string) (io.Writer, error) {

	dir, _ := filepath.Split(file)
	if err := os.MkdirAll(dir, 0744); err != nil {
		return nil, errors.New("Can't create log directory at:'" + dir + "' with error:" + err.Error())
	}

	return &lumberjack.Logger{
		Filename:   file,
		MaxBackups: 10, // files
		MaxSize:    50, // megabytes
		MaxAge:     31, // days
	}, nil
}
