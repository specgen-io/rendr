package input

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/specgen-io/rendr/blueprint"
)

func NoInput(arg blueprint.NamedArg) (blueprint.ArgValue, error) {
	return nil, errors.New(fmt.Sprintf(`no value provided for argument: "%s"`, arg.Name))
}
