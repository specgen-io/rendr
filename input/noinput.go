package input

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/specgen-io/rendr/blueprint"
	"github.com/specgen-io/rendr/values"
)

func NoInput(arg blueprint.NamedArg) (values.ArgValue, error) {
	return nil, errors.New(fmt.Sprintf(`no value provided for argument: "%s"`, arg.Name))
}
