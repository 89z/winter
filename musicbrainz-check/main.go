package main

import (
   "fmt"
   "os"
)

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      return
   }
   name, file := os.Args[1], os.Getenv("WINTER")
   local, err := newLocalArtist(name, file)
   if err != nil {
      panic(err)
   }
   remote, err := remoteAlbums(local.id)
   if err != nil {
      panic(err)
   }
   fmt.Println(remote)
   /*
   index, ok := remote[release.Group.Id]
   if ok {
      // add release to group
      remotes[index].release[release.Title] = true
   } else {
      // add group
      remotes = append(remotes, winterRemote{
         date: release.Group.FirstRelease,
         release: map[string]bool{release.Title: true},
         title: release.Group.Title,
      })
      remote[release.Group.Id] = len(remotes) - 1
   }
   */
}
