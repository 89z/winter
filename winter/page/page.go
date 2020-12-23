package main

import (
   "fmt"
   "github.com/tigrawap/slit"
   "time"
)

const bonobo = `
artist_n | 58
check_s  | 2020-12-22
mb_s     | 9a709693-b4f8-4da9-8cc1-038c911a61be

album_n | 195
album_s | Animal Magic
date_s  | 2000-07-24
--------+----------------------------------------------+-------
song_n  | song_s                                       | note_s
--------+----------------------------------------------+-------
   2216 | Gypsy                                        | bad
   2217 | Sugar Rhyme                                  | bad
   2218 | Intro                                        | short
   2219 | Dinosaurs                                    | bad
   2220 | Terrapin                                     | good
   2221 | Shadowtricks                                 | good
   2222 | Sleepy Seven                                 | good
   2223 | Kota                                         | good
   2224 | The Plug                                     | good
   2225 | Silver                                       | good

album_n | 194
album_s | One Offs… Remixes & B Sides
date_s  | 2002-08-12
--------+----------------------------------------------+-------
song_n  | song_s                                       | note_s
--------+----------------------------------------------+-------
   2211 | Scuba                                        | good
   2212 | Dismantling Frank                            | good
   2213 | The Sicilian                                 | good
   2214 | The Shark                                    | good
   2215 | Magicman                                     | good

album_n | 193
album_s | Dial 'M' for Monkey
date_s  | 2003-06-09
--------+----------------------------------------------+-------
song_n  | song_s                                       | note_s
--------+----------------------------------------------+-------
   2202 | Noctuary                                     | bad
   2203 | D Song                                       | good
   2204 | Change Down                                  | good
   2205 | Wayward Bob                                  | good
   2206 | Pick Up                                      | good
   2207 | Something for Windy                          | short
   2208 | Nothing Owed                                 | good
   2209 | Light Pattern                                | good
   2210 | Flutter                                      | good

album_n | 192
album_s | Days to Come
date_s  | 2006-10-02
--------+----------------------------------------------+-------
song_n  | song_s                                       | note_s
--------+----------------------------------------------+-------
   2184 | Recurring                                    | good
   2185 | Nightlite (demo version)                     | bad
   2186 | Intro                                        | short
   2187 | Walk in the Sky                              | bad
   2188 | Ketto                                        | good
   2189 | On Your Marks                                | good
   2190 | Between the Lines                            | bad
   2191 | Transmission 94, Parts 1 & 2                 | good
   2192 | If You Stayed Over                           | bad
   2193 | If You Stayed Over (reprise)                 | good
   2194 | Days to Come                                 | good
   2195 | Nightlite                                    | good
   2196 | Between the Lines (instrumental)             | bad
   2197 | If You Stayed Over (instrumental)            | bad
   2198 | Walk in the Sky (instrumental)               | good
   2199 | Hatoa                                        | good
   2200 | The Fever                                    | good
   2201 | Days to Come (instrumental)                  | bad

album_n | 198
album_s | Black Sands
date_s  | 2010-03-29
--------+----------------------------------------------+-------
song_n  | song_s                                       | note_s
--------+----------------------------------------------+-------
   2251 | El Toro                                      | good
   2252 | We Could Forever                             | good
   2253 | All in Forms                                 | bad
   2254 | The Keeper                                   | short
   2255 | Stay the Same                                | good
   2256 | Animals                                      | good
   2257 | Prelude                                      | short
   2258 | Kiara                                        | good
   2259 | Kong                                         | bad
   2260 | Eyesdown                                     | good
   2261 | 1009                                         | good
   2262 | Black Sands                                  | good

album_n | 197
album_s | The North Borders
date_s  | 2013-03-21
--------+----------------------------------------------+-------
song_n  | song_s                                       | note_s
--------+----------------------------------------------+-------
   2238 | Jets                                         | bad
   2239 | Towers                                       | bad
   2240 | Don’t Wait                                   | bad
   2241 | Antenna                                      | bad
   2242 | Emkay                                        | bad
   2243 | Heaven for the Sinner                        | bad
   2244 | Sapphire                                     | good
   2245 | Know You                                     | bad
   2246 | Ten Tigers                                   | bad
   2247 | Pieces                                       | bad
   2248 | Transits                                     | good
   2249 | First Fires                                  | bad
   2250 | Cirrus                                       | good

album_n | 1202
album_s | Migration
date_s  | 2017-01-13
--------+----------------------------------------------+-------
song_n  | song_s                                       | note_s
--------+----------------------------------------------+-------
  13770 | Migration                                    | bad
  13771 | Break Apart                                  | bad
  13772 | Outlier                                      | bad
  13773 | Grains                                       | bad
  13774 | Second Sun                                   | bad
  13775 | Surface                                      | bad
  13776 | Bambro Koyo Ganda                            | good
  13777 | Kerala                                       | bad
  13778 | Ontario                                      | good
  13779 | No Reason                                    |
  13780 | 7th Sevens                                   |
  13781 | Figures                                      |
`

func main() {
   ch := make(chan string)
   s, err := slit.NewFromStream(ch)
   if err != nil {
      panic(err)
   }
   go func() {
      for i := 0; i < 100; i++ {
         ch <- fmt.Sprintf("line %d", i)
         time.Sleep(10 * time.Millisecond)
      }
      close(ch)
   }()
   s.Display()
   s.Shutdown()
}
