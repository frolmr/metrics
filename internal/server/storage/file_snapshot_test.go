package storage

import (
	"bufio"
	"errors"
	"os"
	"testing"

	"github.com/frolmr/metrics/internal/server/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSaveData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSnap := mocks.NewMockSnapshooter(ctrl)

	tmpFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	mockSnap.EXPECT().
		SaveToSnapshot(gomock.Any()).
		DoAndReturn(func(writer *bufio.Writer) error {
			_, writerErr := writer.WriteString("test data")
			return writerErr
		}).
		Times(1)

	fileSnapshot := NewFileSnapshot(mockSnap, tmpFile.Name())

	err = fileSnapshot.SaveData()
	assert.NoError(t, err)

	fileContent, err := os.ReadFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "test data", string(fileContent))
}

func TestSaveData_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSnap := mocks.NewMockSnapshooter(ctrl)

	tmpFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	mockSnap.EXPECT().
		SaveToSnapshot(gomock.Any()).
		Return(errors.New("save error")).
		Times(1)

	fileSnapshot := NewFileSnapshot(mockSnap, tmpFile.Name())

	err = fileSnapshot.SaveData()
	assert.Error(t, err)
	assert.Equal(t, "save error", err.Error())
}

func TestRestoreData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSnap := mocks.NewMockSnapshooter(ctrl)

	tmpFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	testData := "test data\n"
	err = os.WriteFile(tmpFile.Name(), []byte(testData), snapshotFilePermissions)
	assert.NoError(t, err)

	mockSnap.EXPECT().
		RestoreFromSnapshot(gomock.Any()).
		DoAndReturn(func(reader *bufio.Reader) error {
			data, readerErr := reader.ReadString('\n')
			assert.NoError(t, readerErr)
			assert.Equal(t, testData, data)
			return nil
		}).
		Times(1)

	fileSnapshot := NewFileSnapshot(mockSnap, tmpFile.Name())

	err = fileSnapshot.RestoreData()
	assert.NoError(t, err)
}

func TestRestoreData_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSnap := mocks.NewMockSnapshooter(ctrl)

	tmpFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	err = os.WriteFile(tmpFile.Name(), []byte("test data"), snapshotFilePermissions)
	assert.NoError(t, err)

	mockSnap.EXPECT().
		RestoreFromSnapshot(gomock.Any()).
		Return(errors.New("restore error")).
		Times(1)

	fileSnapshot := NewFileSnapshot(mockSnap, tmpFile.Name())

	err = fileSnapshot.RestoreData()
	assert.Error(t, err)
	assert.Equal(t, "restore error", err.Error())
}

func TestRestoreData_FileNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSnap := mocks.NewMockSnapshooter(ctrl)

	fileSnapshot := NewFileSnapshot(mockSnap, "nonexistentfile")

	err := fileSnapshot.RestoreData()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}
