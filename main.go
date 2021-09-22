package main

import (
	"code-complexity/calculate"
	"code-complexity/options"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const VERSION = "1.0.0"

func main() {
	cli.AppHelpTemplate =
		`NAME:
   {{.Name}} - {{.Version}} - {{.Usage}}

USAGE:
   {{.Name}}{{range .Flags}}{{if and (not (eq .Name "help")) (not (eq .Name "version")) }} {{if .Required}}--{{.Name}} value{{end}}{{end}}{{end}} [optional flags]

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stderr)
	app := &cli.App{
		Name:    "code-complexity",
		Usage:   "Estimate source code complexity",
		Flags:   options.Flags,
		Version: VERSION,
		Action: func(ctx *cli.Context) error {
			opts, err := options.ParseOptions(ctx)
			if err != nil {
				return err
			}
			summary, err := calculate.Complexity(opts)
			if err != nil {
				return err
			}
			asJson, err := json.Marshal(summary)
			if err != nil {
				return fmt.Errorf("failed to serialize summary to json: %v", err)
			}
			log.Printf("completed successfully at %v", opts.CodePath)
			println(string(asJson))
			if len(opts.OutputPath) > 0 {
				err = os.WriteFile(opts.OutputPath, asJson, 0777)
				if err != nil {
					return fmt.Errorf("failed to write output to %v: %v", opts.OutputPath, err)
				}
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Printf("failed: %v", err)
		os.Exit(1)
	}
}
