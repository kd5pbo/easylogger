// Package easylogger provides an easy-to-use logging interface using the log
// package in the standard go distribution.  Basic usage is to call Generate()
// to generate debug and verbose functions, and then to use the generated
// functions as needed.
//
//    /* Generate verbose() and debug() */
//    var verbose, debug = easylogger.Generate(false)
//
//    func main(){
//
//            /* No logging is on by default */
//            verbose("This message will not be logged.")
//            debug("This message won't, either.")
//            log.Printf("This one still will, though.")
//
//            /* Turn on debug logging */
//            easylogger.LogDebug()
//            debug("Debugging messages will be logged.")
//            verbose("Verbose messages will be, too.")
//
//            /* Turn on verbose logging */
//            easylogger.LogVerbose()
//            verbose("Verbose messages will be logged.")
//            debug("Debugging messages will not be logged.")
//
//            /* Turn off esaylogger logging */
//            easylogger.LogNone()
//            debug("This message will not be logged.")
//            verbose("This one won't, either.")
//            log.Printf("This one still will, though.")
//    }
//
// Optionally, easylogger can add two flags, "debug" and "verbose" to the
// default set of flags if the flag package in the standard go distribution is
// being used:
//
//     verbose, debug := easyLogger.New(true)
//
// The program may be invoked with -verbose or -debug with the same effect as
// calling LogVerbose() or LogDebug(), respectively.  LogVerbose() and
// LogDebug() may still be called later to change the behavior of the generated
// functions.
//
// The generated functions take arguments in the same format as log.Printf
// (and indeed are wrappers around log.Printf).
//
//    wd, err := os.Getwd()
//    if nil != err {
//            verbose("Unable to determine working directory: %v", err)
//    } else {
//            debug("The current working directory is %v", wd)
//    }
package easylogger

import (
	"flag"
	"log"
)

/*
 * easylogger.go
 * Library to make for easy logging
 * by J. Stuart McMurray
 * created 20141218
 * last modified 20141218
 *
 * Copyright (c) 2012 J. Stuart McMurray. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *    * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *    * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *    * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

var (
	// L is the default LogSet used when the top-level functions (which
	// are wrappers for L's methods and variables) are called.
	L         = new(LogSet)
	madeFlags = false
)

// Generate verbose() and debug() functions.
//
// If makeFlags is true, the
// appropriate functions from the flag package will be called to add -verbose
// and -debug, with the effect of turning on verbose() and debug() output.
//
// Calling this function with makeFlags set to true after a call to any of the
// Log* functions will result in the values from the command line being
// restored.
func Generate(makeFlags bool) (verbose,
	debug func(format string, args ...interface{})) {
	/* Set flags if we're meant to */
	if makeFlags {
		L.VerboseOn = flag.Bool("verbose", false, "Log verbosely")
		L.DebugOn = flag.Bool("debug", false, "Log debugging messages")
	}

	return L.Verbose, L.Debug
}

// Set a specific logger (e.g. the output of syslog.New()) to be used.  May
// be nil to use the standard logger.
func SetLog(l *log.Logger) {
	L.Logger = l
}

/* logSwitch switches on/off verbose and debug logging */
func logSwitch(v, d bool) {
	/* Make sure we have bools allocated */
	if nil == L.VerboseOn {
		b := false
		L.VerboseOn = &b
	}
	if nil == L.DebugOn {
		b := false
		L.DebugOn = &b
	}
	/* Switch the switches */
	*L.VerboseOn = v
	*L.DebugOn = d
}

// LogVerbose turns on Verbose logging
// (verbose will log messages, debug won't).
func LogVerbose() { logSwitch(true, false) }

// LogDebug turns on Debuging log messages
// (both verbose and debug will log messages).
func LogDebug() { logSwitch(true, true) }

/* LogNone turns off both verbose and debug logging. */
func LogNone() { logSwitch(false, false) }

// LogDebugOnly logs verbose messages, but not debug messages
// (verbose will not log, debug will).
func LogDebugOnly() { logSwitch(false, true) }

// LogSet is a self-contained set of logging functions and variables.  It can
// be used to turn on and off logging for various parts of large programs.
type LogSet struct {
	VerboseOn *bool       /* Enables verbose logging */
	DebugOn   *bool       /* Enables debug logging */
	Logger    *log.Logger /* Alternate logger (such as syslog). */
}

// New returns a pointer to a new LogSet with both VerboseOn and DebugOn
// allocated.
func New() *LogSet {
	/* Storage for the Ons */
	v := false
	d := false
	return &LogSet{
		VerboseOn: &v,
		DebugOn:   &d,
	}
}

/* Emit a message if doit is true */
func (l *LogSet) log(doit *bool, format string, args ...interface{}) {
	/* Do it only if we're supposed to do it */
	if nil == doit || !*doit {
		return
	}
	/* Work out which logger to use */
	if l.Logger != nil { /* User-assigned logger */
		l.Logger.Printf(format, args...)
	} else { /* Default logger */
		log.Printf(format, args...)
	}
}

/* Verbose logs a message if verbose messages are turned on */
func (l *LogSet) Verbose(format string, args ...interface{}) {
	l.log(l.VerboseOn, format, args)
}

/* Debug logs a message if debugging messages are turned on */
func (l *LogSet) Debug(format string, args ...interface{}) {
	l.log(l.DebugOn, format, args)
}
