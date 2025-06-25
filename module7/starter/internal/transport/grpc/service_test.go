package grpc_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chat_mock "github.com/orlandorode97/grpc-stuff/module7/starter/internal/mocks/chat"
	"github.com/orlandorode97/grpc-stuff/module7/starter/internal/transport/grpc"
)

func TestNewService(t *testing.T) {
	t.Run("returns error if chat service is nil", func(t *testing.T) {
		service, err := grpc.NewService(nil)
		require.Error(t, err)

		assert.Empty(t, service)
	})

	t.Run("successfully initialises grpc service", func(t *testing.T) {
		service, err := grpc.NewService(chat_mock.NewMockChatService(gomock.NewController(t)))
		require.NoError(t, err)

		assert.NotEmpty(t, service)
	})
}

func TestService_SendMessage(t *testing.T) {
	// test cases

	for _, tc := range []struct{
		name string
		req *proto.
	} {

	}
}

func TestService_Subscribe(t *testing.T) {

}
