package re

import "github.com/coel-lang/coel/src/lib/core"

// Module is a module in the language.
var Module = map[string]core.Value{
	"find":    find,
	"match":   match,
	"replace": replace,
}
