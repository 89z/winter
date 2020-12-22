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
select count(1) as count_n, artist_s
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
select count(1) from song_t
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
select count(1) as count_n, album_s from album_t
natural join song_album_t
natural join song_t
where note_s = 'good'
group by album_n
order by count_n;

/*
8|6|The Richest Man in Babylon
8|2|Vegas
9|4|Hold Your Colour
10|0|Selected Ambient Works, Volume II
10|0|Days to Come
10|0|Felt Mountain
10|3|Morning View
12|0|The Elder Scrolls IV: Oblivion
*/
select
   count(1) filter (where note_s = 'good') as good_n,
   count(1) filter (where note_s = '') as unrated_n,
   album_s
from album_t
natural join song_album_t
natural join song_t
group by album_n
order by good_n;

/*
After the Night Falls|2|0
Ambient 2: The Plateaux of Mirror|1|4
Avalon Sutra|0|0
Bandits of Stature|0|3
Before the Day Breaks|3|6
Bordeaux|2|7
By the Dawn's Early Light|1|6
In the Mist|0|5
Jane 1-11|4|5
Jane 12-21|1|0
La Bella Vista|0|1
Lovely Thunder|0|4
Luxa|1|10
Perhaps|1|0
The Pavilion of Dreams|0|0
The Pearl|1|8
The White Arcades|1|0
Through the Hill|2|9
*/
select
   album_s,
   count(1) filter (where note_s = 'good') as good_n,
   count(1) filter (where note_s = '') as unrated_n
from album_t
natural join song_album_t
natural join song_t
natural join song_artist_t
natural join artist_t
where artist_s = 'Harold Budd'
group by album_n;

/*
unrated tracks | good tracks | color
---------------|-------------|------
0              | 0           | red
0              | 1           | green
1              | 0           | light red
1              | 1           | light green
*/
