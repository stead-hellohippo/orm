package script

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/svett/gom"
)

var _ FileGenerator = &Generator{}

type Generator struct {
	Dir string
}

func (g *Generator) Create(container, name string) (string, error) {
	if err := os.MkdirAll(g.Dir, 0700); err != nil {
		return "", err
	}

	provider := &gom.CmdProvider{
		Repository: make(map[string]string),
	}

	if err := provider.LoadDir(g.Dir); err != nil {
		return "", err
	}

	if _, err := provider.Command(name); err == nil {
		return "", fmt.Errorf("Command '%s' already exists", name)
	}

	if container == "" {
		container = time.Now().Format(format)
	}

	path := filepath.Join(g.Dir, fmt.Sprintf("%s.sql", container))
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return "", err
	}

	defer func() {
		if ioErr := file.Close(); err != nil {
			path = ""
			err = ioErr
		}
	}()

	fmt.Fprintln(file, "-- Auto-generated at", time.Now().Format(time.UnixDate))
	fmt.Fprintf(file, "-- name: %s", name)
	fmt.Fprintln(file)
	fmt.Fprintln(file)

	return path, err
}
