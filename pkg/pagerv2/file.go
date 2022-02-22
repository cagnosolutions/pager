package pagerv2

import (
	"os"
	"path/filepath"
)

func OpenFile(path string) (*os.File, error) {
	// sanitize path
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	// split path
	dir, name := filepath.Split(filepath.ToSlash(path))
	// init PageManager and dirs
	var fp *os.File
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// create dir
		err = os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return nil, err
		}
		// create PageManager
		fp, err = os.Create(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		// close PageManager
		err = fp.Close()
		if err != nil {
			return nil, err
		}
	}
	// open existing PageManager
	fp, err = os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return fp, nil
}
