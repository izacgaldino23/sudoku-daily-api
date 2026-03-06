package persistence

import (
	"math"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"time"

	"github.com/uptrace/bun"
)

type (
	Sudoku struct {
		bun.BaseModel `bun:"table:sudoku"`

		ID         string    `bun:"id,pk"`
		Size       int       `bun:",notnull"`
		Difficulty string    `bun:",notnull"`
		Board      []byte    `bun:"type:,notnull"`
		Solution   []byte    `bun:"type:,notnull"`
		Date       time.Time `bun:"type:date,notnull"`
	}

	User struct {
		bun.BaseModel `bun:"table:user"`

		ID            string  `bun:"id,pk"`
		Username      string  `bun:",unique,notnull"`
		Email         string  `bun:",unique,notnull"`
		PasswordHash  []byte  `bun:",notnull"`
		Provider      string  `bun:",notnull"`
		ProviderID    *string `bun:",notnull"`
		EmailVerified bool    `bun:",notnull"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	RefreshToken struct {
		bun.BaseModel

		ID        string    `bun:"id,pk"`
		UserID    string    `bun:"user_id,notnull"`
		TokenHash string    `bun:",unique,notnull"`
		ExpiresAt time.Time `bun:"type:timestamp,notnull"`
		Revoked   bool      `bun:"notnull"`
		CreatedAt time.Time `bun:"type:timestamp,notnull"`
	}
)

func (s *Sudoku) FromDomain(sudoku *entities.Sudoku) {
	s.ID = sudoku.ID
	s.Size = sudoku.GetSize()
	s.Difficulty = string(sudoku.Difficulty)
	s.Board = boardFromDomain(&sudoku.Board)
	s.Solution = boardFromDomain(&sudoku.Solution)
	s.Date = sudoku.Date
}

func (s *Sudoku) ToDomain() *entities.Sudoku {
	return &entities.Sudoku{
		ID:         s.ID,
		Size:       entities.BoardSize(s.Size),
		Difficulty: entities.Difficulty(s.Difficulty),
		Board:      boardToDomain(s.Board),
		Solution:   boardToDomain(s.Solution),
		Date:       s.Date,
	}
}

func boardToDomain(boardData []byte) entities.Board {
	size := int(math.Sqrt(float64(len(boardData))))

	boardFilled := make([][]int, size)

	for i := 0; i < size; i++ {
		row := make([]int, size)
		for j := 0; j < size; j++ {
			row[j] = int(boardData[i*size+j])
		}
		boardFilled[i] = row
	}

	return entities.NewFilledBoard(boardFilled)
}

func boardFromDomain(board *entities.Board) []byte {
	size := board.GetSize()
	linearBoard := make([]byte, 0, size*size)

	for _, row := range board.GetBoard() {
		for _, val := range row {
			linearBoard = append(linearBoard, byte(val))
		}
	}

	return linearBoard
}

func (u *User) FromDomain(user *entities.User) {
	u.ID = string(user.ID)
	u.Email = user.Email.String()
	u.Username = user.Username
	u.Provider = string(user.Provider)
	u.ProviderID = user.ProviderID
	u.EmailVerified = user.EmailVerified
	u.CreatedAt = user.CreatedAt

	if user.PasswordHash != nil {
		u.PasswordHash = []byte(*user.PasswordHash)
	}
}

func (u *User) ToDomain() *entities.User {
	var passwordHash *string
	if u.PasswordHash != nil {
		hashStr := string(u.PasswordHash)
		passwordHash = &hashStr
	}
	return &entities.User{
		ID:            vo.UUID(u.ID),
		Email:         entities.Email(u.Email),
		Username:      u.Username,
		PasswordHash:  passwordHash,
		Provider:      entities.AuthProvider(u.Provider),
		ProviderID:    u.ProviderID,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
	}
}
