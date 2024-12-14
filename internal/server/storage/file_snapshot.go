package storage

import (
	"bufio"
	"os"
)

const (
	snapshotFilePermissions = 0600
)

type FileSnapshooter interface {
	SaveData() error
	RestoreData() error
}

type FileSnapshot struct {
	snap     Snapshooter
	fileName string
}

func NewFileSnapshot(snap Snapshooter, fileName string) *FileSnapshot {
	return &FileSnapshot{
		snap:     snap,
		fileName: fileName,
	}
}

func (fs *FileSnapshot) SaveData() error {
	fileInfo, err := os.Stat(fs.fileName)

	var file *os.File
	if err == nil && fileInfo.Mode().IsRegular() {
		file, err = os.OpenFile(fs.fileName, os.O_RDWR, snapshotFilePermissions)
		if err != nil {
			return err
		}
	} else {
		file, err = os.Create(fs.fileName)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if err := fs.snap.SaveToSnapshot(writer); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func (fs *FileSnapshot) RestoreData() error {
	if _, err := os.Stat(fs.fileName); err != nil {
		return err
	}

	file, err := os.OpenFile(fs.fileName, os.O_RDONLY, snapshotFilePermissions)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	if err := fs.snap.RestoreFromSnapshot(reader); err != nil {
		return err
	}
	return nil
}
