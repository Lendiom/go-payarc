package payarc

import (
	"errors"
	"fmt"
	"strings"
)

// ErrMissingParameter reports that a required argument was blank or malformed, so no request
// was sent. Match it with errors.Is.
var ErrMissingParameter = errors.New("required parameter is missing")

// RequireParam rejects an identifier that cannot safely be put into a request, before the
// request is built. It returns an error wrapping ErrMissingParameter, naming the argument.
//
// Two things are rejected, and both fail the same way: the URL still forms, so the request goes
// out and lands somewhere other than the caller meant.
//
//   - Blank (empty, or only whitespace). An empty id collapses the path: Void("") asked PayArc
//     for v1/charges//void, which it answered with "The route v1/charges//void could not be
//     found." Worse is a blank id on a delete — Delete("") aims DELETE at v1/customers/, the
//     collection rather than one record.
//   - Containing a slash. A PayArc identifier is an opaque token and never contains one, so a
//     value that does is malformed data being spliced into the path, where it silently
//     retargets the request at a different endpoint.
func RequireParam(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w: %s", ErrMissingParameter, name)
	}

	if strings.Contains(value, "/") {
		return fmt.Errorf("%w: %s contains a slash, which would change which endpoint is called", ErrMissingParameter, name)
	}

	return nil
}
