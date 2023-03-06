package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

// ErrorPrinter implements xerrors.Printer interface
type ErrorPrinter struct {
	*bytes.Buffer
}

func (p *ErrorPrinter) Detail() bool {
	return true
}
func (p *ErrorPrinter) Print(args ...interface{}) {
	fmt.Fprint(p.Buffer, args...)
}
func (p *ErrorPrinter) Printf(format string, args ...interface{}) {
	fmt.Fprintf(p.Buffer, format, args...)
}

type ErrorJson struct {
	Message  string          `json:"message"`
	Type     string          `json:"type"`
	Function string          `json:"function,omitempty"`
	File     string          `json:"file,omitempty"`
	Original json.RawMessage `json:"original,omitempty"`
}

func EncodeErrorToJSON(err error) ([]byte, error) {
	if err == nil {
		return nil, nil
	}

	var original error
	if w, ok := err.(xerrors.Wrapper); ok {
		original = w.Unwrap()
	}

	var originalBytes []byte
	if original != nil {
		var eerr error
		originalBytes, eerr = EncodeErrorToJSON(original)
		if eerr != nil {
			return nil, xerrors.Errorf("encoding original error: %w", err)
		}
	}

	var p = &ErrorPrinter{
		Buffer: &bytes.Buffer{},
	}
	if f, ok := err.(xerrors.Formatter); ok {
		f.FormatError(p)
	}

	errMsg := err.Error()
	if original != nil {
		if strings.HasSuffix(errMsg, original.Error()) {
			errMsg = errMsg[:len(errMsg)-len(original.Error())]
			errMsg = strings.TrimSuffix(errMsg, ": ")
		}
	}

	errLoc := p.String()
	errLoc = strings.TrimPrefix(errLoc, errMsg)

	errLocs := strings.Split(errLoc, "\n")
	var errFunc, errFile string
	if len(errLocs) > 0 {
		errFunc = strings.TrimSpace(errLocs[0])
	}
	if len(errLocs) > 1 {
		errFile = strings.TrimSpace(errLocs[1])
	}

	var ejson = ErrorJson{
		Message:  errMsg,
		Type:     fmt.Sprintf("%T", err),
		Function: errFunc,
		File:     errFile,
		Original: json.RawMessage(originalBytes),
	}

	return json.MarshalIndent(ejson, "", "  ")
}
