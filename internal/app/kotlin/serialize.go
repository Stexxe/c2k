package kotlin

import (
	"fmt"
	"io"
)

var defaultIndent = "    "

func Serialize(w io.Writer, file *KtFile) (err error) {
	level := 0

	sep := ""
	for _, imp := range file.ImportList {
		_, err = fmt.Fprintf(w, "%s", sep)
		sep = "\n"
		_, err = fmt.Fprint(w, "import ")
		err = writeFqn(w, &imp.Id)
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
		for i := 0; i < level; i++ {
			_, err = fmt.Fprintf(w, defaultIndent)
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
	default:
		err = writeExpr(w, &st, level)
	}
	return
}

func writeExpr(w io.Writer, expr *any, level int) (err error) {
	switch expr := (*expr).(type) {
	case CtorInvoke:
		err = writeFqn(w, (*Fqn)(&expr.Type))
		_, err = fmt.Fprint(w, "()")
	case CallExpr:
		if expr.Receiver != "" {
			_, err = fmt.Fprintf(w, "%s.%s", expr.Receiver, expr.Method)
		} else {
			_, err = fmt.Fprintf(w, "%s", expr.Method)
		}

		onlyLambda := false
		if len(expr.ValueArgs) == 1 {
			_, onlyLambda = expr.ValueArgs[0].(LambdaLiteral)
		}

		if onlyLambda {
			err = writeExpr(w, &expr.ValueArgs[0], level)
		} else {
			_, err = fmt.Fprint(w, "(")
			sep := ""
			for _, va := range expr.ValueArgs {
				_, err = fmt.Fprintf(w, sep)
				sep = ", "
				err = writeExpr(w, &va, level)
			}

			_, err = fmt.Fprint(w, ")")
		}
	case StringLiteral:
		_, err = fmt.Fprintf(w, "\"%s\"", expr) // TODO: Escape " symbol
	case LambdaLiteral:
		_, err = fmt.Fprint(w, " {\n")
		err = writeStatements(w, expr.Statements, level+1)
		_, err = fmt.Fprint(w, "\n}")
	}

	return
}
