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
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
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
		Run:   runFirestore,
	}

	// -- Flags
	cmd.Flags().StringVarP(&opts.serviceAccountPath, "service-account-path", "s", "", "the path to the service account")
	cmd.Flags().StringVarP(&opts.projectName, "project-name", "n", "", "the name of the firebase project")
	cmd.Flags().StringVarP(&opts.chatsCollection, "chats-collection", "c", "", "the name of the collections where chats are stored")

	// -- Mark as required
	cmd.MarkFlagRequired("service-account-path")
	cmd.MarkFlagRequired("project-name")
	cmd.MarkFlagRequired("chats-collection")

	return cmd
}

func runFirestore(cmd *cobra.Command, args []string) {
	l := log.With().Str("func", "runFirestore").Logger()
	l.Debug().Msg("debug mode requested")

	stopChan := make(chan struct{})
	ctx, canc := context.WithCancel(context.Background())

	// TODO
	_ = ctx
	_ = stopChan

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

	canc()

	l.Info().Msg("goodbye!")
}
