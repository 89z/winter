package main

import (
   "io"
   "os"
   "os/exec"
)

func less() (*exec.Cmd, io.WriteCloser, error) {
   cmd := exec.Command("less")
   pipe, e := cmd.StdinPipe()
   if e != nil {
      return nil, nil, e
   }
   cmd.Stdout = os.Stdout
   return cmd, pipe, cmd.Start()
}
