package data

import (
	"sync"
)

type StubMessageModel struct {
	mu           sync.Mutex
	conversation *StubConversationModel
	// Maps conversation id to array of message ids
	conversationMessages map[int64][]int64
	messages             map[int64]Message
	IdCount              int64
}

func NewStubMessageModel(conversation *StubConversationModel, messages []Message) *StubMessageModel {
	model := &StubMessageModel{
		conversation:         conversation,
		conversationMessages: map[int64][]int64{},
		messages:             map[int64]Message{},
	}
	for _, msg := range messages {
		model.Insert(&msg)
	}
	return model
}

func (s *StubMessageModel) Insert(msg *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	conversation, err := s.conversation.GetById(msg.ConversationId)
	if err != nil {
		return ErrRecordNotFound
	}
	s.insertMessage(msg)
	s.updateReadMessage(conversation, msg)
	s.conversation.Update(conversation)
	return nil
}

func (s *StubMessageModel) updateReadMessage(conversation Conversation, msg *Message) {
	for i, readMessage := range conversation.LastReadMessages {
		if readMessage.UserId == msg.SenderId {
			conversation.LastReadMessages[i] = LastReadMessage{
				UserId:  readMessage.UserId,
				Message: *msg,
			}
		}
	}
}

func (s *StubMessageModel) insertMessage(msg *Message) {
	s.IdCount++
	msg.Id = s.IdCount
	s.messages[msg.Id] = *msg
	if _, ok := s.conversationMessages[msg.ConversationId]; !ok {
		s.conversationMessages[msg.ConversationId] = []int64{}
	}
	s.conversationMessages[msg.ConversationId] = append(s.conversationMessages[msg.ConversationId], msg.Id)
}

func (s *StubMessageModel) GetAllByConversationId(id int64) ([]Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := []Message{}
	conversationMessageIds := s.conversationMessages[id]
	for _, msgId := range conversationMessageIds {
		result = append(result, s.messages[msgId])
	}
	return result, nil
}

func (s *StubMessageModel) GetById(id int64) (Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if msg, ok := s.messages[id]; !ok {
		return Message{}, ErrRecordNotFound
	} else {
		return msg, nil
	}
}
