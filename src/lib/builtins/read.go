package builtins

import (
	"io/ioutil"
	"os"

	"github.com/coel-lang/coel/src/lib/core"
)

// Read reads a string from stdin or a file.
var Read = core.NewLazyFunction(
	core.NewSignature(
		nil, []core.OptionalArgument{core.NewOptionalArgument("file", core.Nil)}, "",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		v := ts[0].Eval()
		file := os.Stdin

		if s, ok := v.(core.StringType); ok {
			var err error
			file, err = os.Open(string(s))

			if err != nil {
				return readError(err)
			}
		} else if _, ok := v.(core.NilType); !ok {
			s, err := core.StrictDump(v)

			if err != nil {
				return err
			}

			return core.ValueError(
				"file optional argument's value must be nil or a filename. Got %s.",
				s)
		}

		s, err := ioutil.ReadAll(file)

		if err != nil {
			return readError(err)
		}

		return core.NewString(string(s))
	})

func readError(err error) *core.Thunk {
	return core.NewError("ReadError", err.Error())
}