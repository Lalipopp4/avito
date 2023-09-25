package user

import (
	"context"
	"errors"
	"sync"

	"github.com/Lalipopp4/test_api/internal/models"
)

func (s *userService) AddUser(ctx context.Context, user *models.UserRequest) error {

	var (
		c    = make(chan error)
		quit = make(chan struct{})
		wg   sync.WaitGroup
	)
	defer close(c)
	defer close(quit)

	// checking if there is no intersection of adding and deleting segments
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-quit:
			wg.Done()
			return
		default:
		}
		intersection := make(map[string]struct{})
		for _, segment := range user.SegmentsToAdd {
			intersection[segment] = struct{}{}
		}
		for _, val := range user.SegmentsToDelete {
			if _, ok := intersection[val]; ok {
				c <- errors.New("Intersection in segments to add and delete.")
				quit <- struct{}{}
				return
			}
		}
	}()

	// checking if segments from add list doesn't exist
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, segment := range user.SegmentsToAdd {

			select {
			case <-quit:
				wg.Done()
				return
			default:
			}
			exists, err := s.redisRepo.CheckSegments(ctx, segment)
			if err != nil {
				c <- err
				quit <- struct{}{}
				return
			}
			if !exists {
				c <- errors.New("Segment " + segment + " already exists.")
				quit <- struct{}{}
				return
			}
		}

	}()

	// checking if segments from delete list exist
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, segment := range user.SegmentsToDelete {

			select {
			case <-quit:
				wg.Done()
				return
			default:
			}
			exists, err := s.redisRepo.CheckSegments(ctx, segment)
			if err != nil {
				c <- err
				quit <- struct{}{}
				return
			}
			if exists {
				c <- errors.New("Segment " + segment + " doesn't exist.")
				quit <- struct{}{}
				return
			}
		}
	}()

	wg.Wait()
	select {
	case err := <-c:
		return err
	default:
	}

	var (
		segmentsToAdd    = make([]int, len(user.SegmentsToAdd))
		segmentsToDelete = make([]int, len(user.SegmentsToDelete))
		errGo, err       error
	)

	wg.Add(len(user.SegmentsToDelete) + len(user.SegmentsToAdd))

	for i, val := range user.SegmentsToAdd {
		go func() {
			defer wg.Done()
			segmentsToAdd[i], err = s.psqlRepo.GetSegmentIdByName(ctx, val)
			if err != nil {
				errGo = err
			}
		}()
	}

	for i, val := range user.SegmentsToDelete {
		go func() {
			defer wg.Done()
			segmentsToDelete[i], err = s.psqlRepo.GetSegmentIdByName(ctx, val)
			if err != nil {
				errGo = err
			}

		}()

	}

	wg.Wait()
	if errGo != nil {
		return errGo
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
	segments := make([]string, len(segmentIds))
	for i, segment := range segmentIds {
		segments[i], err = s.redisRepo.GetElementById(ctx, segment)
		if err != nil {
			return nil, err
		}
	}
	return segments, nil
}
