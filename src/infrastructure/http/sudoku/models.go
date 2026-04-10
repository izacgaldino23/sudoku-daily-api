package sudoku

import (
	"time"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
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

	MySolvesResponse struct {
		Solves []Solve `json:"solves"`
	}

	Solve struct {
		ID        vo.UUID   `json:"id"`
		Duration  int       `json:"duration"`
		StartedAt time.Time `json:"started_at"`
		Size      int       `json:"size"`
		Date      time.Time `json:"date"`
	}
)

func (g *SudokuResponse) FromDomain(s *entities.Sudoku, playToken string) {
	g.ID = s.ID.String()
	g.Size = s.GetSize()
	g.Board = BoardFromDomain(s.Board)
	g.Date = s.Date.Format(time.DateOnly)
	g.PlayToken = playToken
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
		Solution: s.Solution,
		UserID:   userID,
	}
}

func (m *MySolvesResponse) FromDomain(solves []entities.Solve) {
	m.Solves = make([]Solve, 0)

	for _, solve := range solves {
		m.Solves = append(m.Solves, Solve{
			ID:        solve.ID,
			Date:      solve.SudokuDate,
			Duration:  solve.Duration,
			StartedAt: solve.StartedAt,
			Size:      solve.Size,
		})
	}
}
