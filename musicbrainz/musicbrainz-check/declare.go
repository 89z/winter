package main
import "winter/snow"

const (
   BLOCK = "\u2587\u2587\u2587\u2587\u2587"
   GREEN_10 = "\x1b[92m" + BLOCK + BLOCK + "\x1b[m"
   GREEN_5 = "\x1b[92m" + BLOCK + "\x1b[90m" + BLOCK + "\x1b[m"
   RED_10 = "\x1b[91m" + BLOCK + BLOCK + "\x1b[m"
   RED_5 = "\x1b[91m" + BLOCK + "\x1b[90m" + BLOCK + "\x1b[m"
)

func Color(url_s string, unrated_n, good_n int) string {
   if snow.Pop(url_s) {
      return GREEN_10
   }
   if unrated_n == 0 && good_n == 0 {
      return RED_10
   }
   if unrated_n == 0 {
      return GREEN_10
   }
   if good_n == 0 {
      return RED_5
   }
   return GREEN_5
}

/* Regarding the title and date:

For the title, we will display the remote Group title, but we also need to get
the remote Release titles to match against the local Release titles.

For the date, if we have a local match, use that date. Otherwise, use use the
remote Group date */
type Group struct {
   Color string
   Date string
   Release map[string]bool
   Title string
}
