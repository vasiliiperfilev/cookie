package data

type StubConversationModel struct {
	conversations []Conversation
	idCount       int64
	userModel     UserModel
}

func NewStubConversationModel(conversations []Conversation, userModel UserModel) *StubConversationModel {
	return &StubConversationModel{conversations: conversations, userModel: userModel}
}

func (s *StubConversationModel) Insert(dto PostConversationDto) (Conversation, error) {
	for _, existingConversation := range s.conversations {
		userIds := Map(existingConversation.Users, func(u User) int64 { return u.Id })
		if EqualArraysContent(userIds, dto.UserIds) {
			return Conversation{}, ErrDuplicateConversation
		}
	}
	s.idCount++
	conversation := Conversation{
		Id: s.idCount,
		Users: Map(dto.UserIds, func(id int64) User {
			u, _ := s.userModel.GetById(id)
			return u
		}),
		LastMessage: Message{Id: 0},
		Version:     1,
	}
	s.conversations = append(s.conversations, conversation)
	return conversation, nil
}

func (s *StubConversationModel) GetAllByUserId(userId int64) ([]Conversation, error) {
	result := []Conversation{}
	for _, conversation := range s.conversations {
		for _, u := range conversation.Users {
			if u.Id == userId {
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
