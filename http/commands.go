package http

import (
	"bufio"
	"context"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/filebrowser/filebrowser/v2/runner"
)

const (
	WSWriteDeadline = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  16384,  // Increased from 1024 to 16KB
	WriteBufferSize: 16384,  // Increased from 1024 to 16KB
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow cross-origin for better compatibility
	},
}

var (
	cmdNotAllowed = []byte("Command not allowed.")
)

//nolint:unparam
func wsErr(ws *websocket.Conn, r *http.Request, status int, err error) {
	txt := http.StatusText(status)
	if err != nil || status >= 400 {
		log.Printf("%s: %v %s %v", r.URL.Path, status, r.RemoteAddr, err)
	}
	if err := ws.WriteControl(websocket.CloseInternalServerErr, []byte(txt), time.Now().Add(WSWriteDeadline)); err != nil {
		log.Print(err)
	}
}

var commandsHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	var raw string

	for {
		_, msg, err := conn.ReadMessage() //nolint:govet
		if err != nil {
			wsErr(conn, r, http.StatusInternalServerError, err)
			return 0, nil
		}

		raw = strings.TrimSpace(string(msg))
		if raw != "" {
			break
		}
	}

	// Fail fast
	if !d.server.EnableExec || !d.user.Perm.Execute {
		if err := conn.WriteMessage(websocket.TextMessage, cmdNotAllowed); err != nil { //nolint:govet
			wsErr(conn, r, http.StatusInternalServerError, err)
		}

		return 0, nil
	}

	command, _, err := runner.ParseCommand(d.settings, raw)
	if err != nil {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil { //nolint:govet
			wsErr(conn, r, http.StatusInternalServerError, err)
		}
		return 0, nil
	}

	// Allow all commands when user has execute permission
	// Note: Removed restrictive command whitelist check to allow all commands

	// Create context with timeout to prevent hanging commands
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, command[0], command[1:]...) //nolint:gosec
	cmd.Dir = d.user.FullPath(r.URL.Path)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		wsErr(conn, r, http.StatusInternalServerError, err)
		return 0, nil
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		wsErr(conn, r, http.StatusInternalServerError, err)
		return 0, nil
	}

	if err := cmd.Start(); err != nil {
		wsErr(conn, r, http.StatusInternalServerError, err)
		return 0, nil
	}

	// Use goroutines for concurrent reading of stdout and stderr
	var wg sync.WaitGroup
	wg.Add(2)

	// Buffer for more efficient output handling
	bufferSize := 8192

	go func() {
		defer wg.Done()
		buffer := make([]byte, bufferSize)
		for {
			n, err := stdout.Read(buffer)
			if n > 0 {
				if err := conn.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
					log.Print(err)
					return
				}
			}
			if err != nil {
				if err != io.EOF {
					log.Print(err)
				}
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		buffer := make([]byte, bufferSize)
		for {
			n, err := stderr.Read(buffer)
			if n > 0 {
				if err := conn.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
					log.Print(err)
					return
				}
			}
			if err != nil {
				if err != io.EOF {
					log.Print(err)
				}
				return
			}
		}
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			conn.WriteMessage(websocket.TextMessage, []byte("\nCommand timed out after 30 minutes\n"))
		} else {
			// Only log error if it's not a normal exit code
			if exitErr, ok := err.(*exec.ExitError); ok {
				conn.WriteMessage(websocket.TextMessage, []byte("\nCommand exited with code: "+string(rune(exitErr.ExitCode()))+"\n"))
			}
		}
	}

	return 0, nil
})
