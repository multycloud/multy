package deploy

import (
	"context"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
	"github.com/stretchr/testify/mock"
)

type MockTerraformCommand struct {
	mock.Mock
}

func (m *MockTerraformCommand) Init(ctx context.Context, dir string) error {
	args := m.Called(ctx, dir)
	return args.Error(0)
}

func (m *MockTerraformCommand) Apply(ctx context.Context, dir string, resources []string) error {
	args := m.Called(ctx, dir, resources)
	return args.Error(0)
}

func (m *MockTerraformCommand) Plan(ctx context.Context, dir string) (string, error) {
	args := m.Called(ctx, dir)
	return args.String(0), args.Error(1)
}

func (m *MockTerraformCommand) Refresh(ctx context.Context, dir string) error {
	args := m.Called(ctx, dir)
	return args.Error(0)
}

func (m *MockTerraformCommand) GetState(ctx context.Context, userId string, reader db.TfStateReader) (*output.TfState, error) {
	args := m.Called(ctx, userId, reader)
	return args.Get(0).(*output.TfState), args.Error(1)
}
