package impl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/model"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Service) GetRecovererMiddleware(f model.Middleware, options ...model.MiddlewareOption) model.Middleware {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					middleware.PrintPrettyStack(rvr)
				}

				switch t := rvr.(type) {
				case string:
					err = errorsg.WithOptions(errors.New(t))
				case error:
					err = t
				default:
					err = errors.New("unknown panic error")
				}
				err = errorsg.WithOptions(err,
					errorsg.WithType(errorsg.ErrorTypePanic),
					errorsg.WithStack(PrintPrettyStack(rvr, debug.Stack())),
				)
			}
		}()

		return f(ctx, rw, r)
	}
}

func PrintPrettyStack(rvr interface{}, debugStack []byte) (prettyStackByte []byte) {
	s := prettyStack{}
	out, err := s.parse(debugStack, rvr)
	if err == nil {
		return out
	} else {
		return debugStack
	}
}

// from /vendor/github.com/go-chi/chi/v5/middleware/recoverer.go
type prettyStack struct {
}

func (s prettyStack) parse(debugStack []byte, rvr interface{}) ([]byte, error) {
	var err error
	useColor := true
	buf := &bytes.Buffer{}

	cW(buf, false, bRed, "\n")
	cW(buf, useColor, bCyan, " panic: ")
	cW(buf, useColor, bBlue, "%v", rvr)
	cW(buf, false, bWhite, "\n \n")

	// process debug stack info
	stack := strings.Split(string(debugStack), "\n")
	lines := []string{}

	// locate panic line, as we may have nested panics
	for i := len(stack) - 1; i > 0; i-- {
		lines = append(lines, stack[i])
		if strings.HasPrefix(stack[i], "panic(0x") {
			lines = lines[0 : len(lines)-2] // remove boilerplate
			break
		}
	}

	// reverse
	for i := len(lines)/2 - 1; i >= 0; i-- {
		opp := len(lines) - 1 - i
		lines[i], lines[opp] = lines[opp], lines[i]
	}

	// decorate
	for i, line := range lines {
		lines[i], err = s.decorateLine(line, useColor, i)
		if err != nil {
			return nil, err
		}
	}

	for _, l := range lines {
		fmt.Fprintf(buf, "%s", l)
	}
	return buf.Bytes(), nil
}

func (s prettyStack) decorateLine(line string, useColor bool, num int) (string, error) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "\t") || strings.Contains(line, ".go:") {
		return s.decorateSourceLine(line, useColor, num)
	} else if strings.HasSuffix(line, ")") {
		return s.decorateFuncCallLine(line, useColor, num)
	} else {
		if strings.HasPrefix(line, "\t") {
			return strings.Replace(line, "\t", "      ", 1), nil
		} else {
			return fmt.Sprintf("    %s\n", line), nil
		}
	}
}

func (s prettyStack) decorateFuncCallLine(line string, useColor bool, num int) (string, error) {
	idx := strings.LastIndex(line, "(")
	if idx < 0 {
		return "", errors.New("not a func call line")
	}

	buf := &bytes.Buffer{}
	pkg := line[0:idx]
	// addr := line[idx:]
	method := ""

	idx = strings.LastIndex(pkg, string(os.PathSeparator))
	if idx < 0 {
		idx = strings.Index(pkg, ".")
		method = pkg[idx:]
		pkg = pkg[0:idx]
	} else {
		method = pkg[idx+1:]
		pkg = pkg[0 : idx+1]
		idx = strings.Index(method, ".")
		pkg += method[0:idx]
		method = method[idx:]
	}
	pkgColor := nYellow
	methodColor := bGreen

	if num == 0 {
		cW(buf, useColor, bRed, " -> ")
		pkgColor = bMagenta
		methodColor = bRed
	} else {
		cW(buf, useColor, bWhite, "    ")
	}
	cW(buf, useColor, pkgColor, "%s", pkg)
	cW(buf, useColor, methodColor, "%s\n", method)
	// cW(buf, useColor, nBlack, "%s", addr)
	return buf.String(), nil
}

func (s prettyStack) decorateSourceLine(line string, useColor bool, num int) (string, error) {
	idx := strings.LastIndex(line, ".go:")
	if idx < 0 {
		return "", errors.New("not a source line")
	}

	buf := &bytes.Buffer{}
	path := line[0 : idx+3]
	lineno := line[idx+3:]

	idx = strings.LastIndex(path, string(os.PathSeparator))
	dir := path[0 : idx+1]
	file := path[idx+1:]

	idx = strings.Index(lineno, " ")
	if idx > 0 {
		lineno = lineno[0:idx]
	}
	fileColor := bCyan
	lineColor := bGreen

	if num == 1 {
		cW(buf, useColor, bRed, " ->   ")
		fileColor = bRed
		lineColor = bMagenta
	} else {
		cW(buf, false, bWhite, "      ")
	}
	cW(buf, useColor, bWhite, "%s", dir)
	cW(buf, useColor, fileColor, "%s", file)
	cW(buf, useColor, lineColor, "%s", lineno)
	if num == 1 {
		cW(buf, false, bWhite, "\n")
	}
	cW(buf, false, bWhite, "\n")

	return buf.String(), nil
}

var (
	// Normal colors
	// nBlack   = []byte{'\033', '[', '3', '0', 'm'}
	// nRed     = []byte{'\033', '[', '3', '1', 'm'}
	// nGreen   = []byte{'\033', '[', '3', '2', 'm'}
	nYellow = []byte{'\033', '[', '3', '3', 'm'}
	// nBlue    = []byte{'\033', '[', '3', '4', 'm'}
	// nMagenta = []byte{'\033', '[', '3', '5', 'm'}
	// nCyan    = []byte{'\033', '[', '3', '6', 'm'}
	// nWhite   = []byte{'\033', '[', '3', '7', 'm'}

	// Bright colors
	// bBlack   = []byte{'\033', '[', '3', '0', ';', '1', 'm'}
	bRed   = []byte{'\033', '[', '3', '1', ';', '1', 'm'}
	bGreen = []byte{'\033', '[', '3', '2', ';', '1', 'm'}
	// bYellow  = []byte{'\033', '[', '3', '3', ';', '1', 'm'}
	bBlue    = []byte{'\033', '[', '3', '4', ';', '1', 'm'}
	bMagenta = []byte{'\033', '[', '3', '5', ';', '1', 'm'}
	bCyan    = []byte{'\033', '[', '3', '6', ';', '1', 'm'}
	bWhite   = []byte{'\033', '[', '3', '7', ';', '1', 'm'}

	reset = []byte{'\033', '[', '0', 'm'}
)

var IsTTY bool

// colorWrite
func cW(w io.Writer, useColor bool, color []byte, s string, args ...interface{}) {
	if IsTTY && useColor {
		w.Write(color)
	}
	fmt.Fprintf(w, s, args...)
	if IsTTY && useColor {
		w.Write(reset)
	}
}
