package util

import (
	"runtime/debug"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	typeSystem string = "system"

	service string = "gateway"

	statusSuccess    string = "success"
	statusInProgress string = "in-progress"
	statusFailed     string = "failed"
)

func PrintRequestLog(method string, requestID string, data interface{}) {
	log.Info().
		Str("type", typeSystem).
		Str("service", service).
		Str("method", method).
		Str("request_id", requestID).
		Str("status", statusInProgress).
		Dict("payload", zerolog.Dict().
			Interface("request", data),
		).
		Msg("")
}

func PrintSuccessResponseLog(method string, requestID string, data interface{}) {
	log.Info().
		Str("type", typeSystem).
		Str("service", service).
		Str("method", method).
		Str("request_id", requestID).
		Str("status", statusSuccess).
		Dict("payload", zerolog.Dict().
			Interface("response", data),
		).
		Msg("")
}

func PrintErrorResponseLog(method string, requestID string, data interface{}, errorMessage string) {
	stack := string(debug.Stack())
	stack = strings.ReplaceAll(stack, "\n\t", "\n")
	stacks := strings.Split(stack, "\n")

	log.Error().
		Str("type", typeSystem).
		Str("service", service).
		Str("method", method).
		Str("request_id", requestID).
		Str("status", statusFailed).
		Dict("payload", zerolog.Dict().
			Interface("response", data).
			Str("message", errorMessage).
			Interface("stacktrace", stacks),
		).
		Msg("")
}

func PrintErrorLog(method string, errorMessage string) {
	stack := string(debug.Stack())
	stack = strings.ReplaceAll(stack, "\n\t", "\n")
	stacks := strings.Split(stack, "\n")

	log.Error().
		Str("type", typeSystem).
		Str("service", service).
		Str("method", method).
		Str("status", statusFailed).
		Dict("payload", zerolog.Dict().
			Str("message", errorMessage).
			Interface("stacktrace", stacks),
		).
		Msg("")
}

func PrintSuccessLog(method string, message string) {
	log.Info().
		Str("type", typeSystem).
		Str("service", service).
		Str("method", method).
		Str("status", statusSuccess).
		Dict("payload", zerolog.Dict().
			Str("message", message),
		).
		Msg("")
}
