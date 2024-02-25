package bringer

import (
	"Guenhwyvar/entities"
	"database/sql"
	"fmt"
	"time"
)

type TimeWithOutPostgres struct {
	db *sql.DB
}

func NewTimeWithOutPostgres(db *sql.DB) *TimeWithOutPostgres {
	return &TimeWithOutPostgres{
		db: db,
	}
}

func (t *TimeWithOutPostgres) GetTimeWithOutList(chatID int64) (list []entities.TimeWithOut, err error) {
	rows, err := t.db.Query("SELECT * FROM time_with_out where chat_id=$1", chatID)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var event entities.TimeWithOut
		err = rows.Scan(&event.ID, &event.Name, &event.Time, &event.ChatId)
		if err != nil {
			fmt.Println(err.Error())
			return list, err
		}
		list = append(list, event)
	}
	return list, nil
}

func (t *TimeWithOutPostgres) GetTimeWithOutTimerByID(id int) (event entities.TimeWithOut, err error) {
	row := t.db.QueryRow("SELECT * FROM time_with_out where id=$1", id)
	err = row.Scan(&event.ID, &event.Name, &event.Time, &event.ChatId)
	return event, err
}

func (t *TimeWithOutPostgres) ResetTimer(id int) (err error) {
	row, err := t.db.Exec("UPDATE time_with_out SET time=$1 WHERE id=$2", time.Now(), id)
	if err != nil {
		fmt.Println(err.Error())
	}
	rows, err := row.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("affected rows = %d\n", rows)
	if err != nil {
		return err
	}
	return nil
}

func (t *TimeWithOutPostgres) AddNewTimer(name string, chatID int64) (err error) {
	_, err = t.db.Exec("INSERT into time_with_out (name, time, chat_id) values($1,$2,$3) RETURNING id", name, time.Now(), chatID)
	if err != nil {
		return err
	}
	return nil
}
