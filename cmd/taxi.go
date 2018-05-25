// Copyright 2017 Pilosa Corp.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright
// notice, this list of conditions and the following disclaimer in the
// documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
// contributors may be used to endorse or promote products derived
// from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
// CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
// CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH
// DAMAGE.

package cmd

import (
	"io"
	"log"
	"time"

	"github.com/pilosa/pdk/usecase/taxi"
	"github.com/spf13/cobra"
)

// TaxiMain is wrapped by NewTaxiCommand. It is exported for testing purposes.
var TaxiMain *taxi.Main

// NewTaxiCommand wraps taxi.Main with cobra.Command for use from a CLI.
func NewTaxiCommand(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	TaxiMain = taxi.NewMain()
	taxiCommand := &cobra.Command{
		Use:   "taxi",
		Short: "taxi - import taxi data to pilosa",
		Long:  `TODO`,
		RunE: func(cmd *cobra.Command, args []string) error {
			start := time.Now()
			err := TaxiMain.Run()
			if err != nil {
				return err
			}
			log.Println("Done: ", time.Since(start))
			select {}
		},
	}
	flags := taxiCommand.Flags()
	flags.IntVarP(&TaxiMain.Concurrency, "concurrency", "c", 8, "Number of goroutines fetching and parsing")
	flags.IntVarP(&TaxiMain.FetchConcurrency, "fetch-concurrency", "e", 8, "Number of goroutines fetching and parsing")
	flags.IntVarP(&TaxiMain.BufferSize, "buffer-size", "b", 1000000, "Size of buffer for importers - heavily affects memory usage")
	flags.BoolVarP(&TaxiMain.UseReadAll, "use-read-all", "", false, "Setting to true uses much more memory, but ensures that an entire file can be read before beginning to parse it.")
	flags.StringVarP(&TaxiMain.PilosaHost, "pilosa", "p", "localhost:10101", "Pilosa host")
	flags.StringVarP(&TaxiMain.Index, "index", "i", TaxiMain.Index, "Pilosa db to write to")
	flags.StringVarP(&TaxiMain.URLFile, "url-file", "f", "usecase/taxi/urls-short.txt", "File to get raw data urls from. Urls may be http or local files.")

	return taxiCommand
}

func init() {
	subcommandFns["taxi"] = NewTaxiCommand
}
