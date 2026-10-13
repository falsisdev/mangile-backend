package queries

import _ "embed"

//go:embed title.groq
var TitleQuery string

//go:embed chapter.groq
var ChapterQuery string

//go:embed list.groq
var ListQuery string

//go:embed publicLists.groq
var PublicListsQuery string

//go:embed userLists.groq
var UserListsQuery string

//go:embed likedLists.groq
var LikedListsQuery string

//go:embed user.groq
var UserQuery string

//go:embed history.groq
var HistoryQuery string

//go:embed article.groq
var ArticleQuery string

//go:embed scan.groq
var ScanQuery string

//go:embed latestTitles.groq
var LatestTitlesQuery string

//go:embed titlesByTag.groq
var TitlesByTagQuery string

//go:embed mangaList.groq
var MangaListQuery string

//go:embed latestChapters.groq
var LatestChaptersQuery string

//go:embed localTitles.groq
var LocalTitlesQuery string

//go:embed graphql/anilist_media.graphql
var AniListMediaDetails string

//go:embed graphql/anilist_manga_list.graphql
var AniListMangaList string
