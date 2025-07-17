package fp

import (
	"fmt"
)

func WrapErrors(err error, causes ...error) error {
	if err == nil && len(causes) == 0 {
		return nil
	}
	if err == nil && len(causes) > 0 {
		err = causes[0]
		causes = causes[1:]
	}
	for _, cause := range causes {
		err = fmt.Errorf("%w: %w", err, cause)
	}
	return err
}
