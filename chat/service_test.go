package chat_test

import (
	"practice-run/chat"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	svc *chat.Service
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSubTest() {
	s.svc = chat.NewService()
}

type MockMember struct {
	username         string
	lastNotification chat.Event
}

func (m *MockMember) Username() string {
	return m.username
}

func (m *MockMember) Notify(event chat.Event) {
	m.lastNotification = event
}
