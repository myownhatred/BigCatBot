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
		return nil
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

		return nil
	}
}

func (p *PolicePostgres) MetatronChatAdd(ChatID int64, ChatName string) (err error) {
	// Check if chat already exists
	var count int
	err = p.db.QueryRow("SELECT COUNT(*) FROM metatron WHERE chat_id = $1", ChatID).Scan(&count)
	if err != nil {
		p.logger.Warn(err.Error())
	}

	if count > 0 {
		p.logger.Info("Chat exist ", ChatID)
	} else {
		// Insert the new chat to table
		stmt, err := p.db.Prepare("INSERT INTO metatron (chat_id, chat_name) VALUES ($1, $2)")
		if err != nil {
			p.logger.Warn(err.Error())
		}
		defer stmt.Close()

		_, err = stmt.Exec(ChatID, ChatName)
		if err != nil {
			p.logger.Warn(err.Error())
			return err
		}

		p.logger.Info("Chat added",
			slog.Int64("ChatID:", ChatID),
			slog.String("chatname:", ChatName))
	}

	return nil
}

func (p *PolicePostgres) MetatronChatList() (IDs []int64, ChatIDs []int64, Names []string, err error) {
	rows, err := p.db.Query("SELECT id, chat_id, chat_name FROM metatron")
	if err != nil {
		p.logger.Warn(err.Error())
		return nil, nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ID int64
		var chatID int64
		var name string
		err := rows.Scan(&ID, &chatID, &name)
		if err != nil {
			p.logger.Warn(err.Error())
			return nil, nil, nil, err
		}
		IDs = append(IDs, ID)
		ChatIDs = append(ChatIDs, chatID)
		Names = append(Names, name)
	}
	return IDs, ChatIDs, Names, nil
}
