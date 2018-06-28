# dev watcher

dev-watcher is a tool using for monitoring source folder then executing scripts (eg bash) while file changes.

## usage

```
go get -u github.com/xinpianchang/dev-watcher
```

```
➜  ~ dev-watcher -h

Usage of dev-watcher:
  -d string
    	folder to watch. (default ".")
  -f string
    	filter file extension, multiple extensions separated by commas. (default "*")
  -help
    	show help
  -s string
    	shell script file which executed after file changed. (default "./.dev-watcher.sh")
  -t int
    	postpone shell execution until after wait milliseconds. (default 2000)
  -version
    	show version
```