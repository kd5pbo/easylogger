// Package easylogger provides an easy-to-use logging interface using the log
// package in the standard go distribution.  Basic usage is to call Generate
// to generate debug and verbose functions, and then to use the generated
// functions as needed.
//
//    /* Generate verbose and debug */
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
// calling LogVerbose or LogDebug, respectively.  LogVerbose and
// LogDebug may still be called later to change the behavior of the generated
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
	// def is the default LogSet used when the top-level functions (which
	// are wrappers for L's methods and variables) are called.
	def = new(LogSet)
)

// Generate verbose and debug functions.
//
// If makeFlags is true, the
// appropriate functions from the flag package will be called to add -verbose
// and -debug, with the effect of turning on verbose and debug output, as
// if LogVerbose LogDebug, respectively, had been called.
//
// Calling this function with makeFlags set to true after a call to any of the
// Log* functions will result in the values from the command line being
// restored.
func Generate(makeFlags bool) (verbose,
	debug func(format string, args ...interface{})) {
	/* Set flags if we're meant to */
	if makeFlags {
		def.verboseOn = flag.Bool("verbose", false, "Log verbosely")
		def.debugOn = flag.Bool("debug", false, "Log debugging messages")
	}

	return def.Verbose, def.Debug
}

// SetLogger causes l to be used for log output.  This may be nil to use the
// default logger.
func SetLogger(l *log.Logger) {
	def.SetLogger(l)
}

// LogVerbose turns on Verbose logging
// (verbose will log messages, debug won't).
func LogVerbose() { def.LogVerbose() }

// LogDebug turns on Debuging log messages
// (both verbose and debug will log messages).
func LogDebug() { def.LogDebug() }

/* LogNone turns off both verbose and debug logging. */
func LogNone() { def.LogNone() }

// LogdebugOnly logs verbose messages, but not debug messages
// (verbose will not log, debug will).
func LogDebugOnly() { def.LogDebugOnly() }

// Pause pauses logging.  Calls to Verbose and Debug will block until Resume
// is called.  Aside from being an excellent source of deadlocks, this allows
// for logfile rotation without risk of losing data.  See Resume for an
// example.
func Pause() {
	def.Pause()
}

// Resume resumes logging.  This should be called soon after Pause.  Pause and
// Resume can be used to safely change logfiles.
//
//    func changeLogFile(f string) {
//            easylogger.Pause()
//            defer easylogger.Resume()
//            o, err := os.OpenFile(f, os.O_CREATE|os.O_APPEND, 0644)
//            /* Error checking goes here */
//            log.SetOutput(o)
//            return
//    }
func Resume() {
	def.Resume()
}

// LogSet is a self-contained set of logging functions and variables.  It can
// be used to turn on and off logging for various parts of large programs.
type LogSet struct {
	verboseOn *bool       /* Enables verbose logging */
	debugOn   *bool       /* Enables debug logging */
	logger    *log.Logger /* Alternate logger (such as syslog). */
	changed   bool        /* One of the Log* functions has been called */
	m         *sync.Mutex /* Mutex held during writes */

}

// New returns a pointer to a new LogSet.
func New() *LogSet {
	/* Storage for the Ons */
	v := false
	d := false
	return &LogSet{
		verboseOn: &v,
		debugOn:   &d,
		logger:    nil,
		changed:   false,
		m:         &sync.Mutex{},
	}
}

/* Emit a message if doit is true */
func (l *LogSet) log(doit *bool, format string, args ...interface{}) {
	/* Do it only if we're supposed to do it */
	if nil == doit || !*doit {
		return
	}
	/* Work out which logger to use */
	if l.logger != nil { /* User-assigned logger */
		l.logger.Printf(format, args...)
	} else { /* Default logger */
		log.Printf(format, args...)
	}
}

/* Verbose logs a message if verbose messages are turned on */
func (l *LogSet) Verbose(format string, args ...interface{}) {
	doit := *l.verboseOn
	/* If the state hasn't been changed (i.e. set by the flags), verbose
	if debug is set */
	if !l.changed && !*l.verboseOn && *l.debugOn {
		doit = true
	}
	l.log(&doit, format, args...)
}

/* Debug logs a message if debugging messages are turned on */
func (l *LogSet) Debug(format string, args ...interface{}) {
	l.log(l.debugOn, format, args...)
}

/* logSwitch switches on/off verbose and debug logging */
func (l *LogSet) logSwitch(v, d bool) {
	/* Make sure we have bools allocated */
	if nil == l.verboseOn {
		b := false
		l.verboseOn = &b
	}
	if nil == l.debugOn {
		b := false
		l.debugOn = &b
	}
	/* Switch the switches */
	*l.verboseOn = v
	*l.debugOn = d
	/* Note there's been a change */
	l.changed = true
}

// LogVerbose turns on Verbose logging
// (verbose will log messages, debug won't).
func (l *LogSet) LogVerbose() { l.logSwitch(true, false) }

// LogDebug turns on Debuging log messages
// (both verbose and debug will log messages).
func (l *LogSet) LogDebug() { l.logSwitch(true, true) }

/* LogNone turns off both verbose and debug logging. */
func (l *LogSet) LogNone() { l.logSwitch(false, false) }

// LogdebugOnly logs verbose messages, but not debug messages
// (verbose will not log, debug will).
func (l *LogSet) LogDebugOnly() { l.logSwitch(false, true) }

// SetLogger causes logger to be used for log output.  This may be nil to use
// the default logger.
func (l *LogSet) SetLogger(logger *log.Logger) {
	l.logger = logger
}

// Pause pauses logging.  Calls to Verbose and Debug will block until Resume
// is called.  Aside from being an excellent source of deadlocks, this allows
// for logfile rotation without risk of losing data.  See Resume for an
// example.
func (l *Logset) Pause() {
	l.m.Lock()
}

// Resume resumes logging.  This should be called soon after Pause.
func (l *Logset) Resume() {
	l.m.Unlock()
}
