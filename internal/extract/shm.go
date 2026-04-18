package extract

import (
	"fmt"
	"os"
	"sync"
	"syscall"
)

type ShmManager struct {
	Path   string
	Size   int64
	file   *os.File
	data   []byte
	offset int64
	mu     sync.Mutex
}

func NewShmManager(path string, size int64) (*ShmManager, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	if err := f.Truncate(size); err != nil {
		_ = f.Close()
		return nil, err
	}
	b, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	return &ShmManager{
		Path: path, Size: size, file: f, data: b, offset: 0,
	}, nil
}

func (s *ShmManager) Write(payload []byte) (int64, int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := int64(len(payload))
	if n > s.Size {
		return 0, 0, fmt.Errorf("payload too large: %d > shm size %d", n, s.Size)
	}

	if s.offset+n > s.Size {
		s.offset = 0
	}
	start := s.offset
	copy(s.data[start:start+n], payload)
	s.offset += n
	return start, n, nil
}

func (s *ShmManager) Close() error {
	if s.data != nil {
		_ = syscall.Munmap(s.data)
	}
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}