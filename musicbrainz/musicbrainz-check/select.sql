/*
129|Red Hot Chili Peppers
135|M83
136|The Black Dog
152|Aphex Twin
175|Madonna
177|Autechre
179|Tori Amos
206|Harold Budd
*/
select count(*) as count_n, artist_s
from song_artist_t
natural join artist_t
group by artist_n
order by count_n;

-- harold budd songs
select * from album_t
natural join song_album_t
natural join song_t
natural join song_artist_t
natural join artist_t
where artist_s = 'Harold Budd';

-- first, lets just try to get a list of good tracks
select * from song_t
where note_s = 'good';

-- now lets try to get a count of good songs
select count(*) from song_t
where note_s = 'good';

/*
8|Sweet Jones
8|Black Sands
8|Yellow House
8|Mezzanine
8|Portishead
8|Kid A
8|The Cosmic Game
8|The Richest Man in Babylon
8|Vegas
9|Hold Your Colour
10|Selected Ambient Works, Volume II
10|Days to Come
10|Felt Mountain
10|Morning View
12|The Elder Scrolls IV: Oblivion
*/
select count(*) as count_n, album_s from album_t
natural join song_album_t
natural join song_t
where note_s = 'good'
group by album_n
order by count_n;

/*
unrated tracks | good tracks | color
---------------|-------------|------
0              | 0           | red
0              | 1           | green
1              | 0           | light red
1              | 1           | light green
*/
