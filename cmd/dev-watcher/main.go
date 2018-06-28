package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	pkg "github.com/xinpianchang/dev-watcher/watcher"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	Build = "devel"

	dir     = flag.String("d", ".", "folder to watch.")
	filters = flag.String("f", "*", "filter file extension, multiple extensions separated by commas.")
	wait    = flag.Int64("t", 2000, "postpone shell execution until after wait milliseconds.")
	script  = flag.String("s", "./.dev-watcher.sh", "shell script file which executed after file changed.")

	V = flag.Bool("version", false, "show version")
	H = flag.Bool("help", false, "show help")

	cmd     *exec.Cmd
	cmdLock sync.Mutex

	l = log.New(os.Stdout, "[dev-watcher] ", log.LstdFlags)
)

func main() {
	flag.Parse()
	if !flag.Parsed() {
		flag.Usage()
		return
	}

	if *H {
		flag.Usage()
		return
	}

	if *V {
		fmt.Println("dev-watcher")
		fmt.Println("  git version:", Build)
		fmt.Println("  go version :", runtime.Version())
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		l.Println("create watcher error:", err)
	}
	defer watcher.Close()

	debounceShellExecutor := pkg.NewDebounce(time.Millisecond*time.Duration(*wait), shellExecutor)

	go func() {
		fileFilter := newFilter(*filters)

		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Rename == fsnotify.Rename ||
					event.Op&fsnotify.Write == fsnotify.Write {

					if info, err := os.Stat(event.Name); err == nil {
						if info.IsDir() {
							continue
						}
					}

					if fileFilter(event.Name) {
						debounceShellExecutor()
					}
				}
			}
		}
	}()

	filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				l.Println("add watcher error:", err)
			}
		}

		return nil
	})

	l.Println("start watch:", *dir)

	// first run
	debounceShellExecutor()

	select {}
}

type filter func(file string) bool

func newFilter(filters string) filter {
	if filters == "" || filters == "*" {
		return func(file string) bool {
			return true
		}
	}

	exts := strings.Split(filters, ",")

	return func(file string) bool {
		for _, ext := range exts {
			if strings.HasSuffix(file, ext) {
				return true
			}
		}

		return false
	}
}

func shellExecutor() {
	cmdLock.Lock()
	defer cmdLock.Unlock()

	if cmd != nil {
		state := cmd.ProcessState
		if state != nil && !state.Exited() {
			log.Println("kill previous shell process pid:", cmd.Process.Pid)
			cmd.Process.Kill()
			cmd.Process = nil
		}
	}

	cmd = exec.Command(*script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		err := cmd.Run()
		if err != nil {
			l.Println(err)
		}
	}()
}
