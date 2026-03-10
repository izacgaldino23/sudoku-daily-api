package sudoku

import (
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"time"
)

type (
	GetDailySudokuRequest struct {
		Size string `query:"size" validate:"required,oneof=four six nine"`
	}

	SudokuResponse struct {
		ID        string  `json:"id"`
		Size      int     `json:"size"`
		Board     []Cell  `json:"board"`
		Date      string  `json:"date"`
		PlayToken string  `json:"session_token,omitempty"`
		SessionID vo.UUID `json:"session_id,omitempty"`
	}

	Cell struct {
		Row   int `json:"row"`
		Col   int `json:"col"`
		Value int `json:"value"`
	}

	VerifySolutionRequest struct {
		Solution  [][]int `json:"solution" validate:"required"`
		PlayToken string  `json:"play_token" validate:"required"`
	}

	VerifySolutionResponse struct {
		Valid      bool `json:"valid"`
		StartedAt  int  `json:"started_at"`
		FinishedAt int  `json:"finished_at"`
	}
)

func (g *GetDailySudokuRequest) GetSize() int {
	switch g.Size {
	case "four":
		return 4
	case "six":
		return 6
	case "nine":
		return 9
	default:
		return 0
	}
}

func (g *SudokuResponse) FromDomain(s *entities.Sudoku, playToken string, sessionID vo.UUID) {
	g.ID = s.ID.String()
	g.Size = s.GetSize()
	g.Board = BoardFromDomain(s.Board)
	g.Date = s.Date.Format(time.DateOnly)
	g.PlayToken = playToken
	g.SessionID = sessionID
}

func BoardFromDomain(board entities.Board) []Cell {
	var cells = make([]Cell, 0)

	for i, row := range board.GetBoard() {
		for j := range row {
			val := board.GetCell(i, j)
			if val == 0 {
				continue
			}
			cells = append(cells, Cell{
				Row:   i,
				Col:   j,
				Value: val,
			})
		}
	}

	return cells
}

func (s *VerifySolutionRequest) ToDomain(userID vo.UUID) *entities.Solve {
	return &entities.Solve{
		Solution:    s.Solution,
		UserID:      userID,
		CompletedAt: time.Now(),
	}
}
