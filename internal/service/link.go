package service

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/repository"
)

const (
	batchSize    = 100
	workersCount = 5
)

type deleteTask struct {
	UserID string
	Alias  string
}

type LinkService struct {
	Storage repository.LinkRepository
	DB      *sql.DB
	taskCh  chan deleteTask
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func NewLinkService(storage repository.LinkRepository, db *sql.DB) *LinkService {
	ctx, cancel := context.WithCancel(context.Background())

	s := &LinkService{
		Storage: storage,
		DB:      db,
		taskCh:  make(chan deleteTask, 1000),
		ctx:     ctx,
		cancel:  cancel,
	}

	s.FanIn(workersCount)

	return s
}

func (s *LinkService) Add(ctx context.Context, link model.Link) error {
	return s.Storage.Add(ctx, link)
}

func (s *LinkService) AddBatch(ctx context.Context, db *sql.DB, links []model.Link) error {
	if db == nil {
		for _, link := range links {
			if err := s.Add(ctx, link); err != nil {
				return err
			}
		}
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, link := range links {
		if err = s.Storage.AddWithTx(ctx, tx, link); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *LinkService) GetLinkByAlias(ctx context.Context, alias string) (*model.Link, error) {
	link, err := s.Storage.GetByAlias(ctx, alias)
	if err != nil {
		return &model.Link{}, err
	}

	return &link, nil
}

func (s *LinkService) GetLinkByURL(ctx context.Context, url string) (*model.Link, error) {
	link, err := s.Storage.GetByURL(ctx, url)
	if err != nil {
		return &model.Link{}, err
	}

	return &link, nil
}

func (s *LinkService) GetLinksByUserID(ctx context.Context, userID string) (*[]model.Link, error) {
	links, err := s.Storage.GetByUserID(ctx, userID)
	if err != nil {
		return &[]model.Link{}, err
	}

	return &links, nil
}

func (s *LinkService) Update(ctx context.Context, link model.Link) error {
	return s.Storage.Update(ctx, link)
}

func (s *LinkService) UpdateBatch(ctx context.Context, db *sql.DB, links []model.Link) error {
	if db == nil {
		for _, link := range links {
			if err := s.Update(ctx, link); err != nil {
				return err
			}
		}
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, link := range links {
		if err = s.Storage.UpdateWithTx(ctx, tx, link); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *LinkService) Stop() {
	s.cancel()
	close(s.taskCh)
	s.wg.Wait()
}

func (s *LinkService) Delete(userID, alias string) {
	select {
	case s.taskCh <- deleteTask{UserID: userID, Alias: alias}:
	case <-s.ctx.Done():
	}
}

func (s *LinkService) FanIn(numWorkers int) {
	taskCh := s.taskCh

	type workerBucket struct {
		bucket []model.Link
	}

	worker := func() {
		defer s.wg.Done()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		wb := workerBucket{
			bucket: make([]model.Link, 0, batchSize),
		}

		flush := func() {
			if len(wb.bucket) == 0 {
				return
			}
			_ = s.UpdateBatch(context.Background(), s.DB, wb.bucket)
			wb.bucket = wb.bucket[:0]
		}

		for {
			select {
			case <-s.ctx.Done():
				flush()
				return
			case task, ok := <-taskCh:
				if !ok {
					flush()
					return
				}

				link, err := s.GetLinkByAlias(context.Background(), task.Alias)
				if err != nil || link.UserID != task.UserID {
					continue
				}

				link.IsDeleted = true
				wb.bucket = append(wb.bucket, *link)

				if len(wb.bucket) >= batchSize {
					flush()
				}
			case <-ticker.C:
				flush()
			}
		}
	}

	s.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}
}
