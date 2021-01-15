package main

import (
   "io"
   "os"
   "os/exec"
)

func Less() (*exec.Cmd, io.WriteCloser, error) {
   less := exec.Command("less")
   pipe, e := less.StdinPipe()
   if e != nil {
      return nil, nil, e
   }
   less.Stdout = os.Stdout
   return less, pipe, less.Start()
}
