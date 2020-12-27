# Database

## CSV import

~~~
sqlite3 music.db
.read music.sql
.mode csv
.import --skip 1 artist_t.csv artist_t
.import --skip 1 album_t.csv album_t
.import --skip 1 song_artist_t.csv song_artist_t
.import --skip 1 song_t.csv song_t
.quit
~~~

<https://sqlite.org/cli.html#importing_csv_files>

## CSV export

~~~
sqlite3 music.db
.mode csv
.once artist_t.csv
select * from artist_t;
.once album_t.csv
select * from album_t;
.once song_t.csv
select * from song_t;
.once song_artist_t.csv
select * from song_artist_t;
~~~

- <https://github.com/mattn/go-sqlite3/tree/master/_example>
- <https://golang.org/pkg/database/sql#DB.Query>
- <https://sqlite.org/cli.html#export_to_csv>
