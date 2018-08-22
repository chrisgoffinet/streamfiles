package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/chrisgoffinet/streamfiles/api"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var port int
var dataDir string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "an example server that listens over grpc",
	Run: func(cmd *cobra.Command, args []string) {
		port := fmt.Sprintf(":%d", port)
		// create a listener on TCP port
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		// create a server instance
		s := api.New(dataDir)

		// create a gRPC server object
		grpcServer := grpc.NewServer()

		// attach the Ping service to the server
		api.RegisterStorageServer(grpcServer, s)
		log.Printf("listening on port: %s", port)

		// start the server
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&port, "port", "p", 7777, "port to listen on for the grpc server")
	serverCmd.Flags().StringVarP(&dataDir, "datadir", "d", "data", "directory to store files")
}
