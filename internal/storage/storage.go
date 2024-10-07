package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type Note struct {
	Id    string `json:"id"`
	Title string `json:"title_note"`
	Text  string `json:"text_note"`
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, err
}

func (s *Storage) CreateNote(title, text, userId string) (string, error) {
	var id string
	query := `INSERT INTO notes(title_note, text_note, userId) VALUES ($1, $2, $3) RETURNING note_id`
	err := s.db.QueryRow(query, title, text, userId).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Storage) GetAllNotes(userId string) (any, error) {
	var notes []Note
	query := `SELECT note_id, title_note, text_note FROM notes WHERE userId = $1;`
	rows, err := s.db.Query(query, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note Note
		err := rows.Scan(&note.Id, &note.Title, &note.Text)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	notes = ([]Note)(notes)
	return notes, nil
}

type User struct {
	Id       string `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

func (s *Storage) CreateUser(login, password string) error {
	query := `INSERT INTO users(login, password) VALUES ($1, $2)`
	_, err := s.db.Exec(query, login, password)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *Storage) CheckUserIfExist(login string) (string, error) {
	var userId string
	query := `SELECT user_id FROM users WHERE login = $1`
	row := s.db.QueryRow(query, login)
	err := row.Scan(&userId)
	if err != nil {
		return "", err
	}
	if row == nil {
		return "", nil
	}

	return userId, nil
}

func (s *Storage) GetUserPassword(login string) (string, error) {
	var password string
	query := `SELECT password FROM users WHERE login = $1`
	err := s.db.QueryRow(query, login).Scan(&password)
	
	if err != nil {
		return "", err
	}
	return password, nil
}
