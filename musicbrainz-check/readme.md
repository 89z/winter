# MusicBrainz check

Do we have example of two groups with same title?

- https://musicbrainz.org/ws/2/release-group/3bb226b7-db0b-3288-bea4-e2fecbe20c20
- https://musicbrainz.org/ws/2/release-group/97375917-d328-3738-a556-9edfb489cdad

Do we have example of two groups with same first release date?

- https://musicbrainz.org/ws/2/release-group/9c76c1fb-e228-3292-884a-d6e18cbd80ef
- https://musicbrainz.org/ws/2/release-group/9c78b909-9949-3d63-890a-ea995b025626

For local primary key:

~~~
album_t (
   album_s
   date_s
)
~~~

For remote primary key:

~~~
Group {
   Title
   FirstRelease
}
~~~
