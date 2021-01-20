package winter
import "strings"

func Pop(s string) bool {
   return strings.HasPrefix(s, "youtube.com/watch?")
}
