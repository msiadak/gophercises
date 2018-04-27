package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/davecgh/go-spew/spew"
)

const defaultErrorMsg = "Something went wrong while serving your page."

type RecoverMux struct {
	*http.ServeMux
}

func NewRecoverMux(mux *http.ServeMux) *RecoverMux {
	return &RecoverMux{mux}
}

func (mux *RecoverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rrw := NewRecoverResponseWriter(w)
	logOutput := io.Writer(os.Stderr)
	if !prod {
		logOutput = io.MultiWriter(rrw, os.Stderr)
	}
	logger := log.New(logOutput, "", log.Flags())
	defer func() {
		if rec := recover(); rec != nil {
			rrw.buf.Reset()
			rrw.WriteHeader(http.StatusInternalServerError)
			if prod {
				rrw.Write([]byte(defaultErrorMsg))
			}
			logger.Printf("Panicked while handling %s\n", r.URL)
			logger.Println(rec)
			logger.Print(spew.Sdump(r))
			logger.Println("Stack Trace:")
			logger.Print(string(debug.Stack()))
			rrw.Flush()
		}
	}()
	mux.ServeMux.ServeHTTP(rrw, r)
	rrw.Flush()
}

type RecoverResponseWriter struct {
	http.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
}

func NewRecoverResponseWriter(w http.ResponseWriter) RecoverResponseWriter {
	rrw := RecoverResponseWriter{ResponseWriter: w}
	rrw.buf = new(bytes.Buffer)
	return rrw
}

func (rrw RecoverResponseWriter) Write(in []byte) (int, error) {
	return rrw.buf.Write(in)
}

func (rrw RecoverResponseWriter) WriteHeader(statusCode int) {
	rrw.statusCode = statusCode
}

func (rrw RecoverResponseWriter) Flush() {
	if rrw.statusCode != 0 {
		rrw.ResponseWriter.WriteHeader(rrw.statusCode)
	}
	rrw.buf.WriteTo(rrw.ResponseWriter)
}
