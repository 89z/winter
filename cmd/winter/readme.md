# Database

## JSON export

~~~
sqlite3 winter.db
.mode json

.once artist.json
select * from artist_t;

.once album.json
select * from album_t;

.once song.json
select * from song_t;

.once song_artist.json
select * from song_artist_t;
~~~

- <https://github.com/mattn/go-sqlite3/tree/master/_example>
- <https://golang.org/pkg/database/sql#DB.Query>
- <https://sqlite.org/cli.html#export_to_csv>
