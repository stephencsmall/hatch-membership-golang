# HATCH Membership App Prototype
## Command Line / Web UI experiment written in Golang

This is a Golang app written with Beego framework.  Originally it was an interactive shell tool that offered a command line interface for maniupuating a simple membership database. I grew it to include the beginning of a simple web UI.

**Web UI probably doesn't work, it's a starting place for future development. Shell mode might work, mostly...**

### Am I in shell mode or Web UI mode?

Check line 198 of `main.go` and do some heinous toggling to decide...

### Running it
 Pretty easy, it's Go:

 ```
 go run main.go
 ```

 If you're in web mode, then open `http://localhost:8080` in a browser. Otherwise you're in a shell, so type "help" and press ENTER.