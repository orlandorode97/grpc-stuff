package uploading

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/orlandorode97/grpc-stuff/module4/exercise/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	files map[string][]byte
	proto.UnimplementedFileUploadServiceServer
}

func New() *Service {
	return &Service{
		files: make(map[string][]byte),
	}
}

func (s *Service) DownloadFile(req *proto.DownloadFileRequest, stream proto.FileUploadService_DownloadFileServer) error {
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name cannot be empty")
	}

	if _, ok := s.files[req.GetName()]; !ok {
		return status.Error(codes.NotFound, "file name not found")
	}

	content := s.files[req.GetName()]

	reader := bytes.NewReader(content)

	buffer := make([]byte, 10)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}

		if err = stream.Send(&proto.DownloadFileResponse{
			Content: buffer[:n],
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) UploadFile(stream proto.FileUploadService_UploadFileServer) error {
	name := fmt.Sprintf("%s.png", uuid.New().String())
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		if len(req.Content) == 0 {
			return status.Error(codes.InvalidArgument, "content is empty; nothing to upload")
		}

		s.files[name] = append(s.files[name], req.Content...)
	}

	if err := stream.SendAndClose(&proto.UploadFileResponse{
		Name: name,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetListOfFiles(ctx context.Context, req *proto.GetListOfFilesRequest) (*proto.GetListOfFilesResponse, error) {
	files := make([]string, 0)
	for i := range s.files {
		files = append(files, i)
	}

	return &proto.GetListOfFilesResponse{
		Files: files,
	}, nil
}
