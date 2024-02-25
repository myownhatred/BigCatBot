package bringer

import (
	"Guenhwyvar/entities"
	"database/sql"
)

type FreeMawPostgres struct {
	db *sql.DB
}

func NewFreeMawPostgres(db *sql.DB) *FreeMawPostgres {
	return &FreeMawPostgres{
		db: db,
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
