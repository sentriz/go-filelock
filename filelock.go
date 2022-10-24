package filelock

import (
	"fmt"
	"io"
	"sync"

	"golang.org/x/sys/unix"
)

type File struct {
	fd int
}

var _ sync.Locker = (*File)(nil)
var _ io.ReadWriteSeeker = (*File)(nil)
var _ io.Closer = (*File)(nil)

func New(path string) (*File, error) {
	fd, err := unix.Open(path, unix.O_CREAT|unix.O_RDWR, 0700)
	if err != nil {
		return nil, fmt.Errorf("open path: %w", err)
	}
	return &File{fd: fd}, nil
}

// Lock implements Lock() from sync.Locker
func (f *File) Lock() {
	if err := unix.Flock(f.fd, unix.LOCK_EX); err != nil {
		panic(fmt.Errorf("lock fd: %w", err))
	}
}

// Unlock implements Unlock() from sync.Locker
func (f *File) Unlock() {
	if err := unix.Flock(f.fd, unix.LOCK_UN); err != nil {
		panic(fmt.Errorf("unlock fd: %w", err))
	}
}

// Read implements Read() from io.Reader
func (f *File) Read(p []byte) (int, error) {
	return unix.Read(f.fd, p)
}

// Writer implements Write() from io.Writer
func (f *File) Write(p []byte) (int, error) {
	return unix.Write(f.fd, p)
}

// Writer implements Seek() from io.Seeker
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return unix.Seek(f.fd, offset, whence)
}

// Close implements Close() from io.Closer
func (f *File) Close() error {
	if err := unix.Flock(f.fd, unix.LOCK_UN); err != nil {
		return fmt.Errorf("unlock fd: %w", err)
	}
	if err := unix.Close(f.fd); err != nil {
		return fmt.Errorf("close fd: %w", err)
	}
	return nil
}
