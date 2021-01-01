package dbstorm

import (
   "io"
   "os"
   "os/exec"
)

const DASH = "-----------------------------------------------------------------"
const SPACE = "                                                                "
const WIDTH = 50
const YELLOW = "\x1b[43m   \x1b[m"

func Less() (*exec.Cmd, io.WriteCloser, error) {
   less := exec.Command("less")
   pipe, e := less.StdinPipe()
   if e != nil {
      return nil, nil, e
   }
   less.Stdout = os.Stdout
   return less, pipe, less.Start()
}
