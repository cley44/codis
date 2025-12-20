package discord

import (
	"codis/models"

	"github.com/samber/oops"
)

func (svc *DiscordService) VerifySession(userID string, session models.DiscordSession) (models.DiscordSession, error) {
	newSession, err := svc.oauthClient.VerifySession(session.Session)
	if err != nil {
		return models.DiscordSession{}, oops.Wrap(err)
	}

	// @TODO might need to optimize this and update only if session change
	s := models.DiscordSession{Session: newSession}
	err = svc.userRepository.UpdateSession(userID, s)
	if err != nil {
		return models.DiscordSession{}, oops.Wrap(err)
	}
	return s, nil
}
