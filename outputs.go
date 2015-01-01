package statsreg

import (
	"os"
)

var newLine = []byte("\n")

type Output interface {
	Write(data []byte) error
}

// Output results to file
type File struct {
	path     string
	truncate bool
}

func (out *File) Write(data []byte) error {
	flags := os.O_CREATE | os.O_WRONLY
	if out.truncate {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_APPEND
	}

	f, err := os.OpenFile(out.path, flags, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return err
	}
	if _, err := f.Write(newLine); err != nil {
		return err
	}
	return nil
}
