package main

const DASH = "-----------------------------------------------------------------"
const SPACE = "                                                                "
const WIDTH = 46
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
