package data

type StubItemModel struct {
	items   map[int64]Item
	idCount int64
}

func NewStubItemModel(items []Item) *StubItemModel {
	itemMap := map[int64]Item{}
	for _, item := range items {
		itemMap[item.Id] = item
	}
	return &StubItemModel{items: itemMap}
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
