package data

import (
	"sync"
)

type StubMessageModel struct {
	mu            sync.Mutex
	conversations storage
}

type storage map[int64]struct {
	UserIds  []int64
	Messages []Message
	IdCount  int64
}

func NewStubMessageModel(conversations []Conversation, messages []Message) *StubMessageModel {
	msgStorage := storage{}
	for _, conversation := range conversations {
		msgStorage[conversation.Id] = struct {
			UserIds  []int64
			Messages []Message
			IdCount  int64
		}{
			UserIds:  conversation.UserIds,
			Messages: []Message{},
			IdCount:  int64(0),
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

func (s *StubMessageModel) Insert(msg *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.conversations[msg.ConversationId]; !ok {
		return ErrRecordNotFound
	} else {
		entry.IdCount++
		msg.Id = entry.IdCount
		entry.Messages = append(entry.Messages, *msg)
		s.conversations[msg.ConversationId] = entry
		return nil
	}
}

func (s *StubMessageModel) GetAllByConversationId(id int64) ([]Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := []Message{}
	for conversationId, conversation := range s.conversations {
		if conversationId == id {
			result = append(result, conversation.Messages...)
		}
	}
	return result, nil
}
