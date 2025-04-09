package kotlin

import (
	"fmt"
	"io"
	"log"
	"strings"
)

var defaultIndent = "    "

func Serialize(w io.Writer, file *KtFile) (err error) {
	level := 0

	sep := ""
	for _, imp := range file.ImportList {
		_, err = fmt.Fprintf(w, "%s", sep)
		sep = "\n"
		_, err = fmt.Fprint(w, "import ")
		err = writeFqn(w, &imp.fqn)
	}

	_, err = fmt.Fprintln(w)

	for _, top := range file.TopLevels {
		switch top := top.(type) {
		case FuncDecl:
			fn := top

			if fn.Expr != nil {
				_, err = fmt.Fprintf(w, "\nfun %s() = ", fn.Name)
				err = writeExpr(w, &fn.Expr, level)
			} else {
				_, err = fmt.Fprintf(w, "\nfun %s() {\n", fn.Name)
				level++
				err = writeStatements(w, fn.Body.Statements, level)
				_, err = fmt.Fprintf(w, "\n}")
				level--
			}
		}
	}

	return
}

func writeFqn(w io.Writer, fqn *Fqn) (err error) {
	sep := ""

	for _, id := range *fqn {
		_, err = fmt.Fprint(w, sep)
		sep = "."
		_, err = fmt.Fprint(w, id)
	}
	return
}

func writeStatements(w io.Writer, statements []any, level int) (err error) {
	sep := ""
	for _, st := range statements {
		_, err = fmt.Fprint(w, sep)
		sep = "\n"

		if _, ok := st.(EmptyStatement); !ok {
			err = writeIdent(w, level)
		}

		err = writeStatement(w, &st, level)
	}

	return
}

func writeStatement(w io.Writer, st *any, level int) (err error) {
	switch st := (*st).(type) {
	case VarDecl:
		vd := st

		if vd.Mutable {
			_, err = fmt.Fprint(w, "var")
		} else {
			_, err = fmt.Fprint(w, "val")
		}

		_, err = fmt.Fprint(w, " ")
		_, err = fmt.Fprintf(w, "%s", vd.Name)

		if vd.Assignment != nil {
			_, err = fmt.Fprint(w, " ")
			_, err = fmt.Fprintf(w, "=")
			_, err = fmt.Fprint(w, " ")
			err = writeExpr(w, &vd.Assignment, level)
		}
	case PropAssignment:
		_, err = fmt.Fprintf(w, "%s = ", st.Prop)
		err = writeExpr(w, &st.Expr, level)
	case EmptyStatement:
		// Do nothing
	case ForInLoop:
		_, err = fmt.Fprint(w, "for (")

		switch bind := st.Bind.(type) {
		case PairDestruct:
			_, err = fmt.Fprint(w, "(")

			sep := ""
			for _, name := range bind {
				_, err = fmt.Fprint(w, sep)
				_, err = fmt.Fprint(w, name)
				sep = ", "
			}

			_, err = fmt.Fprint(w, ")")
		case string:
			_, err = fmt.Fprint(w, bind)
		}

		_, err = fmt.Fprint(w, " in ")
		err = writeExpr(w, &st.Expr, level)
		_, err = fmt.Fprint(w, ") ")
		err = writeBlock(w, st.Statements, level)
	default:
		err = writeExpr(w, &st, level)
	}
	return
}

func writeExpr(w io.Writer, expr *any, level int) (err error) {
	switch expr := (*expr).(type) {
	case CtorInvoke:
		err = writeFqn(w, (*Fqn)(&expr.Type))
		err = writeValueArgs(w, expr.ValueArgs, level)
	case MethodCall:
		switch rec := expr.Receiver.(type) {
		case string:
			_, err = fmt.Fprint(w, rec)
		default:
			err = writeExpr(w, &rec, level)
		}

		_, err = fmt.Fprintf(w, ".%s", expr.Method)
		err = writeValueArgs(w, expr.ValueArgs, level)
	case FuncCall:
		_, err = fmt.Fprintf(w, "%s", expr.Name)
		err = writeValueArgs(w, expr.ValueArgs, level)
	case PropAccess:
		_, err = fmt.Fprintf(w, "%s.%s", expr.Object, expr.Prop)
	case Id:
		_, err = fmt.Fprint(w, expr)
	case string:
		if strings.Contains(expr, "\n") {
			err = writeMultilineStr(w, expr, level)
		} else {
			err = writeStr(w, expr)
		}
	case bool:
		if expr {
			_, err = fmt.Fprint(w, "true")
		} else {
			_, err = fmt.Fprint(w, "false")
		}
	case LambdaLiteral:
		err = writeBlock(w, expr.Statements, level)
	case InlineLambdaLiteral:
		if len(expr.Statements) != 1 {
			log.Fatalf("expected 1 statement for InlineLambdaLiteral, got %d", len(expr.Statements))
		}

		_, err = fmt.Fprint(w, "{ ")
		err = writeStatement(w, &expr.Statements[0], 0)
		_, err = fmt.Fprint(w, " }")
	}

	return
}

func writeBlock(w io.Writer, statements []any, level int) (err error) {
	_, err = fmt.Fprint(w, "{\n")
	err = writeStatements(w, statements, level+1)
	_, err = fmt.Fprintln(w)
	err = writeIdent(w, level)
	_, err = fmt.Fprint(w, "}")
	return
}

func writeMultilineStr(w io.Writer, str string, level int) (err error) {
	lines := strings.Split(str, "\n")

	_, err = fmt.Fprintf(w, "\"\"\"%s", lines[0])

	for _, l := range lines[1 : len(lines)-1] {
		_, err = fmt.Fprint(w, "\n")
		err = writeIdent(w, level)
		_, err = fmt.Fprint(w, l)
	}
	_, err = fmt.Fprint(w, "\n")
	err = writeIdent(w, level)
	_, err = fmt.Fprintf(w, "%s\"\"\".trimIndent()", lines[len(lines)-1])
	return
}

func writeStr(w io.Writer, str string) (err error) {
	_, err = fmt.Fprint(w, "\"")
	for _, r := range str {
		if r != '"' && r != '\\' {
			_, err = fmt.Fprintf(w, "%c", r)
		} else {
			_, err = fmt.Fprintf(w, "\\%c", r)
		}
	}

	_, err = fmt.Fprint(w, "\"")
	return
}

func writeValueArgs(w io.Writer, args []any, level int) (err error) {
	writeSimpleArg := func(va any) (err error) {
		if arg, ok := va.(NamedArg); ok {
			_, err = fmt.Fprintf(w, "%s = ", arg.Name)
			err = writeExpr(w, &arg.Value, level)
		} else {
			err = writeExpr(w, &va, level)
		}

		return
	}

	onlyLambda := false
	if len(args) == 1 {
		_, onlyLambda = args[0].(LambdaLiteral)
	}

	if onlyLambda {
		_, err = fmt.Fprint(w, " ")
		err = writeExpr(w, &args[0], level)
	} else {
		_, err = fmt.Fprint(w, "(")

		if len(args) == 0 {
			_, err = fmt.Fprint(w, ")")
		} else {
			sep := ""
			var i int
			var va any
			for i, va = range args {
				if i == len(args)-1 { // Last
					_, isLambdaLiteral := va.(LambdaLiteral)
					_, isInlineLambdaLiteral := va.(InlineLambdaLiteral)
					if isLambdaLiteral || isInlineLambdaLiteral {
						_, err = fmt.Fprint(w, ") ")
						err = writeExpr(w, &va, level)
					} else {
						_, err = fmt.Fprintf(w, sep)
						sep = ", "
						err = writeSimpleArg(va)
						_, err = fmt.Fprint(w, ")")
					}
				} else {
					_, err = fmt.Fprintf(w, sep)
					sep = ", "
					err = writeSimpleArg(va)
				}
			}
		}
	}

	return
}

func writeIdent(w io.Writer, level int) (err error) {
	for i := 0; i < level; i++ {
		_, err = fmt.Fprintf(w, defaultIndent)
	}

	return
}
