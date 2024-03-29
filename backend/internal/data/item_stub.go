package data

type StubItemModel struct {
	items   map[int64]Item
	idCount int64
}

func NewStubItemModel(items []Item) *StubItemModel {
	itemMap := map[int64]Item{}
	idCount := 0
	for _, item := range items {
		itemMap[item.Id] = item
		idCount++
	}
	return &StubItemModel{items: itemMap, idCount: int64(idCount)}
}

func (s *StubItemModel) Insert(item *Item) error {
	s.idCount++
	item.Id = s.idCount
	s.items[item.Id] = *item
	return nil
}

func (s *StubItemModel) GetById(id int64) (Item, error) {
	if item, ok := s.items[id]; ok {
		return item, nil
	}
	return Item{}, ErrRecordNotFound
}

func (s *StubItemModel) GetAllBySupplierId(id int64) ([]Item, error) {
	result := []Item{}
	for _, item := range s.items {
		if item.SupplierId == id {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return result, ErrRecordNotFound
	}
	return result, nil
}

func (s *StubItemModel) Update(item Item) (Item, error) {
	updated := false
	for i, existingItem := range s.items {
		if existingItem.Id == item.Id {
			s.items[i] = item
			updated = true
		}
	}
	if !updated {
		return Item{}, ErrRecordNotFound
	}
	return item, nil
}

func (s *StubItemModel) Delete(id int64) error {
	if _, ok := s.items[id]; ok {
		delete(s.items, id)
		return nil
	}
	return ErrRecordNotFound
}
