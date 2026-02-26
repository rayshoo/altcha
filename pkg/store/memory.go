package store

import "sync"

type MemoryStore struct {
	mu         sync.Mutex
	records    []string
	maxRecords int
}

func NewMemoryStore(maxRecords int) *MemoryStore {
	return &MemoryStore{maxRecords: maxRecords}
}

func (s *MemoryStore) Exists(token string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, r := range s.records {
		if r == token {
			return true, nil
		}
	}
	return false, nil
}

func (s *MemoryStore) Add(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = append(s.records, token)
	if len(s.records) > s.maxRecords {
		s.records = s.records[1:]
	}
	return nil
}

func (s *MemoryStore) Ping() error  { return nil }
func (s *MemoryStore) Close() error { return nil }
