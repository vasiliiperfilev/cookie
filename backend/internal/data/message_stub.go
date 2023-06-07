package data

import (
	"sync"

	"golang.org/x/exp/slices"
)

type StubMessageModel struct {
	mu            sync.Mutex
	conversations storage
	idCount       int64
}

type storage map[int64]struct {
	UserIds  []int64
	Messages []Message
}

func NewStubMessageModel(conversations []Conversation, messages []Message) *StubMessageModel {
	msgStorage := storage{}
	for _, conversation := range conversations {
		msgStorage[conversation.Id] = struct {
			UserIds  []int64
			Messages []Message
		}{
			UserIds:  conversation.UserIds,
			Messages: []Message{},
		}
	}
	for _, message := range messages {
		if entry, ok := msgStorage[message.ConversationId]; ok {
			entry.Messages = append(entry.Messages, message)
			msgStorage[message.ConversationId] = entry
		}
	}

	return &StubMessageModel{conversations: msgStorage}
}

func (s *StubMessageModel) Insert(msg Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.conversations[msg.ConversationId]; !ok {
		return ErrRecordNotFound
	} else {
		s.idCount++
		msg.Id = s.idCount
		entry.Messages = append(entry.Messages, msg)
		s.conversations[msg.ConversationId] = entry
		return nil
	}
}

func (s *StubMessageModel) GetAllByUserId(id int64) ([]Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := []Message{}
	for _, conversation := range s.conversations {
		if slices.Contains(conversation.UserIds, id) {
			result = append(result, conversation.Messages...)
		}
	}
	return result, nil
}
