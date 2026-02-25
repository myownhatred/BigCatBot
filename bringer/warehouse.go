package bringer

import (
	"Guenhwyvar/entities"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
)

type WarehouseMariaDB struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewFreeMawMariaDB(db *sql.DB, logger *slog.Logger) *WarehouseMariaDB {
	return &WarehouseMariaDB{
		db:     db,
		logger: logger,
	}
}

func (w *WarehouseMariaDB) GetRandomBoyan() (*entities.Boyan, error) {
	var totalCount int
	err := w.db.QueryRow("SELECT COUNT(*) FROM img_storage").Scan(&totalCount)
	if err != nil {
		w.logger.Error("Failed to count boyans", "error", err)
		return nil, fmt.Errorf("failed to count boyans: %w", err)
	}

	if totalCount == 0 {
		return nil, fmt.Errorf("no boyans found")
	}

	randomOffset := rand.Intn(totalCount)
	query := `
        SELECT id, link, name 
        FROM img_storage 
        LIMIT 1 OFFSET ?
    `

	var boyan entities.Boyan
	err = w.db.QueryRow(query, randomOffset).Scan(&boyan.ID, &boyan.Link, &boyan.Name)
	if err != nil {
		w.logger.Error("Failed to get random boyan", "error", err)
		return nil, fmt.Errorf("failed to get random boyan: %w", err)
	}

	// Получаем теги для этой картинки
	tags, err := w.GetTagsForPicture(boyan.ID)
	if err != nil {
		w.logger.Warn("Failed to get tags for boyan", "boyanID", boyan.ID, "error", err)
	}
	boyan.Tags = tags

	w.logger.Debug("Retrieved random boyan", "boyanID", boyan.ID, "name", boyan.Name)
	return &boyan, nil
}

func (w *WarehouseMariaDB) GetTagsForPicture(pictureID int) ([]string, error) {
	query := `
        SELECT bt.tag 
        FROM BoyanPictureTags bt
        JOIN TagPicture tp ON bt.tagID = tp.tagID
        WHERE tp.pictureID = ?
    `

	rows, err := w.db.Query(query, pictureID)
	if err != nil {
		w.logger.Error("Failed to query tags for picture", "pictureID", pictureID, "error", err)
		return nil, fmt.Errorf("failed to query tags for picture: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			w.logger.Error("Failed to scan tag", "pictureID", pictureID, "error", err)
			continue
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		w.logger.Error("Error iterating tag rows for picture", "pictureID", pictureID, "error", err)
		return nil, fmt.Errorf("error iterating tag rows for picture: %w", err)
	}

	return tags, nil
}
