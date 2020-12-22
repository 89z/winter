package main
import "strings"

const WIDTH = 48
const YELLOW = "\x1b[43m   \x1b[m"

type Row struct {
   Album int
   AlbumStr string
   Artist int
   Check string
   Date string
   Note string
   Song int
   SongStr string
   URL string
}

func Pop(s string) bool {
   return strings.HasPrefix(s, "youtube.com/watch?")
}
