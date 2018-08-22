package api

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Server represents the gRPC server
type Server struct {
	DataDirectory string
}

// New handles creating a new Server
func New(dataDir string) *Server {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0700); err != nil {
			log.Fatal(err)
		}
	}
	return &Server{DataDirectory: dataDir}
}

// Upload handles storing a file to disk
func (s *Server) Upload(stream Storage_UploadServer) error {
	var f *os.File
	var totalBytes int

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
			return err
		}
		if f == nil {
			f, err = os.Create(filepath.Join(s.DataDirectory, chunk.Filename))
			if err != nil {
				return err
			}
			defer f.Close()
		}
		f.Write(chunk.Content)
		totalBytes += len(chunk.Content)
	}

	err := stream.SendAndClose(&UploadStatus{
		Message: "upload received with success",
		Code:    UploadStatusCode_Ok,
	})
	f.Sync()

	log.Printf("wrote file: %s (%d bytes)\n", f.Name(), totalBytes)
	return err
}
