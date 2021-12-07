package format

import (
	"fmt"
	"strings"

	"github.com/harpyd/thestis/pkg/nerrors"
)

func SpecificationError(err error, indent int) string {
	return formatErrorWithStartIndent(err, indent, indent)
}

func formatErrorWithStartIndent(err error, indent, startIndent int) string {
	errMsg := nerrors.CommonError(err)

	nested := nerrors.NestedErrors(err)
	if len(nested) != 0 {
		errMsg = fmt.Sprintf("%s:", errMsg)
	}

	for _, e := range nested {
		errMsg = fmt.Sprintf(
			"%s\n%s%s",
			errMsg,
			strings.Repeat(" ", indent),
			formatErrorWithStartIndent(e, startIndent+indent, startIndent),
		)
	}

	return errMsg
}
