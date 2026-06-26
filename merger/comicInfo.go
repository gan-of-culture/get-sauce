package merger

import (
	"encoding/xml"

	"github.com/gan-of-culture/get-sauce/static"
	"github.com/pkg/errors"
)

// generated from https://github.com/anansi-project/comicinfo/blob/main/schema/v2.0/ComicInfo.xsd

// YesNo defines the schema's YesNo restriction base.
type YesNo string

const (
	YesNoUnknown YesNo = "Unknown"
	YesNoNo      YesNo = "No"
	YesNoYes     YesNo = "Yes"
)

// Manga defines the schema's Manga restriction base.
type Manga string

const (
	MangaUnknown           Manga = "Unknown"
	MangaNo                Manga = "No"
	MangaYes               Manga = "Yes"
	MangaYesAndRightToLeft Manga = "YesAndRightToLeft"
)

// AgeRating defines the schema's AgeRating restriction base.
type AgeRating string

const (
	AgeRatingUnknown        AgeRating = "Unknown"
	AgeRatingAdultsOnly18   AgeRating = "Adults Only 18+"
	AgeRatingEarlyChildhood AgeRating = "Early Childhood"
	AgeRatingEveryone       AgeRating = "Everyone"
	AgeRatingEveryone10     AgeRating = "Everyone 10+"
	AgeRatingG              AgeRating = "G"
	AgeRatingKidsToAdults   AgeRating = "Kids to Adults"
	AgeRatingM              AgeRating = "M"
	AgeRatingMA15           AgeRating = "MA15+"
	AgeRatingMature17       AgeRating = "Mature 17+"
	AgeRatingPG             AgeRating = "PG"
	AgeRatingR18            AgeRating = "R18+"
	AgeRatingRatingPending  AgeRating = "Rating Pending"
	AgeRatingTeen           AgeRating = "Teen"
	AgeRatingX18            AgeRating = "X18+"
)

// ComicPageType defines the schema's ComicPageType restriction base.
type ComicPageType string

const (
	PageTypeFrontCover    ComicPageType = "FrontCover"
	PageTypeInnerCover    ComicPageType = "InnerCover"
	PageTypeRoundup       ComicPageType = "Roundup"
	PageTypeStory         ComicPageType = "Story"
	PageTypeAdvertisement ComicPageType = "Advertisement"
	PageTypeEditorial     ComicPageType = "Editorial"
	PageTypeLetters       ComicPageType = "Letters"
	PageTypePreview       ComicPageType = "Preview"
	PageTypeBackCover     ComicPageType = "BackCover"
	PageTypeOther         ComicPageType = "Other"
	PageTypeDeleted       ComicPageType = "Deleted"
)

// --- Complex Types ---

// ComicInfo represents the root element and main metadata layout.
type ComicInfo struct {
	XMLName             xml.Name        `xml:"ComicInfo"`
	XmlnsXsi            string          `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation   string          `xml:"xsi:noNamespaceSchemaLocation,attr"`
	Title               string          `xml:"Title,omitempty"`
	Series              string          `xml:"Series,omitempty"`
	Number              string          `xml:"Number,omitempty"`
	Count               int             `xml:"Count,omitempty"`  // Default: -1
	Volume              int             `xml:"Volume,omitempty"` // Default: -1
	AlternateSeries     string          `xml:"AlternateSeries,omitempty"`
	AlternateNumber     string          `xml:"AlternateNumber,omitempty"`
	AlternateCount      int             `xml:"AlternateCount,omitempty"` // Default: -1
	Summary             string          `xml:"Summary,omitempty"`
	Notes               string          `xml:"Notes,omitempty"`
	Year                int             `xml:"Year,omitempty"`  // Default: -1
	Month               int             `xml:"Month,omitempty"` // Default: -1
	Day                 int             `xml:"Day,omitempty"`   // Default: -1
	Writer              string          `xml:"Writer,omitempty"`
	Penciller           string          `xml:"Penciller,omitempty"`
	Inker               string          `xml:"Inker,omitempty"`
	Colorist            string          `xml:"Colorist,omitempty"`
	Letterer            string          `xml:"Letterer,omitempty"`
	CoverArtist         string          `xml:"CoverArtist,omitempty"`
	Editor              string          `xml:"Editor,omitempty"`
	Publisher           string          `xml:"Publisher,omitempty"`
	Imprint             string          `xml:"Imprint,omitempty"`
	Genre               string          `xml:"Genre,omitempty"`
	Web                 string          `xml:"Web,omitempty"`
	PageCount           int             `xml:"PageCount,omitempty"` // Default: 0
	LanguageISO         string          `xml:"LanguageISO,omitempty"`
	Format              string          `xml:"Format,omitempty"`
	BlackAndWhite       YesNo           `xml:"BlackAndWhite,omitempty"` // Default: Unknown
	Manga               Manga           `xml:"Manga,omitempty"`         // Default: Unknown
	Characters          string          `xml:"Characters,omitempty"`
	Teams               string          `xml:"Teams,omitempty"`
	Locations           string          `xml:"Locations,omitempty"`
	ScanInformation     string          `xml:"ScanInformation,omitempty"`
	StoryArc            string          `xml:"StoryArc,omitempty"`
	SeriesGroup         string          `xml:"SeriesGroup,omitempty"`
	AgeRating           AgeRating       `xml:"AgeRating,omitempty"`       // Default: Unknown
	Pages               []ComicPageInfo `xml:"Pages>Page,omitempty"`      // Handles ArrayOfComicPageInfo mapping
	CommunityRating     float64         `xml:"CommunityRating,omitempty"` // xs:decimal (0.00 to 5.00)
	MainCharacterOrTeam string          `xml:"MainCharacterOrTeam,omitempty"`
	Review              string          `xml:"Review,omitempty"`
}

// ComicPageInfo represents the attributes associated with an individual page.
type ComicPageInfo struct {
	Image       int           `xml:"Image,attr"` // Required
	Type        ComicPageType `xml:"Type,attr,omitempty"`
	DoublePage  bool          `xml:"DoublePage,attr,omitempty"`
	ImageSize   int64         `xml:"ImageSize,attr,omitempty"`
	Key         string        `xml:"Key,attr,omitempty"`
	Bookmark    string        `xml:"Bookmark,attr,omitempty"`
	ImageWidth  int           `xml:"ImageWidth,attr,omitempty"`
	ImageHeight int           `xml:"ImageHeight,attr,omitempty"`
}

func NewComicInfo(data *static.Data) ([]byte, error) {
	comicInfo := ComicInfo{}
	comicInfo.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	comicInfo.XsiSchemaLocation = "https://raw.githubusercontent.com/anansi-project/comicinfo/main/schema/v2.0/ComicInfo.xsd"

	comicInfo.Title = data.Title
	comicInfo.Web = data.URL
	comicInfo.Format = "Digital"
	comicInfo.Manga = MangaYesAndRightToLeft

	stream, ok := data.Streams["0"]
	if ok {
		comicInfo.PageCount = len(stream.URLs)
		for i := range stream.URLs {
			if i == 0 {
				comicInfo.Pages = append(comicInfo.Pages, ComicPageInfo{Image: i, Type: PageTypeFrontCover})
				continue
			}
			comicInfo.Pages = append(comicInfo.Pages, ComicPageInfo{Image: i, Type: PageTypeStory})
		}
	}

	xmlData, err := xml.MarshalIndent(comicInfo, "", "    ")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return append([]byte(xml.Header), xmlData...), nil
}
