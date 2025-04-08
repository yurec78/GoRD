package documentstore2

type Store struct {
	collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	if _, exists := s.collections[name]; exists {
		return false, nil
	}

	collection := &Collection{
		config:    cfg,
		documents: make(map[string]Document),
	}
	s.collections[name] = collection
	return true, collection
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	col, ok := s.collections[name]
	return col, ok
}

func (s *Store) DeleteCollection(name string) bool {
	_, ok := s.collections[name]
	if !ok {
		return false
	}
	delete(s.collections, name)
	return true
}
