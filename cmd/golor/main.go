package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := execute(os.Stdout, os.Stderr, os.Args[1:]); err != nil {
		os.Exit(1)
	}
}

func execute(stdout, stderr io.Writer, args []string) error {
	cmd := newRootCommand(stdout, stderr)
	cmd.SetArgs(args)
	err := cmd.Execute()
	if err == nil {
		return nil
	}
	jsonOut := false
	if flag := cmd.Flags().Lookup("json"); flag != nil {
		jsonOut = flag.Value.String() == "true"
	}
	if flag := cmd.PersistentFlags().Lookup("json"); flag != nil {
		jsonOut = jsonOut || flag.Value.String() == "true"
	}
	if jsonOut {
		_ = json.NewEncoder(stderr).Encode(map[string]string{"error": err.Error()})
	} else {
		_, _ = fmt.Fprintln(stderr, err)
	}
	return err
}
