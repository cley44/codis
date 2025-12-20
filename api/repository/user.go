package repository

import (
	"codis/models"

	"github.com/samber/do/v2"
)

type UserRepository struct {
	postgresDatabaseService *PostgresDatabaseService
}

func NewUserRepository(injector do.Injector) (*UserRepository, error) {
	u := UserRepository{
		postgresDatabaseService: do.MustInvoke[*PostgresDatabaseService](injector),
	}

	return &u, nil
}

func (u UserRepository) CreateOrUpdate(username string, displayUsername *string, discordID string, discordAvatar *string, discordSession *models.DiscordSession, email string) (user models.User, err error) {

	q := `INSERT INTO public.user
			(username, display_username, discord_id, discord_avatar, discord_session, email)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT ON CONSTRAINT user_discord_id_key
		DO UPDATE SET
			display_username = excluded.display_username,
			discord_avatar = excluded.discord_avatar,
			discord_session = excluded.discord_session,
			email = excluded.email,
			username = excluded.username
		RETURNING *;`

	err = u.postgresDatabaseService.Get(&user, q, username, displayUsername, discordID, discordAvatar, discordSession, email)

	return user, err
}

func (u UserRepository) GetByID(ID string) (user models.User, err error) {

	q := `SELECT * from public.user WHERE id = $1`

	err = u.postgresDatabaseService.Get(&user, q, ID)
	return
}

func (u UserRepository) UpdateSession(userID string, session models.DiscordSession) (err error) {

	q := `UPDATE public.user SET discord_session = $2 WHERE id = $1;`

	err = u.postgresDatabaseService.Exec(q, userID, session)
	return
}
