package application

import (
	"task1/items_manager/pkg/models"
	"time"

	"github.com/google/uuid"
)

type sessionManager struct {
	keyToUser map[string]session
	userToKey map[string]session
}

type session struct {
	Username string
	Token    string
	Expiry   time.Time
}

func (s *sessionManager) isTokenValid(api_key string) (session, bool) {
	val, ok := s.keyToUser[api_key]

	if !ok {
		// fmt.Println("api key does not exist")
		panic(models.ErrInvalidApiKey)
	}
	ok = isSessionExpired(val)

	if ok {
		return session{}, false
	}

	return val, true
}

func isSessionExpired(val session) bool {
	return time.Now().After(val.Expiry)
}

func (s *sessionManager) currentActiveSession(username string) bool {
	val, ok := s.userToKey[username]

	if !ok {
		return false
	}
	return isSessionExpired(val)

}

func (s *sessionManager) removeSession(val session) {
	delete(s.keyToUser, val.Token)
	delete(s.userToKey, val.Username)
}

func (s *sessionManager) getSessionToken(username string) string {
	for {
		newToken := uuid.New().String()
		_, ok := s.keyToUser[newToken]
		if !ok {
			newSession := session{Username: username,
				Token:  newToken,
				Expiry: time.Now().Add(2 * time.Hour)}
			s.keyToUser[newToken] = newSession
			s.userToKey[username] = newSession
			return newToken
		}
	}
}
