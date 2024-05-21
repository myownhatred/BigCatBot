package bringer

import (
	"database/sql"
	"log/slog"
)

type PolicePostgres struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPolicePostgres(db *sql.DB, logger *slog.Logger) *PolicePostgres {
	return &PolicePostgres{
		db:     db,
		logger: logger,
	}
}

func (p *PolicePostgres) UserDefaultCheck(UserID int64, username, firstname, lastname, command string) (err error) {
	// Check if user already exists
	var count int
	err = p.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", UserID).Scan(&count)
	if err != nil {
		p.logger.Warn(err.Error())
	}

	if count > 0 {
		p.logger.Info("User exist ", UserID)
	} else {
		// Insert the new user
		stmt, err := p.db.Prepare("INSERT INTO users (id, first_name, last_name, username) VALUES ($1, $2, $3, $4)")
		if err != nil {
			p.logger.Warn(err.Error())
		}
		defer stmt.Close()

		_, err = stmt.Exec(UserID, firstname, lastname, username)
		if err != nil {
			p.logger.Warn(err.Error())
			return err
		}

		p.logger.Info("User added",
			slog.Int64("userID:", UserID),
			slog.String("username:", username),
			slog.String("firstname:", firstname),
			slog.String("lastname:", lastname))
	}

	return nil
}
