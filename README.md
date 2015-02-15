easylogger
==========

A small go library providing verbose() and debug() printf-like functions.  Aimed to be easy to use.

Please see https://godoc.org/github.com/kd5pbo/easylogger for documentation.

Not very well tested, but tested.  Test well before use.

Quickstart
----------
```go
var verbose, debug = easylogger.Generate(true)
func main() {
        
        /* Enable verbose logging */
        easylogger.VerboseOn()
        log.Printf("This message is logged")
        verbose("This message is logged, too")
        debug("This message is not")
        
        /* Enable debug logging */
        easylogger.DebugOn()
        log.Printf("This message is logged")
        verbose("This message is, too")
        debug("So is this one")
        
        /* Switch off logging */
        easylogger.LogNone()
        log.Printf("This message is logged")
        verbose("This one isn't")
        debug("Nor is this one")
}
```
