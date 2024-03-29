package utils

import (
	"errors"
	"io/fs"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/debdutdeb/gopark/pkg/progressbar"
)

var ErrNotFound = errors.New("not found")

// DownloadSilent is your normal file download path, shows no progress
func DownloadSilent(url string, path string) error {
	return download(url, path, false, "")
}

// DownloadWithProgressBar shows a very basic progress bar for the file
func DownloadWithProgressBar(label, url, path string) error {
	return download(url, path, true, label)
}

func download(url, path string, showProgress bool, label string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	dir := filepath.Dir(path)

	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}

		err := os.MkdirAll(dir, 0750)
		if err != nil {
			return err
		}
	}

	w, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0750)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if showProgress {
		bar, err := progressbar.NewWriteProgressBar(label, resp.ContentLength, w, nil)
		if err != nil {
			return err
		}

		_, err = io.Copy(bar, resp.Body)
		if err != nil {
			return err
		}

		return w.Sync()
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return w.Sync()
}

// DumbInstall copies files and directories from source to destination
// follows symlinks
func DumbInstall(dst, src string) (err error) {
	var links map[string]string = make(map[string]string)

	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)

		info, _ := d.Info()

		if info.Mode()&fs.ModeSymlink == fs.ModeSymlink {
			linkSrc, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}

			if relSrc, err := filepath.Rel(src, linkSrc); err == nil {
				links[filepath.Join(dst, relSrc)] = filepath.Join(dst, rel)
				return nil
			}
		}

		if d.IsDir() {
			if err := os.MkdirAll(filepath.Join(dst, rel), info.Mode().Perm()); err != nil {
				return err
			}
			return nil
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}

		defer srcFile.Close()

		dstFile, err := os.OpenFile(filepath.Join(dst, rel), os.O_CREATE|os.O_WRONLY, info.Mode().Perm())
		if err != nil {
			return err
		}

		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}

		return dstFile.Sync()
	})

	if err != nil {
		return
	}

	for src, dst := range links {
		if err = os.Symlink(src, dst); err != nil && os.IsExist(err) {
			os.Remove(dst)
			err = os.Symlink(src, dst)
			if err != nil {
				break
			}
		} else if err != nil {
			break
		}
	}

	return
}

// MkdirTemp creates a temporary directory and returns the followed path of the new directory
// Useful for environments where TMP_DIR is a symlink and you need the true path
func MkdirTemp(dir, pattern string) (string, error) {
	dir, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(dir)
}
