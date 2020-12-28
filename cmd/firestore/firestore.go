// Copyright © 2020 Elis Lulja
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package firestore

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/SunSince90/kube-scraper-backend/pkg/firestore"
	"github.com/SunSince90/kube-scraper-backend/pkg/pb"
	"github.com/SunSince90/kube-scraper-backend/pkg/server"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	log zerolog.Logger
)

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	log = zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// NewFirestoreCommand returns an instance of the firestore command
func NewFirestoreCommand() *cobra.Command {
	opts := &firestoreOptions{}
	cmd := &cobra.Command{
		Use: `firestore --service-account-path|-s <path> --project-name|-n <project-name>
--chats-collection|-c <chats-collection-name>`,
		Example: `backend firestore -s ./credentials/service-account.json -n my-project-name
-c chats`,
		Short: "establish a connection with firestore",
		Long:  `firestore connects to firestore and loads data from there`,
		Run: func(cmd *cobra.Command, args []string) {
			addr, err := cmd.Parent().Flags().GetString("address")
			if err != nil {
				log.Fatal().Msg("could not get address flag")
			}
			port, err := cmd.Parent().Flags().GetInt("port")
			if err != nil {
				log.Fatal().Msg("could not get port flag")
			}

			opts.address = addr
			opts.port = port
			runFirestore(opts)
		},
	}

	// -- Flags
	cmd.Flags().StringVarP(&opts.serviceAccountPath, "service-account-path", "s", "", "the path to the service account")
	cmd.Flags().StringVarP(&opts.projectID, "project-id", "i", "", "the id of the firebase project")
	cmd.Flags().StringVarP(&opts.chatsCollection, "chats-collection", "c", "", "the name of the collections where chats are stored")
	cmd.Flags().BoolVar(&opts.listen, "listen", true, "whether to listening to chat updates or not")

	// -- Mark as required
	cmd.MarkFlagRequired("service-account-path")
	cmd.MarkFlagRequired("project-name")
	cmd.MarkFlagRequired("chats-collection")

	return cmd
}

func runFirestore(opts *firestoreOptions) {
	// -- Init
	l := log.With().Str("func", "runFirestore").Logger()
	l.Info().Msg("starting...")
	l.Debug().Msg("debug mode requested")

	ctx, canc := context.WithCancel(context.Background())
	defer canc()
	stopChan := make(chan struct{})
	exitChan := make(chan struct{})
	endpoint := fmt.Sprintf("%s:%d", opts.address, opts.port)

	// -- Get the backend
	fs, err := firestore.NewBackend(ctx, opts.serviceAccountPath, &firestore.Options{
		ProjectID:       opts.projectID,
		ChatsCollection: opts.chatsCollection,
		UseCache:        true,
	})
	if err != nil {
		l.Fatal().Err(err).Msg("error while getting firestore as backend")
	}
	defer fs.Close()

	if opts.listen {
		go fs.ListenForChats(ctx, stopChan)
	} else {
		close(stopChan)
	}

	// -- Start the grpc server
	serv, err := server.New(fs)
	if err != nil {
		l.Fatal().Err(err).Msg("error while starting the backend")
		return
	}

	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to listen")
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBackendServer(grpcServer, serv)
	go func() {
		grpcServer.Serve(lis)
		close(exitChan)
	}()

	// -- Graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	<-signalChan
	fmt.Println()
	l.Info().Msg("exit requested")
	grpcServer.GracefulStop()
	<-stopChan
	<-exitChan
	l.Info().Msg("goodbye!")
}
