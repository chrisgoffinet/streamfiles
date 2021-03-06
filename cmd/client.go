package cmd

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/chrisgoffinet/streamfiles/api"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var conn *grpc.ClientConn
var writing = true
var chunkSize = 1024 * 1024 // 1MB chunk size

var clientCmd = &cobra.Command{
	Use:   "client [host:port] [filename]",
	Short: "an example client that connects over grpc",
	Args:  cobra.MinimumNArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		hostport, filename := args[0], args[1]

		conn, err := grpc.Dial(hostport, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		c := api.NewStorageClient(conn)

		// open the file we should upload
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}

		// start uploading by chunks to server
		stream, err := c.Upload(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		// allocate a []byte by defined chunk size
		buf := make([]byte, chunkSize)
		start := time.Now()

		for writing {
			n, err := file.Read(buf)
			if err != nil {
				if err == io.EOF {
					writing = false
				} else {
					log.Fatal(err)
				}
			}
			stream.Send(&api.Chunk{
				Filename: filename,
				Content:  buf[:n],
			})
		}

		// close
		status, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatal(err)
		}
		duration := time.Since(start)
		log.Println(status)
		log.Printf("total time for upload: %s\n", duration)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
