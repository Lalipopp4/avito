package redis

import (
	"context"
	"strconv"
)

const (
	SETNAME = "segments"
	MAPNAME = "segmentID"
)

func (r *repository) CountUsers(ctx context.Context) (int, error) {
	res, err := r.cur.Get(ctx, "users").Result()
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(res)
}

func (r *repository) CheckSegments(ctx context.Context, segment string) (bool, error) {
	res, err := r.cur.SIsMember(ctx, SETNAME, segment).Result()
	if err != nil {
		return false, err
	}
	return res, nil
}

func (r *repository) GetElementById(ctx context.Context, id int) (string, error) {
	res, err := r.cur.HGet(ctx, MAPNAME, strconv.Itoa(id)).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}
