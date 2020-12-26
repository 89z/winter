select
   count(1) filter (where note_s = 'good') as count_n,
   artist_s
from artist_t
natural join song_artist_t
natural join song_t
where check_s < '2019-12-25'
group by artist_n
order by count_n desc;
