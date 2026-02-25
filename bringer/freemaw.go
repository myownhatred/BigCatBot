package bringer

import (
	"Guenhwyvar/entities"
	"database/sql"
	"fmt"
	"log/slog"
)

type FreeMawPostgres struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewFreeMawPostgres(db *sql.DB, logger *slog.Logger) *FreeMawPostgres {
	return &FreeMawPostgres{
		db:     db,
		logger: logger,
	}
}

func (m *FreeMawPostgres) PutFreeMawToDB(maw entities.FreeMaw) (err error) {
	_, err = m.db.Exec("INSERT into free_maw (typ, description, link) values($1,$2,$3) RETURNING id", maw.Type, maw.Description, maw.Link)
	if err != nil {
		return err
	}
	return nil
}

// create type checking function to provide proper error in
// case wrong or no type provided

func (m *FreeMawPostgres) GetRandomMawFromDB(typ string) (maw entities.FreeMaw, err error) {
	row := m.db.QueryRow("SELECT id, typ, description, link FROM free_maw where typ=$1 ORDER BY random() LIMIT 1", typ)
	err = row.Scan(&maw.ID, &maw.Type, &maw.Description, &maw.Link)
	if err != nil {
		return maw, err
	}
	return maw, nil
}

func (m *FreeMawPostgres) FreeMawDBReport() (report string, err error) {
	return "not yet enplementash", nil
}

// add new quiz type/name

func (m *FreeMawPostgres) FreeMawVectorTypeAdd(qtype entities.VectorType) (err error) {
	_, err = m.db.Exec("INSERT into question_types (type_name, protected) values($1,$2) RETURNING id", qtype.Name, qtype.Protected)
	if err != nil {
		m.logger.Warn("bringer freemaw", "error inserting new vector type:", err.Error())
		return err
	}
	return nil
}

// get type by ID
func (m *FreeMawPostgres) FreeMawVectorTypeByID(ID int) (qtype entities.VectorType, err error) {
	row := m.db.QueryRow("SELECT id, type_name, protected FROM question_types where id=$1", ID)
	err = row.Scan(&qtype.ID, &qtype.Name, &qtype.Protected)
	if err != nil {
		m.logger.Warn("bringer freemaw", "error getting vector type by ID:", err.Error())
		return qtype, err
	}
	return qtype, nil
}

// list of possible types/quiz names
func (m *FreeMawPostgres) FreeMawVectorTypes() (report []entities.VectorType, err error) {
	rows, err := m.db.Query("SELECT * FROM question_types")
	if err != nil {
		m.logger.Warn(err.Error())
		return report, err
	}
	defer rows.Close()

	for rows.Next() {
		var tipe entities.VectorType
		err = rows.Scan(&tipe.ID, &tipe.Name, &tipe.Protected)
		if err != nil {
			m.logger.Warn(err.Error())
			return report, err
		}
		report = append(report, tipe)
	}
	return report, nil
}

// add question
func (m *FreeMawPostgres) FreeMawVectorAdd(vec entities.FreeVector) (err error) {
	mainID := 0

	err = m.db.QueryRow("INSERT into question (typeid, pic_link, question_string, userid) values($1,$2,$3,$4) RETURNING id", vec.TypeID, vec.PicLink, vec.Question, vec.UserID).Scan(&mainID)
	if err != nil {
		m.logger.Warn("bringer freemaw", "error inserting new vector :", err.Error())
		return err
	}
	// vec comes with answers ID's set to 6666 and own IDs are local indexes
	// rearranging it here
	for _, answer := range vec.Answers {
		_, err = m.db.Exec("INSERT into answer (questionid, answer_text) values($1,$2) RETURNING id", mainID, answer.Answer)
		if err != nil {
			m.logger.Warn("bringer freemaw", "error inserting answer :", err.Error(),
				slog.Int("Question ID", mainID))
			return err
		}
	}
	return nil
}

// get random question for s
func (m *FreeMawPostgres) FreeMawVectorGetRandomByType(typ int) (vec entities.FreeVector, err error) {
	row := m.db.QueryRow("SELECT id, typeid, pic_link, question_string, userid FROM question where typeid=$1 ORDER BY random() LIMIT 1", typ)
	err = row.Scan(&vec.ID, &vec.TypeID, &vec.PicLink, &vec.Question, &vec.UserID)
	if err != nil {
		m.logger.Warn("bringer freemaw", "error getting random vector:", err.Error(),
			slog.Int("Vector type ID", typ))
		return vec, err
	}
	m.logger.Info("bringer freemaw - vector entrails",
		slog.Int("Vector ID", vec.ID),
		slog.Int("Vector type ID", vec.TypeID),
		slog.String("Vector question", vec.Question),
		slog.String("Vector pic link", vec.PicLink),
		slog.Int64("Vector uid", vec.UserID),
	)
	rows, err := m.db.Query("SELECT id, questionid, answer_text FROM answer where questionid=$1", vec.ID)
	if err != nil {
		m.logger.Warn("bringer freemas error getting answers for question",
			slog.Int("question ID:", vec.ID))
		return vec, err
	}
	defer rows.Close()
	for rows.Next() {
		var ans entities.VectorAnswer
		err = rows.Scan(&ans.ID, &ans.QuestionID, &ans.Answer)
		if err != nil {
			m.logger.Warn("bringer freemas error reading rows of answeres for question",
				slog.Int("question ID:", vec.ID),
				slog.Int("answer ID:", ans.ID))
			return vec, err
		}
		vec.Answers = append(vec.Answers, ans)
	}

	return vec, nil
}

func (m *FreeMawPostgres) FreeMawVectorUpsertScore(uid int64, vectorType int, score int) error {
	query := `
		INSERT INTO vector_scores (uid, vector_type, score)
		VALUES ($1, $2, $3)
		ON CONFLICT (uid, vector_type)
		DO UPDATE SET score = vector_scores.score + $3
	`

	_, err := m.db.Exec(query, uid, vectorType, score)
	if err != nil {
		m.logger.Error("bringer freemas error - upserting score",
			slog.String("query", query),
			slog.Int64("UID", uid),
			slog.Int("vector type", vectorType),
			slog.Int("score", score))
		return fmt.Errorf("error on upsert score: %w", err)
	}
	return nil
}

func (m *FreeMawPostgres) FreeMawVectorGetTopScores(limit int) (scores []entities.VectorScore, err error) {
	query := `
		SELECT uid, 0 as vec_type, SUM(score) as total_score
		FROM vector_scores
		GROUP BY uid
		ORDER BY total_score DESC
		LIMIT $1
	`

	rows, err := m.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting all scores from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var score entities.VectorScore
		if err = rows.Scan(&score.UID, &score.VectorType, &score.Score); err != nil {
			return nil, fmt.Errorf("error getting single scores row from DB: %w", err)
		}
		scores = append(scores, score)
	}
	return scores, nil
}
