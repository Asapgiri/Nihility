package logic

import (
	"nyarrent/dbase"

	"github.com/er-azh/go-animeschedule"
	"time"
)

// copied from as start

// TimetableShow is a single entry in a Timetable
type TimetableShowOriginal struct {
	AirType                 animeschedule.AirType           `json:"AirType"`
	AiringStatus            string            `json:"AiringStatus"`
	Chinese                 bool              `json:"Chinese"`
	DelayedFrom             time.Time         `json:"DelayedFrom"`
	DelayedUntil            time.Time         `json:"DelayedUntil"`
	English                 *string           `json:"English,omitempty"`
	EpisodeDate             time.Time         `json:"EpisodeDate"`
	EpisodeNumber           int               `json:"EpisodeNumber"`
	Episodes                *int              `json:"Episodes,omitempty"`
	ImageVersionRoute       string            `json:"ImageVersionRoute"`
	Japanese                *string           `json:"Japanese,omitempty"`
	LengthMin               *int              `json:"LengthMin,omitempty"`
	Romaji                  *string           `json:"Romaji,omitempty"`
	Route                   string            `json:"Route"`
	Status                  string            `json:"Status"`
	SubtractedEpisodeNumber int               `json:"SubtractedEpisodeNumber"`
	Title                   string            `json:"Title"`
}

// copied from as end

type AnimeSearchPage struct {
    Page        int
    TotalAmount int
    SearchText  string
    Anime       []animeschedule.ShowDetail
    Added       []bool
}

type TimetableShow struct {
    Anime   TimetableShowOriginal
    Added   bool
    Filled  bool
    Aired   bool
}

type AnimeWeek [7]TimetableShow

type AnimeTimetableFilter struct {
    OnlyOnList  bool
    SendBack    bool
    Hash        string
}

type AnimeTimetablePage struct {
    AnimeWeek   []AnimeWeek
    Filter      AnimeTimetableFilter
}

type EpisodeTorrent struct {
    Torrent     dbase.AnimeDownload
    Info        TorrentInfo
    Progress    Progress
    Url         string
}

type Episodes struct {
    Index       int
    Title       string
    Torrents    []EpisodeTorrent
    Nyaa        []dbase.NyaaData
    NyaaText    string
}

type NyaaFilter struct {
    Group       string
    NameParams  string
    Category    string
    SubCategory string
    ResultCount string
    Resolution  string
    EpisodeFmt  string
    SeasonFmt   string
    ForseRefrsh bool
}

type EpisodeFilter struct {
    Nyaa    NyaaFilter
    Hash    string
}

type PageCounter struct {
	Index		int
	Value		int
	Selected 	bool
}

type TorrentsFilter struct {
	Page            int
	EpisodesPerPage int
	PageCounter     []PageCounter
    Hash            string
}

type Anime struct {
    Anime       dbase.Anime
    Episodes    []Episodes
    Filter      EpisodeFilter
}

type DtoAnime struct {
    SearchText  string
    Anime       []Anime
}
