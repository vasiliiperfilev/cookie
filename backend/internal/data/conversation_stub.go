package data

type StubConversationModel struct {
	conversations []Conversation
	idCount       int64
}

func NewStubConversationModel(conversations []Conversation) *StubConversationModel {
	model := &StubConversationModel{conversations: []Conversation{}}
	for _, c := range conversations {
		dto := PostConversationDto{UserIds: c.UserIds}
		model.Insert(dto)
	}
	return model
}

func (s *StubConversationModel) Insert(dto PostConversationDto) (Conversation, error) {
	for _, existingConversation := range s.conversations {
		if EqualArrays(existingConversation.UserIds, dto.UserIds) {
			return Conversation{}, ErrDuplicateConversation
		}
	}
	lastReadMessages := []LastReadMessage{}
	for _, userId := range dto.UserIds {
		lastReadMessages = append(lastReadMessages, LastReadMessage{
			UserId: userId,
			Message: Message{
				Id: 0,
			},
		})
	}
	s.idCount++
	conversation := Conversation{
		Id:               s.idCount,
		UserIds:          dto.UserIds,
		LastMessageId:    0,
		Version:          1,
		LastReadMessages: lastReadMessages,
	}
	s.conversations = append(s.conversations, conversation)
	return conversation, nil
}

func (s *StubConversationModel) GetAllByUserId(userId int64) ([]Conversation, error) {
	result := []Conversation{}
	for _, conversation := range s.conversations {
		for _, id := range conversation.UserIds {
			if id == userId {
				result = append(result, conversation)
			}
		}
	}
	return result, nil
}

func (s *StubConversationModel) GetById(id int64) (Conversation, error) {
	for _, conversation := range s.conversations {
		if conversation.Id == id {
			return conversation, nil
		}
	}
	return Conversation{}, ErrRecordNotFound
}

func (s *StubConversationModel) Update(conversation Conversation) error {
	for i, c := range s.conversations {
		if conversation.Id == c.Id {
			s.conversations[i] = conversation
			return nil
		}
	}
	return ErrRecordNotFound
}
