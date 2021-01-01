package model

type Artist struct {
   ID int
   Name string
   Album []int
}

type Album struct {
   ID int
   Name string
   Song []int
}

type Song struct {
   ID int
   Name string
}
