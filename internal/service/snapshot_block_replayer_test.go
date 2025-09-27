package service

import (
	"context"
	"fmt"
	"testing"
)

func TestSnapshotBlockReplayer_Init(t *testing.T) {
	validator := "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg"
	delegator := "omni195tf7tuvxagrsq7seh8579qrjvxc6j9wgnpz86"
	startTime := int64(1748507477910)
	endTime := int64(1748598144482)
	ctx := context.Background()
	snapshotBlockReplayer := NewRewardSnapshotTimetangeReplayer(validator, delegator, startTime, endTime)
	snapshotBlockReplayer.Init(ctx)
	snapshotBlockReplayer.Replay(ctx)
	fmt.Println(snapshotBlockReplayer.GetResultRewardAmount())
}
