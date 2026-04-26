package strategies

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/infrastructure/logging"
)

type (
	HideStrategy interface {
		Hide(ctx context.Context, board *entities.Sudoku, r *rand.Rand) bool
	}

	hideBacktracking struct {
	}
)

func NewHideStrategy() HideStrategy {
	return &hideBacktracking{}
}

func (s *hideBacktracking) Hide(ctx context.Context, board *entities.Sudoku, r *rand.Rand) bool {
	targetToHide := s.defineToHideCount(board, r)

	const (
		maxTries          = 1000
		maxTargetDecrease = 15
		parallelism       = 4
	)

	var (
		tries     int
		wg        sync.WaitGroup
		startTime = time.Now()
		jobs      = make(chan struct{}, parallelism)
		result    = make(chan *entities.Board, 1)
		success   atomic.Bool
		mu        sync.Mutex
	)

	wg.Add(parallelism)
	for i := 0; i < parallelism; i++ {
		go func() {
			defer wg.Done()

			solver := newSolver()
			for range jobs {
				if success.Load() {
					return
				}

				clone, err := board.Board.Clone()
				mu.Lock()
				tries++
				logging.Log(ctx).Debug().Msgf("Executing job on attempt %v", tries)
				mu.Unlock()

				if err != nil {
					continue
				}
				if s.hideCells(solver, clone, r, targetToHide) && !success.Load() {
					success.Store(true)
					logging.Log(ctx).Info().Msgf("Successfully hidden %v cells, tries: %v, board size: %v, time: %v", targetToHide, tries, board.Size, time.Since(startTime))
					select {
					case result <- clone:
					default:
					}
					return
				}
			}
		}()
	}

	for target := targetToHide; target >= targetToHide-maxTargetDecrease && !success.Load(); target-- {
		for i := 0; i < maxTries && !success.Load(); i++ {
			jobs <- struct{}{}

			select {
			case clone := <-result:
				board.Board = *clone
				close(jobs)
				wg.Wait()
				return true
			default:
			}
		}
	}

	close(jobs)
	wg.Wait()

	logging.Log(ctx).Error().Msgf("Failed to hide %v cells, tries: %v, board size: %v", targetToHide, tries, board.Size)

	return false
}

func (s *hideBacktracking) hideCells(solver *solver, board *entities.Board, r *rand.Rand, target int) bool {
	var (
		hiddenCount = 0
		cells       = s.getCellShuffled(board, r)
	)

	for i := range cells {
		if hiddenCount >= target || i+(target-hiddenCount) >= len(cells) {
			break
		}

		cell := cells[i]
		row, col, value := cell[0], cell[1], cell[2]

		board.SetCell(row, col, 0)
		if solver.Execute(board) == 1 {
			hiddenCount++
		} else {
			board.SetCell(row, col, value)
		}
	}

	hidedTarget := hiddenCount >= target

	if !hidedTarget {
		for i := range cells {
			cell := cells[i]
			row, col, value := cell[0], cell[1], cell[2]

			board.SetCell(row, col, value)
		}
	}

	return hidedTarget
}

func (s *hideBacktracking) defineToHideCount(board *entities.Sudoku, r *rand.Rand) int {
	if board.Difficulty == "" {
		// generate random difficulty
		difficulties := []entities.Difficulty{entities.DifficultyEasy, entities.DifficultyMedium, entities.DifficultyHard}
		board.Difficulty = difficulties[r.Intn(len(difficulties))]
	}

	min, max := entities.GetClue(board.Size, board.Difficulty)

	clueCount := r.Intn(max-min+1) + min
	return board.GetSize()*board.GetSize() - clueCount
}

func (s *hideBacktracking) getCellShuffled(board *entities.Board, r *rand.Rand) [][3]int {
	cellReference := make([][3]int, 0)

	for row := range board.GetBoard() {
		for col := range board.GetBoard()[row] {
			value := board.GetCell(row, col)
			if value != 0 {
				cellReference = append(cellReference, [3]int{row, col, value})
			}
		}
	}

	r.Shuffle(len(cellReference), func(i, j int) {
		cellReference[i], cellReference[j] = cellReference[j], cellReference[i]
	})

	return cellReference
}
