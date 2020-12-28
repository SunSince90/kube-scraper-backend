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

package root

import (
	"os"

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

// NewRootCommand returns an instance of the telegram command
func NewRootCommand() *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:     `backend [backend]`,
		Example: `backend firestore [options]`,
		Short:   "establish a connection with the backend",
		Long: `backend establishes a connection with the backend so that all other
scrapers can use this as the place where to get data.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if opts.debug {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
		},
	}

	// -- Flags
	cmd.PersistentFlags().StringVarP(&opts.address, "address", "a", "", "the address to serve from")
	cmd.PersistentFlags().IntVarP(&opts.port, "port", "p", 80, "the port to serve from")
	cmd.PersistentFlags().BoolVar(&opts.debug, "debug", false, "whether to log debug lines")

	// -- Mark as required
	cmd.MarkPersistentFlagRequired("address")
	cmd.MarkPersistentFlagRequired("port")

	return cmd
}
