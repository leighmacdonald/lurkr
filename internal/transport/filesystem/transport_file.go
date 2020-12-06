package filesystem

import (
	"github.com/leighmacdonald/golib"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

type fileTransport struct {
}

func (t *fileTransport) Send(reader io.Reader, path string) error {
	dir := filepath.Dir(path)
	if !golib.Exists(dir) {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return errors.Wrapf(err, "Failed to create base directory for dest")
		}
	}
	dst, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "Could not create dest file")
	}
	if _, err := io.Copy(dst, reader); err != nil {
		return errors.Wrapf(err, "Failed to write file")
	}
	log.Infof("Copied file successfully: %v", path)
	return nil
}

func NewFileTransport() (*fileTransport, error) {
	return &fileTransport{}, nil
}
