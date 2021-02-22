package main

import (
   "log"
   "os"
)

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   name, file := os.Args[1], os.Getenv("WINTER")
   local, e := newLocalArtist(name, file)
   if e != nil {
      log.Fatal(e)
   }
   remote, e := remoteAlbums(local.id)
   if e != nil {
      log.Fatal(e)
   }
   log.Print(remote)
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
