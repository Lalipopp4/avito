package user

import (
	"context"
	"errors"
	"sync"

	"github.com/Lalipopp4/test_api/internal/models"
	"golang.org/x/sync/errgroup"
)

func (s *userService) AddUser(ctx context.Context, user *models.UserRequest) error {

	var (
		wg, errCtx = errgroup.WithContext(ctx)
	)

	// checking if there is no intersection of adding and deleting segments
	wg.Go(func() error {
		if errCtx != nil {
			return errCtx.Err()
		}
		intersection := make(map[string]struct{})
		for _, segment := range user.SegmentsToAdd {
			intersection[segment] = struct{}{}
		}
		for _, val := range user.SegmentsToDelete {
			if _, ok := intersection[val]; ok {
				return errors.New("Intersection in segments to add and delete.")
			}
		}
		return nil
	})

	// checking if segments from add list doesn't exist
	for _, segment := range user.SegmentsToAdd {
		segment := segment
		wg.Go(func() error {
			if errCtx != nil {
				return errCtx.Err()
			}

			exists, err := s.redisRepo.CheckSegments(ctx, segment)
			if err != nil {
				return err
			}
			if !exists {
				return errors.New("Segment " + segment + " already exists.")
			}
			return nil
		})

	}

	// checking if segments from delete list exist
	for _, segment := range user.SegmentsToDelete {
		segment := segment
		wg.Go(func() error {
			if errCtx != nil {
				return errCtx.Err()
			}

			exists, err := s.redisRepo.CheckSegments(ctx, segment)
			if err != nil {
				return err
			}
			if exists {
				return errors.New("Segment " + segment + " doesn't exist.")
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	var (
		segmentsToAdd    = make([]int, len(user.SegmentsToAdd))
		segmentsToDelete = make([]int, len(user.SegmentsToDelete))
		err              error
	)

	for i, val := range user.SegmentsToAdd {
		i, val := i, val
		wg.Go(func() error {

			if errCtx != nil {
				return errCtx.Err()
			}
			segmentsToAdd[i], err = s.psqlRepo.GetSegmentIdByName(ctx, val)
			if err != nil {
				return err
			}
			return nil
		})
	}

	for i, val := range user.SegmentsToDelete {
		i, val := i, val
		wg.Go(func() error {

			if errCtx != nil {
				return errCtx.Err()
			}
			segmentsToDelete[i], err = s.psqlRepo.GetSegmentIdByName(ctx, val)
			if err != nil {
				return err
			}
			return nil
		})

	}

	if err := wg.Wait(); err != nil {
		return err
	}

	// adding user in segments
	err = s.psqlRepo.AddUser(ctx, user, segmentsToAdd, segmentsToDelete)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) GetSegmentsByUser(ctx context.Context, user *models.User) ([]string, error) {
	segmentIds, err := s.psqlRepo.Get(ctx, user.Id, true, -1)
	if err != nil {
		return nil, err
	}
	var (
		segments = make([]string, len(segmentIds))
		wg       sync.WaitGroup
	)
	wg.Add(len(segmentIds))
	for i, segment := range segmentIds {
		defer wg.Done()
		i, segment := i, segment
		go func() {
			segments[i], err = s.redisRepo.GetElementById(ctx, segment)
			if err != nil {
				return
			}
		}()
	}
	wg.Wait()
	return segments, nil
}
