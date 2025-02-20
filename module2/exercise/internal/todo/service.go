package todo

import (
	"context"

	"github.com/google/uuid"
	"github.com/orlandorode97/grpc-stuff/module2/exercise/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	tasks map[string]*proto.Task
	proto.UnimplementedTodoServiceServer
}

func New() *Service {
	return &Service{
		tasks: make(map[string]*proto.Task),
	}
}

func (s *Service) AddTask(ctx context.Context, req *proto.AddTaskRequest) (*proto.AddTaskResponse, error) {
	if req.GetTask() == "" {
		return nil, status.Error(codes.InvalidArgument, "task description cannot be empty")
	}

	for _, task := range s.tasks {
		if task.GetTask() == req.GetTask() {
			return nil, status.Errorf(codes.AlreadyExists, "the task: %s already exists", req.GetTask())
		}
	}

	id := uuid.NewString()
	s.tasks[id] = &proto.Task{
		Id:   id,
		Task: req.GetTask(),
	}

	return &proto.AddTaskResponse{
		Id: id,
	}, nil
}

func (s *Service) CompleteTask(ctx context.Context, req *proto.CompleteTaskRequest) (*proto.CompleteTaskResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "task id cannot be empty")
	}

	if _, ok := s.tasks[req.GetId()]; !ok {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	delete(s.tasks, req.GetId())

	return &proto.CompleteTaskResponse{}, nil
}

func (s *Service) ListTasks(ctx context.Context, req *proto.ListTasksRequest) (*proto.ListTasksResponse, error) {
	tasks := make([]*proto.Task, 0)
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return &proto.ListTasksResponse{
		Tasks: tasks,
	}, nil
}
