package main
import "github.com/asdine/storm/v3"

type Artist struct {
   ID int `storm:"increment"`
   Name string
   Song []int
}

type Song struct {
   ID int `storm:"increment"`
   Name string
   Artist []int
}

var (
   artist Artist
   song Song
)

func Save(db *storm.DB) error {
   // 1
   artist = Artist{Name: "Zero 7", Song: []int{0, 1}}
   e = db.Save(&artist)
   if e != nil {
      return e
   }
   // 2
   artist = Artist{Name: "Sia", Song: []int{0, 2}}
   e = db.Save(&artist)
   if e != nil {
      return e
   }
   // 3
   song = Song{Name: "Destiny", Artist: []int{0, 1}}
   e = db.Save(&song)
   if e != nil {
      return e
   }
   // 4
   song = Song{Name: "Give It Away", Artist: []int{0}}
   e = db.Save(&song)
   if e != nil {
      return e
   }
   // 5
   song = Song{Name: "Elastic Heart", Artist: []int{1}}
   return db.Save(&song)
}
