package data

import (
	"reflect"
	"sort"
)

type StubConversationModel struct {
	conversations []Conversation
	idCount       int64
}

func NewStubConversationModel(conversations []Conversation) *StubConversationModel {
	return &StubConversationModel{conversations: conversations}
}

func (s *StubConversationModel) Insert(dto PostConversationDto) (Conversation, error) {
	for _, existingConversation := range s.conversations {
		sort.Slice(existingConversation.UserIds, func(i, j int) bool {
			return existingConversation.UserIds[i] >= existingConversation.UserIds[j]
		})
		sort.Slice(dto.UserIds, func(i, j int) bool {
			return dto.UserIds[i] >= dto.UserIds[j]
		})
		if reflect.DeepEqual(existingConversation.UserIds, dto.UserIds) {
			return Conversation{}, ErrDuplicateConversation
		}
	}
	s.idCount++
	conversation := Conversation{
		Id:            s.idCount,
		UserIds:       dto.UserIds,
		LastMessageId: 0,
		Version:       1,
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
