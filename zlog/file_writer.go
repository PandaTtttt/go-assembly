package zlog

import (
	"github.com/PandaTtttt/go-assembly/errs"
	"github.com/PandaTtttt/go-assembly/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// FileWriter is a infrastructure of zlog, performs rotate operation if necessary,
// it's forked from github.com/natefinch/lumberjack
// but fix https://github.com/natefinch/lumberjack/issues/56, https://github.com/natefinch/lumberjack/issues/81.

const (
	backupNameFormat = "2006-01-02T15-04-05"
	backupSuffix     = ".bak"
)

// If MaxBackups and MaxLifetime are both 0, no old log files will be deleted.
type FileWriter struct {
	// File represents a log file where we write logs to.
	File string
	// BackupDir is the directory where backup log files will be retained,
	// defaults to (directory of file)/backup/.
	BackupDir string
	// MaxSize is the maximum size in bytes of the log file before it gets
	// rotated, if MaxSize equals zero, FileWriter does not perform rotation.
	MaxSize int64
	// MaxLifetime is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename. The default is not to remove old log
	// files based on lifetime.
	MaxLifetime int
	// MaxBackups is the maximum number of old log files to retain,
	// defaults to retain all old log files.
	MaxBackups int
	// UTCTime determines if the time used for formatting the timestamps in
	// backup files is the UTC time, defaults to use the computer's local time.
	UTCTime bool
	// Truncation determines whether the log files should be truncated instead of rotation
	// defaults to false.
	Truncation bool

	fd   *os.File
	size int64
	mu   sync.Mutex
}

func (w *FileWriter) callback(logger *Logger) error {
	if w.File == "" {
		return errs.InvalidConfig.New("missing FileWriter.File")
	}
	return nil
}

func (w *FileWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.fd == nil {
		if w.File == "" {
			return 0, errs.InvalidConfig.New("missing FileWriter.File")
		}

		util.PrepareFileDir(w.File)
		if w.MaxSize != 0 && !w.Truncation {
			util.PrepareDir(w.backupDir())
		}

		fd, err := os.OpenFile(w.File,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return 0, err
		}
		info, err := os.Stat(w.File)
		if err != nil {
			return 0, err
		}
		w.fd = fd
		w.size = info.Size()
	}

	writeLen := int64(len(b))
	if w.MaxSize != 0 && writeLen > w.MaxSize {
		return 0, errs.Internal.Newf("write length %d exceeds maximum file size %d", writeLen, w.MaxSize)
	}

	if w.MaxSize != 0 && w.size+writeLen > w.MaxSize {
		if w.Truncation {
			err := os.Truncate(w.File, 0)
			if err != nil {
				return 0, err
			}
			w.size = 0
		} else {
			err := w.rotate()
			if err != nil {
				return 0, err
			}
		}
	}

	n, err := w.fd.Write(b)
	w.size += int64(n)

	return n, err
}

func (w *FileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.close()
}

func (w *FileWriter) close() error {
	if w.fd == nil {
		return nil
	}
	err := w.fd.Close()
	w.fd = nil
	return err
}

func (w *FileWriter) rotate() error {
	err := w.close()
	if err != nil {
		return err
	}
	err = os.Rename(w.File, w.backupName())
	if err != nil {
		return err
	}
	fd, err := os.OpenFile(w.File,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	w.fd = fd
	w.size = 0

	return w.mill()
}

func (w *FileWriter) mill() error {
	if w.MaxBackups == 0 && w.MaxLifetime == 0 {
		return nil
	}
	remaining, err := w.oldLogFiles()
	if err != nil {
		return err
	}
	var remove logFiles

	if w.MaxBackups > 0 && w.MaxBackups < len(remaining) {
		remove = remaining[w.MaxBackups:]
		remaining = remaining[:w.MaxBackups]
	}

	if w.MaxLifetime > 0 {
		diff := 24 * time.Hour * time.Duration(w.MaxLifetime)
		cutoff := time.Now().Add(-1 * diff)
		for _, f := range remaining {
			if f.timestamp.Before(cutoff) {
				remove = append(remove, f)
			}
		}
	}

	for _, f := range remove {
		err := os.Remove(filepath.Join(w.backupDir(), f.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *FileWriter) backupName() string {
	timeStr := time.Now().Format(backupNameFormat)
	if w.UTCTime {
		timeStr = time.Now().UTC().Format(backupNameFormat)
	}
	backupFile := filepath.Join(w.backupDir(), timeStr) + backupSuffix
	util.PrepareFileDir(backupFile)
	return backupFile
}

func (w *FileWriter) backupDir() string {
	filename := util.NameWithoutExt(w.File)
	if w.BackupDir != "" {
		return filepath.Join(w.BackupDir, filename)
	}
	fileDir := filepath.Dir(w.File)
	return filepath.Join(fileDir, "backup", filename)
}

func (w *FileWriter) oldLogFiles() (logFiles, error) {
	fileInfos, err := ioutil.ReadDir(w.backupDir())
	if err != nil {
		return nil, errs.Internal.Newf("can't read backup log file directory: %s", err)
	}
	var res logFiles
	for _, fi := range fileInfos {
		if fi.IsDir() {
			continue
		}
		t, err := w.timeFromName(fi.Name())
		if err != nil {
			// ignore error
			continue
		}
		res = append(res, struct {
			timestamp time.Time
			os.FileInfo
		}{t, fi})
	}
	sort.Sort(res)

	return res, nil
}

func (w *FileWriter) timeFromName(name string) (time.Time, error) {
	if !strings.HasSuffix(name, backupSuffix) {
		return time.Time{}, errs.Internal.New("mismatched suffix")
	}
	ts := name[:len(name)-len(backupSuffix)]
	if w.UTCTime {
		return time.ParseInLocation(backupNameFormat, ts, time.UTC)
	}
	return time.ParseInLocation(backupNameFormat, ts, time.Local)
}

type logFiles []struct {
	timestamp time.Time
	os.FileInfo
}

func (b logFiles) Less(i, j int) bool {
	return b[i].timestamp.After(b[j].timestamp)
}

func (b logFiles) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b logFiles) Len() int {
	return len(b)
}
