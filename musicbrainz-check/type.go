package main

/* Regarding the title and date:

For the title, we will display the remote Group title, but we also need to get
the remote Release titles to match against the local Release titles.

For the date, if we have a local match, use that date. Otherwise, use use the
remote Group date */
type winterRemote struct{
   color string
   date string
   release map[string]bool
   title string
}

type winterLocal struct{
   color string
   date string
}

type mbRelease struct{
   ReleaseCount int `json:"release-count"`
   Releases []struct{
      Date string
      Group struct{
         FirstRelease string `json:"first-release-date"`
         Id string
         SecondaryTypes []string `json:"secondary-types"`
         Title string
      } `json:"release-group"`
      Title string
   }
}
