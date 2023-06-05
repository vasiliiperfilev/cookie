package data

import (
	"reflect"
	"sort"
)

type StubConversationModel struct {
	conversations []*Conversation
	idCount       int64
}

func NewStubConversationModel(conversations []*Conversation) *StubConversationModel {
	return &StubConversationModel{conversations: conversations}
}

func (s *StubConversationModel) Insert(conversation *Conversation) error {
	for _, existingConversation := range s.conversations {
		sort.Slice(existingConversation.UserIds, func(i, j int) bool {
			return existingConversation.UserIds[i] >= existingConversation.UserIds[j]
		})
		sort.Slice(conversation.UserIds, func(i, j int) bool {
			return conversation.UserIds[i] >= conversation.UserIds[j]
		})
		if reflect.DeepEqual(existingConversation.UserIds, conversation.UserIds) {
			return ErrDuplicateConversation
		}
	}
	s.idCount++
	conversation.Id = s.idCount
	conversation.LastMessageId = -1
	s.conversations = append(s.conversations, conversation)
	return nil
}

func (s *StubConversationModel) GetAllByUserId(userId int64) ([]*Conversation, error) {
	result := []*Conversation{}
	for _, conversation := range s.conversations {
		for _, id := range conversation.UserIds {
			if id == userId {
				result = append(result, conversation)
			}
		}
	}
	return result, nil
}
