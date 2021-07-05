package util

import (
	"github.com/PandaTtttt/go-assembly/util/must"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func WalkPath(dirPath string, goDeep bool,
	fileF func(file string), dirF func(dir string)) {
	var wg sync.WaitGroup
	fs, err := ioutil.ReadDir(dirPath)
	must.Must(err)
	for _, fi := range fs {
		path := filepath.Join(dirPath, fi.Name())
		if fi.IsDir() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				dirF(path)
			}()
			if goDeep {
				WalkPath(path, goDeep, fileF, dirF)
			}
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fileF(path)
			}()
		}
	}
	wg.Wait()
}

// PrepareDir makes all dirs if not exists and changes perm to 777.
func PrepareDir(dirs ...string) {
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		must.Must(os.MkdirAll(dir, os.ModePerm))
		must.Must(os.Chmod(dir, os.ModePerm))
	}
}

// PrepareFileDir makes all dirs of given file if not exists and changes perm to 777.
func PrepareFileDir(files ...string) {
	for _, file := range files {
		if file == "" {
			continue
		}
		dir := filepath.Dir(file)
		must.Must(os.MkdirAll(dir, os.ModePerm))
		must.Must(os.Chmod(dir, os.ModePerm))
	}
}
