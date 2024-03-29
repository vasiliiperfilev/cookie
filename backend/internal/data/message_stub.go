package data

import (
	"sync"
)

type StubMessageModel struct {
	mu            sync.Mutex
	conversations storage
}

type storage map[int64]struct {
	Users    []User
	Messages []Message
	IdCount  int64
}

func NewStubMessageModel(conversations []Conversation, messages []Message) *StubMessageModel {
	msgStorage := storage{}
	for _, conversation := range conversations {
		msgStorage[conversation.Id] = struct {
			Users    []User
			Messages []Message
			IdCount  int64
		}{
			Users:    conversation.Users,
			Messages: []Message{{Id: 0, SenderId: 0, ConversationId: conversation.Id, PrevMessageId: 0}},
			IdCount:  int64(0),
		}
	}
	for _, message := range messages {
		if entry, ok := msgStorage[message.ConversationId]; ok {
			entry.IdCount++
			message.Id = entry.IdCount
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

func (s *StubMessageModel) GetById(id int64) (Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, conversation := range s.conversations {
		for _, msg := range conversation.Messages {
			if msg.Id == id {
				return msg, nil
			}
		}
	}
	return Message{}, ErrRecordNotFound
}

func (s *StubMessageModel) DeleteById(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, conversation := range s.conversations {
		for i, msg := range conversation.Messages {
			if msg.Id == id {
				conversation.Messages = append(conversation.Messages[:i], conversation.Messages[i+1:]...)
				s.conversations[key] = conversation
				return nil
			}
		}
	}
	return ErrRecordNotFound
}

func (s *StubMessageModel) Update(msg Message) error {
	return nil
}
