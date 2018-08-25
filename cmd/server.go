package cmd

import (
	"log"
	"net"
	"os"

	"github.com/almostmoore/gbquestion/question"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var server = &cobra.Command{
	Use:   "serve",
	Short: "Run a questions grpc server",
	Long:  "Server will read .env file and work with it",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := bolt.Open(os.Getenv("DB_PATH"), 0600, nil)
		if err != nil {
			log.Fatalf("Couldn't load database file (%s): %s", os.Getenv("DB_PATH"), err)
		}
		defer db.Close()

		qs := question.NewStorage(db)

		service := question.NewRPCService(qs)
		srv := grpc.NewServer()

		question.RegisterQuestionsServer(srv, service)
		l, err := net.Listen("tcp", os.Getenv("LISTEN"))
		if err != nil {
			log.Fatalf("Couldn't start listening a port: %v", err)
		}

		return srv.Serve(l)
	},
}
