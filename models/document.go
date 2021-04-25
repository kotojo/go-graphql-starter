package models

import "database/sql"

type Document struct {
	Id     string `json:"id"`
	Text   string `json:"text"`
	UserId string `json:"user"`
}

func NewDocumentService(db *sql.DB) *DocumentService {
	return &DocumentService{
		DB: db,
	}
}

type DocumentService struct {
	DB *sql.DB
}

func (d *DocumentService) Create(documentInput NewDocument) (*Document, error) {
	sqlStatement := `
		INSERT INTO documents (text, user_id)
		VALUES ($1, $2)
		RETURNING id
	`
	id := ""
	err := d.DB.QueryRow(sqlStatement, documentInput.Text, documentInput.UserID).Scan(&id)
	if err != nil {
		return nil, err
	}
	doc := &Document{
		Id:     id,
		Text:   documentInput.Text,
		UserId: documentInput.UserID,
	}
	return doc, nil
}

func (d *DocumentService) GetAll() ([]*Document, error) {
	query := `SELECT id, text, user_id FROM documents;`
	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	documents := make([]*Document, 0)
	for rows.Next() {
		var document Document
		err = rows.Scan(&document.Id, &document.Text, &document.UserId)
		if err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return documents, nil
}
