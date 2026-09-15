package payarc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ErrNotFound reports that PayArc has no record of the resource that was asked for. Match it
// with errors.Is rather than comparing status codes at the call site.
//
// This matters most on the charge lookups. PayArc answers an unknown charge id with an error
// body and a 404, but the response still decodes cleanly into the success type — leaving a
// zero-valued Charge whose Captured reads as "not captured" and whose ID is empty. A caller
// that took that at face value went on to void a charge with no id at all, which PayArc
// rejected as "The route v1/charges//void could not be found."
var ErrNotFound = errors.New("payarc has no record of the requested resource")

// StatusError is returned when PayArc answers with a status outside the success range. It
// carries PayArc's own explanation where the body parsed as one, and the raw body where it did
// not, so a caller can report something better than the status number alone.
type StatusError struct {
	// Op is what was being attempted, e.g. "get charge".
	Op         string
	StatusCode int

	// Message and Detail are PayArc's "message" and "error" fields; either may be empty.
	Message string
	Detail  string
	Fields  RequestErrorErrors

	// Body is the raw response body, set only when it did not parse as a PayArc error.
	Body string
}

func (e *StatusError) Error() string {
	reason := e.Message
	if reason == "" {
		reason = e.Detail
	}

	if reason == "" {
		reason = e.Body
	}

	if reason == "" {
		return fmt.Sprintf("%s failed with status %d", e.Op, e.StatusCode)
	}

	return fmt.Sprintf("%s failed with status %d: %s", e.Op, e.StatusCode, reason)
}

// Is reports StatusError as ErrNotFound on a 404 so callers can use errors.Is instead of
// reaching for the status code.
func (e *StatusError) Is(target error) bool {
	return target == ErrNotFound && e.StatusCode == http.StatusNotFound
}

// CheckResponse returns a *StatusError when resp carries a status outside the success range,
// and nil otherwise. The body is read and put back either way, so the caller can decode the
// success payload from resp.Body exactly as it would have.
//
// Call it after Do and before decoding. Without it a non-2xx body is decoded into the success
// type, and because the error shape shares no field names with the success shape the result is
// a zero value returned with a nil error — a failure that reads as a successful empty record.
//
// op names the operation for the error message, e.g. "get charge".
func CheckResponse(resp *http.Response, op string) error {
	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusIMUsed {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s failed with status %d and its body could not be read: %w", op, resp.StatusCode, err)
	}

	resp.Body = io.NopCloser(bytes.NewReader(body))

	statusErr := &StatusError{Op: op, StatusCode: resp.StatusCode}

	var reqErr RequestError
	if err := json.Unmarshal(body, &reqErr); err != nil {
		// Not a PayArc error body: a gateway or CDN page, most likely. Keep it, trimmed, so
		// the caller can tell an upstream outage from a PayArc business rejection.
		statusErr.Body = strings.TrimSpace(string(body))

		return statusErr
	}

	statusErr.Message = reqErr.Message
	statusErr.Detail = reqErr.Error
	statusErr.Fields = reqErr.Errors

	return statusErr
}
