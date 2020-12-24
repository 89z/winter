package main

const DASH = "-----------------------------------------------------------------"
const SPACE = "                                                                "

const DUPLICATE = "\x1b[30;43mduplicate\x1b[m"
const YELLOW = "\x1b[43m   \x1b[m"

const WIDTH = 46

type Row struct {
   AlbumInt int
   AlbumStr string
   DateStr string
   NoteStr string
   SongInt int
   SongStr string
   UrlStr string
}
