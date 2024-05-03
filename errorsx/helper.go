package errorsx

import (
	"errors"
	"net/http"
)

func HttpStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var p *XError
	if errors.As(err, &p) {
		return p.HttpStatusCode()
	}

	return http.StatusInternalServerError
}
