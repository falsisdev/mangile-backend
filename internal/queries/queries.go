// Package queries, harici API'lerde kullanilan GROQ ve GraphQL sorgularini
// derleme zamani ikili dosyaya gomulu tutar; sorgular kod icinde degil,
// ayrik dosyalarda yasar.
//
// Parametrelendirme kurallari:
//   - GROQ: kullanici kaynakli degerler asla sorgu metnine serpistirilmez;
//     $param olarak gonderilir (bkz. services.sanityClient.query).
//   - Yalnizca GROQ'un parametrelendiremedigi dinamik parcalar icin
//     (siralama alani gibi) REPLACE_* placeholder'lari kullanilir ve servis
//     katmaninda beyaz listeyle degistirilir.
package queries

import _ "embed"

// ---------------------------------------------------------------
// GROQ
// ---------------------------------------------------------------

// TitleQuery, /api/titles ucunda bir seriyi (manga/lightNovel) ayrik
// mangaChapter/novelChapter dokumanlariyla birlikte dondurur.
//
//go:embed title.groq
var TitleQuery string

// ChapterQuery, /api/chapter ucunda tek bir bölümü ve ayni eserin tum
// bolumlerini dondurur.
//
//go:embed chapter.groq
var ChapterQuery string

// ListQuery, /api/list/:id ucunda tek bir listeyi sahibi ve begenilerle
// birlikte dondurur.
//
//go:embed list.groq
var ListQuery string

// UserQuery, /api/user/:id ucunda bir kullaniciyi ve listelerini dondurur.
//
//go:embed user.groq
var UserQuery string

// ArticleQuery, /api/article/:slug ucunda slug'a gore tek bir makale
// dondurur.
//
//go:embed article.groq
var ArticleQuery string

// ScanQuery, /api/scan/:id ucunda bir ceviri ekibini uyeleriyle birlikte
// dondurur.
//
//go:embed scan.groq
var ScanQuery string

// SanityListQuery, /api/sanityList ucunda tum yerel serileri dondurur.
// REPLACE_SORT_FIELD alani servis katmaninda beyaz listeden gelir.
//
//go:embed sanityList.groq
var SanityListQuery string

// MangaListQuery, /api/mangaList ucunda AniList sonuclarini yerel icerikle
// zenginlestirmek icin kullanilir.
//
//go:embed mangaList.groq
var MangaListQuery string

// ---------------------------------------------------------------
// GraphQL
// ---------------------------------------------------------------

// AniListMediaDetails, /api/titles zenginlestirmesi icin tek bir manga
// medyasini detaylariyla dondurur.
//
//go:embed graphql/anilist_media.graphql
var AniListMediaDetails string

// AniListMangaList, /api/mangaList ucunda AniList katalogunu sayfali dondurur.
//
//go:embed graphql/anilist_manga_list.graphql
var AniListMangaList string
