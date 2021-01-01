package dbstorm

type Artist struct {
   Name string
   MB string
   Check string
   ID int `storm:"increment"` // ARTIST ID
   Album []int // NEW ALBUM IDS
}

type Album struct {
   Name string
   Date string
   URL string
   ID int `storm:"increment"` // NEW ALBUM ID
   Song []int // NEW SONG IDS
}

type Song struct {
   Name string
   Note string
   ID int `storm:"increment"` // NEW SONG ID
}
