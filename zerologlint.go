package zerologlint

import (
	"flag"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"

	"github.com/gostaticanalysis/comment/passes/commentmap"
)

var Analyzer = &analysis.Analyzer{
	Name: "zerologlint",
	Doc:  "Detects the wrong usage of `zerolog` that a user forgets to dispatch with `Send` or `Msg`",
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
		commentmap.Analyzer,
	},
}

func init() {
	Analyzer.Flags.Init("zerologlint", flag.ContinueOnError)
	Analyzer.Flags.String("prefix", "", "comma-separated list of additional module path prefixes to treat as zerolog (e.g., myorg/myzerolog)")
}

type posser interface {
	Pos() token.Pos
}

// callDefer is an interface just to hold both ssa.Call and ssa.Defer in our set
type callDefer interface {
	Common() *ssa.CallCommon
	Pos() token.Pos
}

type linter struct {
	// eventSet holds all the ssa block that is a zerolog.Event type instance
	// that should be dispatched.
	// Everytime the zerolog.Event is dispatched with Msg() or Send(),
	// deletes that block from this set.
	// At the end, check if the set is empty, or report the not dispatched block.
	eventSet    map[posser]struct{}
	// deleteLater holds the ssa block that should be deleted from eventSet after
	// all the inspection is done.
	// this is required because `else` ssa block comes after the dispatch of `if`` block.
	// e.g., if err != nil { log.Error() } else { log.Info() } log.Send()
	//       deleteLater takes care of the log.Info() block.
	deleteLater map[posser]struct{}
	recLimit    uint
	prefixes    []string
}

func run(pass *analysis.Pass) (interface{}, error) {
	srcFuncs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs

	prefixes := []string{"github.com/rs/zerolog"}
	if ff := pass.Analyzer.Flags.Lookup("prefix"); ff != nil {
		for _, p := range strings.Split(ff.Value.String(), ",") {
			if p = strings.TrimSpace(p); p != "" {
				prefixes = append(prefixes, p)
			}
		}
	}

	l := &linter{
		eventSet:    make(map[posser]struct{}),
		deleteLater: make(map[posser]struct{}),
		recLimit:    100,
		prefixes:    prefixes,
	}

	for _, sf := range srcFuncs {
		for _, b := range sf.Blocks {
			for _, instr := range b.Instrs {
				if c, ok := instr.(*ssa.Call); ok {
					l.inspect(c)
				} else if c, ok := instr.(*ssa.Defer); ok {
					l.inspect(c)
				}
			}
		}
	}

	// apply deleteLater to envetSet for else branches of if-else cases

	for k := range l.deleteLater {
		delete(l.eventSet, k)
	}

	// At the end, if the set is clear -> ok.
	// Otherwise, there must be a left zerolog.Event var that weren't dispatched. So report it.
	for k := range l.eventSet {
		pass.Reportf(k.Pos(), "must be dispatched by Msg or Send method")
	}

	return nil, nil
}

func (l *linter) inspect(cd callDefer) {
	c := cd.Common()

	// check if it's in github.com/rs/zerolog/log since there's some
	// functions in github.com/rs/zerolog that returns zerolog.Event
	// which should not be included. However, zerolog.Logger receiver is an exception.
	if l.isInLogPkg(*c) || l.isLoggerRecv(*c) {
		if l.isZerologEvent(c.Value) {
			// this ssa block should be dispatched afterwards at some point
			l.eventSet[cd] = struct{}{}
			return
		}
	}

	// if the call does not return zerolog.Event,
	// check if the base is zerolog.Event.
	// if so, check if the StaticCallee is Send() or Msg().
	// if so, remove the arg[0] from the set.
	f := c.StaticCallee()
	if f == nil {
		return
	}
	if !isDispatchMethod(f) {
		shouldReturn := true
		for _, p := range f.Params {
			if l.isZerologEvent(p) {
				// check if this zerolog.Event as a parameter is dispatched in the function
				// TODO: technically, it can be dispatched in another function that is called in this function, and
				//       this algorithm cannot track that. But I'm tired of thinking about that for now.
				for _, b := range f.Blocks {
					for _, instr := range b.Instrs {
						switch v := instr.(type) {
						case *ssa.Call:
							if l.inspectDispatchInFunction(v.Common()) {
								shouldReturn = false
								break
							}
						case *ssa.Defer:
							if l.inspectDispatchInFunction(v.Common()) {
								shouldReturn = false
								break
							}
						}
					}
				}
			}
		}
		if shouldReturn {
			return
		}
	}
	for _, arg := range c.Args {
		if l.isZerologEvent(arg) {
			// if there's branch, track both ways
			// this is for the case like:
			//   logger := log.Info()
			//   if err != nil {
			//     logger = log.Error()
			//   }
			//   logger.Send()
			//
			// Similar case like below goes to the same root but that doesn't
			// have any side effect.
			//   logger := log.Info()
			//   if err != nil {
			//     logger = logger.Str("a", "b")
			//   }
			//   logger.Send()
			if phi, ok := arg.(*ssa.Phi); ok {
				for _, edge := range phi.Edges {
					l.dfsEdge(edge, make(map[ssa.Value]struct{}), 0)
				}
			} else {
				val := l.getRootSsaValue(arg)
				delete(l.eventSet, val)
			}
		}
	}
}

func (l *linter) dfsEdge(v ssa.Value, visit map[ssa.Value]struct{}, cnt uint) {
	// only for safety
	if cnt > l.recLimit {
		return
	}
	cnt++

	if _, ok := visit[v]; ok {
		return
	}
	visit[v] = struct{}{}

	val := l.getRootSsaValue(v)
	phi, ok := val.(*ssa.Phi)
	if !ok {
		l.deleteLater[val] = struct{}{}
		return
	}
	for _, edge := range phi.Edges {
		l.dfsEdge(edge, visit, cnt)
	}
}

func (l *linter) inspectDispatchInFunction(cc *ssa.CallCommon) bool {
	if isDispatchMethod(cc.StaticCallee()) {
		for _, arg := range cc.Args {
			if l.isZerologEvent(arg) {
				return true
			}
		}
	}
	return false
}

func (l *linter) isInLogPkg(c ssa.CallCommon) bool {
	switch v := c.Value.(type) {
	case ssa.Member:
		p := v.Package()
		if p == nil {
			return false
		}
		for _, prefix := range l.prefixes {
			if strings.HasSuffix(p.Pkg.Path(), prefix+"/log") {
				return true
			}
		}
	}
	return false
}

func (l *linter) isLoggerRecv(c ssa.CallCommon) bool {
	switch f := c.Value.(type) {
	case *ssa.Function:
		if recv := f.Signature.Recv(); recv != nil {
			ts := types.TypeString(recv.Type(), nil)
			for _, prefix := range l.prefixes {
				if strings.HasSuffix(ts, prefix+".Logger") {
					return true
				}
			}
		}
	}
	return false
}

func (l *linter) isZerologEvent(v ssa.Value) bool {
	ts := v.Type().String()
	for _, prefix := range l.prefixes {
		if strings.HasSuffix(ts, prefix+".Event") {
			return true
		}
	}
	return false
}

func isDispatchMethod(f *ssa.Function) bool {
	if f == nil {
		return false
	}
	m := f.Name()
	if m == "Send" || m == "Msg" || m == "Msgf" || m == "MsgFunc" {
		return true
	}
	return false
}

func (l *linter) getRootSsaValue(v ssa.Value) ssa.Value {
	if c, ok := v.(*ssa.Call); ok {
		v := c.Value()

		// When there is no receiver, that's the block of zerolog.Event
		// eg. Error() method in log.Error().Str("foo", "bar").Send()
		if len(v.Call.Args) == 0 {
			return v
		}

		// Even when there is a receiver, if it's a zerolog.Logger instance, return this block
		// eg. Info() method in zerolog.New(os.Stdout).Info()
		root := v.Call.Args[0]
		if !l.isZerologEvent(root) {
			return v
		}

		// Ok to just return the receiver because all the method in this
		// chain is zerolog.Event at this point.
		return l.getRootSsaValue(root)
	}
	return v
}
