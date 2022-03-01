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

			log := logger.With().Logger()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				t2 := time.Now()
				requestID := r.Header.Get("X-Request-ID")

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.Error().
						Str("type", "error").
						Timestamp().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Str("request_id", requestID).
						Msg("log system error")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// log end request
				log.Info().
					Str("type", "access").
					Timestamp().
					Fields(map[string]interface{}{
						"url":        r.URL.Path,
						"method":     r.Method,
						"status":     ww.Status(),
						"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
						"bytes_out":  ww.BytesWritten(),
						"request_id": requestID,
					}).
					Msg("incoming_request")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func NewRollingFile(file string) (io.Writer, error) {

	dir, _ := filepath.Split(file)
	if err := os.MkdirAll(dir, 0744); err != nil {
		return nil, errors.New("Can't create log directory at:'" + dir + "' with error:" + err.Error())
	}

	return &lumberjack.Logger{
		Filename:   file,
		MaxBackups: 3,  // files
		MaxSize:    1,  // megabytes
		MaxAge:     10, // days
	}, nil
}

func Initialize(logfile string) {

	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	logFile, err := NewRollingFile(logfile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to allocate new rolling file")
	}

	var writers []io.Writer
	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	writers = append(writers, logFile)

	mw := io.MultiWriter(writers...)
	logger := zerolog.New(mw).With().Timestamp().Logger()
	log.Logger = logger
}
