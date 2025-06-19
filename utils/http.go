package utils

import (
	"log/slog"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	CodeMapper = map[codes.Code]int {
		codes.OK: 200,
		codes.InvalidArgument: 400,
		codes.Unauthenticated: 401,
		codes.PermissionDenied: 403,
		codes.NotFound: 404,
		codes.AlreadyExists: 409,
		codes.Internal: 500,
	}
)

func HandleResponseErr(w http.ResponseWriter, logger *slog.Logger, msg string, err error) {
	code := 500
	st, ok := status.FromError(err)
	if !ok {
		logger.Warn("Non-gRPC error", slog.Any("error", err))
		http.Error(w, err.Error(), code)
		return 
	}
	code = CodeMapper[st.Code()]
	http.Error(w, msg + st.Message(), code)
}

func GetClientIp(r *http.Request) string {
	ipList := r.Header.Get("X-Forwarded-For")
	if ipList != "" {
		return strings.Split(ipList, ",")[0]
	}
	return r.RemoteAddr
}
