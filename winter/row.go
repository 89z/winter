package main

const DASH = "-----------------------------------------------------------------"
const SPACE = "                                                                "
const WIDTH = 48
const YELLOW = "\x1b[43m   \x1b[m"

type Row struct {
   AlbumInt int
   AlbumStr string
   DateStr string
   NoteStr string
   SongInt int
   SongStr string
   UrlStr string
}

func Pop(s string) bool {
   return s[:18] == "youtube.com/watch?"
}
