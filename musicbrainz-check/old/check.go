package main

import (
   "fmt"
   "sort"
   "strings"
)

/* Regarding the title and date:

For the title, we will display the remote Group title, but we also need to get
the remote Release titles to match against the local Release title.

For the date, if we have a local match, use that date. Otherwise, use use the
remote Group date */
type winterRemote struct {
   color, date, title string
   release map[string]bool
}

func color(url string, unrated, good int) string {
   const (
      block = "\u2587\u2587\u2587\u2587\u2587"
      greenFive = "\x1b[92m" + block + "\x1b[90m" + block + "\x1b[0m"
      greenTen = "\x1b[92m" + block + block + "\x1b[0m"
      redFive = "\x1b[91m" + block + "\x1b[90m" + block + "\x1b[0m"
      redTen = "\x1b[91m" + block + block + "\x1b[0m"
   )
   if strings.HasPrefix(url, "youtube.com/watch?") {
      return greenTen
   }
   if unrated == 0 && good == 0 {
      return redTen
   }
   if unrated == 0 {
      return greenTen
   }
   if good == 0 {
      return redFive
   }
   return greenFive
}

func main() {
   for n, group := range remotes {
      for release := range group.release {
         local, ok := locals[strings.ToUpper(release)]
         if ok {
            remotes[n].date = local.date
            remotes[n].color = local.color
         }
      }
   }
   sort.Slice(remotes, func (d, e int) bool {
      return remotes[d].date < remotes[e].date
   })
   for _, group := range remotes {
      fmt.Printf("%-10v | %10v | %v\n", group.date, group.color, group.title)
   }
}
