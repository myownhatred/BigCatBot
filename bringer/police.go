package bringer

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
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

func (p *PolicePostgres) UserDefaultCheck(UserID int64, username, firstname, lastname, command string) (ID int64, err error) {
	// Check if user already exists
	var count int
	err = p.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", UserID).Scan(&count)
	if err != nil {
		p.logger.Warn(err.Error())
		return 0, err
	}

	if count > 0 {
		p.logger.Info("User exist ", UserID)
		return UserID, nil
	} else {
		// Insert the new user
		stmt, err := p.db.Prepare("INSERT INTO users (id, first_name, last_name, username) VALUES ($1, $2, $3, $4)")
		if err != nil {
			p.logger.Warn(err.Error())
			return 0, err
		}
		defer stmt.Close()

		_, err = stmt.Exec(UserID, firstname, lastname, username)
		if err != nil {
			p.logger.Warn(err.Error())
			return 0, err
		}

		p.logger.Info("User added",
			slog.Int64("userID:", UserID),
			slog.String("username:", username),
			slog.String("firstname:", firstname),
			slog.String("lastname:", lastname))

		return 0, nil
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

func (p *PolicePostgres) Achieves(GRID int) (IDs []int, GRIDs []int, Names []string, Ranks []int, Descrs []string, err error) {
	// if GRID id == 0 we need to get all achievements
	if GRID == 0 {
		rows, err := p.db.Query("SELECT id, groupid, name, rank, description FROM achieves")
		if err != nil {
			p.logger.Warn(err.Error())
			return nil, nil, nil, nil, nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var ID int
			var grID int
			var name string
			var rank int
			var description string
			err := rows.Scan(&ID, &grID, &name, &rank, &description)
			if err != nil {
				p.logger.Warn(err.Error())
				return nil, nil, nil, nil, nil, err
			}
			IDs = append(IDs, ID)
			GRIDs = append(GRIDs, grID)
			Names = append(Names, name)
			Ranks = append(Ranks, rank)
			Descrs = append(Descrs, description)
		}
		return IDs, GRIDs, Names, Ranks, Descrs, nil
	} else {
		rows, err := p.db.Query("SELECT id, groupid, name, rank, description FROM achieves where groupid = ?", GRID)
		if err != nil {
			p.logger.Warn(err.Error())
			return nil, nil, nil, nil, nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var ID int
			var grID int
			var name string
			var rank int
			var description string
			err := rows.Scan(&ID, &grID, &name, &rank, &description)
			if err != nil {
				p.logger.Warn(err.Error())
				return nil, nil, nil, nil, nil, err
			}
			IDs = append(IDs, ID)
			GRIDs = append(GRIDs, grID)
			Names = append(Names, name)
			Ranks = append(Ranks, rank)
			Descrs = append(Descrs, description)
		}
		return IDs, GRIDs, Names, Ranks, Descrs, nil
	}
}

func (p *PolicePostgres) AchGroups() (IDs []int, GroupNames []string, err error) {
	rows, err := p.db.Query("SELECT id, groupname FROM achievegroups")
	if err != nil {
		p.logger.Warn(err.Error())
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ID int
		var groupname string
		err := rows.Scan(&ID, &groupname)
		if err != nil {
			p.logger.Warn(err.Error())
			return nil, nil, err
		}
		IDs = append(IDs, ID)
		GroupNames = append(GroupNames, groupname)
	}
	return IDs, GroupNames, nil
}

func (p *PolicePostgres) UserAchs(UserID int64) (IDs []int, UserIDs []int64, AchIDs []int, Dates []time.Time, Chats []string, ChatIDs []int64, err error) {
	//SELECT achlist.id, achlist.uid, achs.name AS ach_name, achgroups.name AS group_name
	//FROM achlist
	//JOIN achs ON achlist.aid = achs.id
	//JOIN achgroups ON achs.grid = achgroups.id
	//WHERE achlist.uid = desired_uid;
	query := fmt.Sprintf("SELECT id, uid, aid, date, chat, chatid FROM achlist where uid = %d", UserID)
	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Warn(err.Error())
		return nil, nil, nil, nil, nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ID int
		var UserID int64
		var AchID int
		var Date time.Time
		var Chat string
		var ChatID int64
		err := rows.Scan(&ID, &UserID, &AchID, &Date, &Chat, &ChatID)
		if err != nil {
			p.logger.Warn(err.Error())
			return nil, nil, nil, nil, nil, nil, err
		}
		IDs = append(IDs, ID)
		UserIDs = append(UserIDs, UserID)
		AchIDs = append(AchIDs, AchID)
		Dates = append(Dates, Date)
		Chats = append(Chats, Chat)
		ChatIDs = append(ChatIDs, ChatID)
	}
	return IDs, UserIDs, AchIDs, Dates, Chats, ChatIDs, nil
}

func (p *PolicePostgres) UserAchAdd(UserID int64, AID int, chat string, chatID int64) (UAID int, err error) {
	// Check if user already have this achive
	var count int
	err = p.db.QueryRow("SELECT id FROM achlist WHERE uid = $1 and aid = $2", UserID, AID).Scan(&count)
	if err != nil {
		p.logger.Warn(err.Error())
	}

	if count > 0 {
		p.logger.Info("User already have some achive ", UserID)
		return count, nil
	} else {
		// Insert the new achive
		stmt, err := p.db.Prepare("INSERT INTO achlist (uid, aid, date, chat, chatid) VALUES ($1, $2, NOW(), $3, $4) RETURNING id")
		if err != nil {
			p.logger.Warn(err.Error())
			return 0, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(UserID, AID, chat, chatID).Scan(&UAID)
		if err != nil {
			p.logger.Warn(err.Error())
			return 0, err
		}

		p.logger.Info("User's achive added",
			slog.Int("UAID:", UAID),
			slog.Int64("userID:", UserID),
			slog.Int("AchID:", AID),
			slog.String("chat:", chat),
			slog.Int64("chatID:", chatID))

		return UAID, nil
	}
}
