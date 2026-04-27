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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		tries     int
		wg        sync.WaitGroup
		startTime = time.Now()
		jobs      = make(chan int)
		result    = make(chan *entities.Board, 1)
		success   atomic.Bool
		mu        sync.Mutex
	)

	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go s.worker(ctx, jobs, result, &success, &mu, &tries, board, r, startTime, &wg, cancel)
	}

	// produce jobs
	go func() {
		defer close(jobs)
		for target := targetToHide; target >= targetToHide-maxTargetDecrease && !success.Load(); target-- {
			for i := 0; i < maxTries && !success.Load(); i++ {
				select {
				case <-ctx.Done():
					return
				case jobs <- target:
				}
			}
		}
	}()

	select {
	case clone := <-result:
		board.Board = *clone
		logging.Log(ctx).Info().Msgf("Successfully hidden %v cells, tries: %v, board size: %v, time: %v", targetToHide, tries, board.Size, time.Since(startTime))
		success.Store(true)
	case <-ctx.Done():
		logging.Log(ctx).Error().Msgf("Failed to hide %v cells, tries: %v, board size: %v", targetToHide, tries, board.Size)
		success.Store(false)
	}

	wg.Wait()
	return success.Load()
}

func (s *hideBacktracking) worker(ctx context.Context, jobs <-chan int, result chan<- *entities.Board, success *atomic.Bool, mu *sync.Mutex, tries *int, board *entities.Sudoku, r *rand.Rand, startTime time.Time, wg *sync.WaitGroup, cancel context.CancelFunc) {
	defer wg.Done()

	solver := newSolver()
	for {
		select {
		case targetToHide := <-jobs:
			if success.Load() {
				return
			}

			clone, err := board.Board.Clone()
			if err != nil {
				continue
			}

			mu.Lock()
			*tries++
			logging.Log(ctx).Debug().Msgf("Executing job on attempt %v", *tries)
			mu.Unlock()

			if s.hideCells(solver, clone, r, targetToHide) && !success.Load() {
				select {
				case result <- clone:
					success.Store(true)
					logging.Log(ctx).Info().Msgf("Successfully hidden %v cells, tries: %v, board size: %v, time: %v", targetToHide, tries, board.Size, time.Since(startTime))
					cancel()
				default:
				}
				return
			}
		case <-ctx.Done():
			return
		}
	}
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
