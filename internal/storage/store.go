package storage

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type File struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	AccessCount  int       `json:"access_count"`
	DownloadUrl  string    `json:"download_url"`

	StoragePath string `json:"-"` // Internal path, do not return to user.
}

var (
	fileStore       = make(map[string]File)
	fileStoreLock   = &sync.RWMutex{}
	baseStoragePath string
)

func Init(storagePath string) error {
	slog.Info("Initialising storage", "path", storagePath)

	if err := os.MkdirAll(storagePath, 0755); err != nil {
		slog.Error("Failed to create storage directory", "path", storagePath, "error", err)
		return err
	}

	baseStoragePath = storagePath
	slog.Info("Storage ready", "path", storagePath)
	return nil
}

func Get(id string) (File, bool) {
	fileStoreLock.RLock()
	defer fileStoreLock.RUnlock()

	file, ok := fileStore[id]
	return file, ok
}

func Save(file File) {
	fileStoreLock.Lock()
	defer fileStoreLock.Unlock()

	fileStore[file.ID] = file
}

func Delete(id string) {
	fileStoreLock.Lock()
	defer fileStoreLock.Unlock()

	delete(fileStore, id)
}

func List() []File {
	fileStoreLock.RLock()
	defer fileStoreLock.RUnlock()

	files := make([]File, 0)
	for _, file := range fileStore {
		files = append(files, file)
	}

	return files
}

func ListExpired() []File {
	fileStoreLock.Lock()
	defer fileStoreLock.Unlock()

	files := make([]File, 0)
	for _, file := range fileStore {
		if file.ExpiresAt.After(time.Now()) {
			files = append(files, file)
		}
	}

	return files
}

func Path(id string) string {
	if baseStoragePath == "" {
		slog.Error("Storage has not been initialised.")
		return ""
	}

	return filepath.Join(baseStoragePath, id)
}
