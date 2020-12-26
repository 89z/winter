package main

import (
   "strings"
   "winter/snow"
)

type Row struct {
   AlbumInt int
   AlbumStr string
   DateStr string
   NoteStr string
   SongInt int
   SongStr string
   UrlStr string
}

const DASH = "-----------------------------------------------------------------"
const SPACE = "                                                                "
const WIDTH = 50
const YELLOW = "\x1b[43m   \x1b[m"

func Note(r Row, song_m map[string]int) (string, string) {
   if song_m[strings.ToUpper(r.SongStr)] > 1 && r.NoteStr == "" {
      return "\x1b[30;43m%v\x1b[m", "duplicate"
   }
   if r.NoteStr == "" && ! snow.Pop(r.UrlStr) {
      return YELLOW + "%6v", ""
   }
   return "%-9v", r.NoteStr
}

func Less() (exec.Cmd, io.Writer, error) {
   less := exec.Command("less")
   pipe, e := less.StdinPipe()
   if e != nil {
      return nil, nil, e
   }
   less.Stdout = os.Stdout
   return less, pipe, less.Start()
}
