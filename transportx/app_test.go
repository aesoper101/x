package transportx

import (
	"context"
	"fmt"
	"github.com/aesoper101/x/interrupt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := interrupt.WithCancel(context.Background())
	go func() {
		timer := time.NewTimer(5 * time.Second)
		select {
		case <-timer.C:
			fmt.Println("timeout")
			cancel()
		}
	}()

	mockServer := NewMockServer(ctrl)
	mockServer.EXPECT().Start(gomock.Any()).Do(
		func(ctx context.Context) {
			fmt.Println("mock server start")
		},
	).Return(nil).AnyTimes()
	mockServer.EXPECT().Stop(gomock.Any()).Do(
		func(ctx context.Context) {
			fmt.Println("mock server stop")
		},
	).Return(nil).AnyTimes()

	mockServer2 := NewMockServer(ctrl)
	mockServer2.EXPECT().Start(gomock.Any()).DoAndReturn(
		func(ctx2 context.Context) error {
			fmt.Println("mock server start2", ctx2)
			info, ok := FromContext(ctx2)
			fmt.Println("=========", info, ok)
			return nil
		},
	).AnyTimes()
	mockServer2.EXPECT().Stop(gomock.Any()).DoAndReturn(
		func(ctx2 context.Context) error {
			fmt.Println("mock server stop2")
			return nil
		},
	).AnyTimes()

	err := Run(
		WithContext(ctx), WithServers(mockServer, mockServer2), AfterStart(
			func(ctx context.Context) error {
				fmt.Println("mock after start2")
				info, ok := FromContext(ctx)
				fmt.Println("--------------", info, ok)
				return nil
			},
		),
	)
	require.NoError(t, err)
}
