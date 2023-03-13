package zerologlint

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"

	"github.com/gostaticanalysis/comment/passes/commentmap"
)

var Analyzer = &analysis.Analyzer{
	Name: "zerologlinter",
	Doc:  "finds cases where zerolog methods are not followed by Msg or Send",
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
		commentmap.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	srcFuncs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs

	// This map holds all the ssa block that is a zerolog.Event type instance.
	// Everytime the zerolog.Event is dispatched with Msg() or Send(),
	// deletes that block from this map.
	// At the end, check if the set is empty, or report the not dispatched block.
	set := make(map[ssa.Value]struct{})

	for _, sf := range srcFuncs {
		for _, b := range sf.Blocks {
			for _, instr := range b.Instrs {
				if c, ok := instr.(*ssa.Call); ok {
					v := c.Value()
					if isInLogPkg(v) {
						if isZerologEvent(v) {
							// check if this is a new instance of zerolog.Event like logger := log.Error()
							// which should be dispatched afterwards at some point
							if len(v.Call.Args) == 0 {
								set[v] = struct{}{}
							}
							continue
						}
					}

					// if the call does not return zerolog.Event,
					// check if the base is zerolog.Event.
					// if so, check if the StaticCallee is Send() or Msg().
					// if so, remove the arg[0] from the set.
					for _, arg := range v.Call.Args {
						if isZerologEvent(arg) {
							if isDispatchMethod(v) {
								val := getRootSsaValue(arg)
								// if there's branch, remove both ways from the set
								if phi, ok := val.(*ssa.Phi); ok {
									for _, edge := range phi.Edges {
										delete(set, edge)
									}
								} else {
									delete(set, val)
								}
							}
						}
					}
				}
			}
		}
	}
	// At the end, if the set is clear -> ok.
	// if the set is not clear -> there must be a left zerolog.Event var that weren't dispached.
	// -> Report it using position.
	for k := range set {
		pass.Reportf(k.Pos(), "missing Msg or Send call for zerolog log method")
	}
	return nil, nil
}

func isInLogPkg(c *ssa.Call) bool {
	switch v := c.Call.Value.(type) {
	case ssa.Member:
		p := removeVendor(v.Package().Pkg.Path())
		return p == "github.com/rs/zerolog/log"
	default:
		return false
	}
}

func isZerologEvent(c ssa.Value) bool {
	ts := c.Type().String()
	t := removeVendor(ts)
	return t == "github.com/rs/zerolog.Event"
}

// RemoVendor removes vendoring information from import path.
func removeVendor(path string) string {
	i := strings.Index(path, "vendor/")
	if i >= 0 {
		return path[i+len("vendor/"):]
	}
	return path
}

func isDispatchMethod(c *ssa.Call) bool {
	m := c.Common().StaticCallee().Name()
	if m == "Send" || m == "Msg" {
		return true
	}
	return false
}

func getRootSsaValue(arg ssa.Value) ssa.Value {
	if c, ok := arg.(*ssa.Call); ok {
		v := c.Value()
		// When there is no receiver, that's the block of zerolog.Event
		// eg. Error() method in log.Error().Str("foo", "bar"). Send()
		if len(v.Call.Args) == 0 {
			return v
		}

		// Ok to just return the receiver because all the method in this
		// chain is zerolog.Event at this point.
		return getRootSsaValue(v.Call.Args[0])
	}
	return arg
}
