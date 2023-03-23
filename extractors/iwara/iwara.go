package iwara

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type File struct {
	ID            string      `json:"id"`
	Type          string      `json:"type"`
	Path          string      `json:"path"`
	Name          string      `json:"name"`
	Mime          string      `json:"mime"`
	Size          int         `json:"size"`
	Width         interface{} `json:"width"`
	Height        interface{} `json:"height"`
	Duration      interface{} `json:"duration"`
	NumThumbnails int         `json:"numThumbnails"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
}

type MediaInfo struct {
	ID              string      `json:"id"`
	Slug            string      `json:"slug"`
	Title           string      `json:"title"`
	Body            string      `json:"body"`
	Status          string      `json:"status"`
	Rating          string      `json:"rating"`
	Private         bool        `json:"private"`
	Unlisted        bool        `json:"unlisted"`
	Thumbnail       interface{} `json:"thumbnail"`
	EmbedURL        interface{} `json:"embedUrl"`
	Liked           bool        `json:"liked"`
	NumLikes        int         `json:"numLikes"`
	NumViews        int         `json:"numViews"`
	NumComments     int         `json:"numComments"`
	File            File        `json:"file"`
	CustomThumbnail interface{} `json:"customThumbnail"`
	User            struct {
		ID         string    `json:"id"`
		Name       string    `json:"name"`
		Username   string    `json:"username"`
		Status     string    `json:"status"`
		Role       string    `json:"role"`
		FollowedBy bool      `json:"followedBy"`
		Following  bool      `json:"following"`
		Friend     bool      `json:"friend"`
		Premium    bool      `json:"premium"`
		SeenAt     time.Time `json:"seenAt"`
		Avatar     struct {
			ID            string      `json:"id"`
			Type          string      `json:"type"`
			Path          string      `json:"path"`
			Name          string      `json:"name"`
			Mime          string      `json:"mime"`
			Size          int         `json:"size"`
			Width         interface{} `json:"width"`
			Height        interface{} `json:"height"`
			Duration      interface{} `json:"duration"`
			NumThumbnails interface{} `json:"numThumbnails"`
			CreatedAt     time.Time   `json:"createdAt"`
			UpdatedAt     time.Time   `json:"updatedAt"`
		} `json:"avatar"`
		CreatedAt time.Time   `json:"createdAt"`
		UpdatedAt time.Time   `json:"updatedAt"`
		DeletedAt interface{} `json:"deletedAt"`
	} `json:"user"`
	Tags []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"tags"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	DeletedAt interface{} `json:"deletedAt"`
	Files     []File      `json:"files"`
	FileURL   string      `json:"fileUrl"`
}

type VideoSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Src  struct {
		View     string `json:"view"`
		Download string `json:"download"`
	} `json:"src"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Type      string    `json:"type"`
}

type SearchResult struct {
	Count   int `json:"count"`
	Limit   int `json:"limit"`
	Page    int `json:"page"`
	Results []struct {
		ID        string      `json:"id"`
		Slug      string      `json:"slug"`
		Title     string      `json:"title"`
		Body      interface{} `json:"body"`
		Thumbnail struct {
			ID            string      `json:"id"`
			Type          string      `json:"type"`
			Path          string      `json:"path"`
			Name          string      `json:"name"`
			Mime          string      `json:"mime"`
			Size          int         `json:"size"`
			Width         int         `json:"width"`
			Height        int         `json:"height"`
			Duration      interface{} `json:"duration"`
			NumThumbnails interface{} `json:"numThumbnails"`
			CreatedAt     time.Time   `json:"createdAt"`
			UpdatedAt     time.Time   `json:"updatedAt"`
		} `json:"thumbnail"`
		Rating      string        `json:"rating"`
		Liked       bool          `json:"liked"`
		NumImages   int           `json:"numImages"`
		NumLikes    int           `json:"numLikes"`
		NumViews    int           `json:"numViews"`
		NumComments int           `json:"numComments"`
		File        File          `json:"file"`
		CreatedAt   time.Time     `json:"createdAt"`
		UpdatedAt   time.Time     `json:"updatedAt"`
		DeletedAt   interface{}   `json:"deletedAt"`
		Files       []interface{} `json:"files"`
		Tags        []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"tags"`
		User struct {
			ID         string    `json:"id"`
			Name       string    `json:"name"`
			Username   string    `json:"username"`
			Status     string    `json:"status"`
			Role       string    `json:"role"`
			FollowedBy bool      `json:"followedBy"`
			Following  bool      `json:"following"`
			Friend     bool      `json:"friend"`
			Premium    bool      `json:"premium"`
			SeenAt     time.Time `json:"seenAt"`
			Avatar     struct {
				ID            string      `json:"id"`
				Type          string      `json:"type"`
				Path          string      `json:"path"`
				Name          string      `json:"name"`
				Mime          string      `json:"mime"`
				Size          int         `json:"size"`
				Width         int         `json:"width"`
				Height        int         `json:"height"`
				Duration      interface{} `json:"duration"`
				NumThumbnails interface{} `json:"numThumbnails"`
				CreatedAt     time.Time   `json:"createdAt"`
				UpdatedAt     time.Time   `json:"updatedAt"`
			} `json:"avatar"`
			CreatedAt time.Time   `json:"createdAt"`
			UpdatedAt time.Time   `json:"updatedAt"`
			DeletedAt interface{} `json:"deletedAt"`
		} `json:"user"`
	} `json:"results"`
}

const site = "https://iwara.tv/"
const api = "https://api.iwara.tv/"
const files = "https://files.iwara.tv/"

type extractor struct{}

// New returns a iwara extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	postIDs := parseURL(URL)
	if len(postIDs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, pID := range postIDs {
		d, err := extractData(pID)
		if err != nil {
			return nil, utils.Wrap(err, pID)
		}
		data = append(data, d...)
	}

	return data, nil
}

func parseURL(URL string) []string {
	if ok, _ := regexp.MatchString(site+`(?:video|image)/`, URL); ok {
		return []string{URL}
	}

	tmpURL := regexp.MustCompile(`page=\d+`).ReplaceAllString(URL, "page=%d")
	if !strings.Contains(tmpURL, "page=%d") {
		tmpURL = URL + "&page=%d"
	}

	tmpURL = strings.Replace(tmpURL, "https://", "https://api.", 1)

	out := []string{}
	count := 0
	for i := 0; ; {
		res, err := request.GetAsBytes(fmt.Sprintf(tmpURL, i))
		if err != nil {
			return nil
		}
		if config.Amount > 0 {
			fmt.Println(count)
		}

		searchResult := SearchResult{}
		err = json.Unmarshal(res, &searchResult)
		if err != nil {
			return nil
		}

		URLs := []string{}
		for _, result := range searchResult.Results {
			mediaType := "image"
			if result.File.Type == "video" {
				mediaType = result.File.Type
			}
			URLs = append(URLs, fmt.Sprintf("%s%s/%s/%s", site, mediaType, result.ID, result.Slug))
		}
		count += len(URLs)
		i += 1
		out = append(out, URLs...)
		if config.Amount == 0 || count >= config.Amount || len(URLs) == 0 {
			break
		}
	}

	if config.Amount > 0 && len(out) > config.Amount {
		out = out[:config.Amount]
	}

	return out
}

func extractData(URL string) ([]*static.Data, error) {
	splitURL := strings.Split(URL, "/")
	if len(splitURL) < 5 {
		return nil, static.ErrURLParseFailed
	}

	mediaType := splitURL[3]
	id := splitURL[4]

	resJson, err := request.GetAsBytesWithHeaders(fmt.Sprintf("%s%s/%s", api, mediaType, id), map[string]string{
		"Referer": site,
	})
	if err != nil {
		return nil, err
	}

	mediaInfo := MediaInfo{}
	err = json.Unmarshal(resJson, &mediaInfo)
	if err != nil {
		return nil, err
	}

	if mediaInfo.File.ID != "" {
		mediaInfo.Files = append(mediaInfo.Files, mediaInfo.File)
	}

	streams := map[string]*static.Stream{}
	mediaStreams := map[string]*static.Stream{}
	for idx, file := range mediaInfo.Files {
		switch file.Type {
		case "image":
			mediaStreams, err = imageFileInfoToStream(&file, idx)
		case "video":
			mediaStreams, err = videoFileInfoToStream(&file, mediaInfo.FileURL, idx)
		}
		if err != nil {
			return nil, err
		}
		for k, v := range mediaStreams {
			streams[k] = v
		}
	}

	return []*static.Data{
		{
			Site:    site,
			Title:   mediaInfo.Title,
			Type:    static.DataType(mediaType),
			Streams: streams,
			URL:     URL,
		},
	}, nil
}

func imageFileInfoToStream(fileInfo *File, idx int) (map[string]*static.Stream, error) {
	fileExt := utils.GetFileExt(fileInfo.Name)

	quality := ""
	if fileInfo.Height != nil && fileInfo.Width != nil {
		quality = fmt.Sprintf("%vx%v", fileInfo.Width, fileInfo.Height)
	}

	return map[string]*static.Stream{
		fmt.Sprint(idx): {
			Type: utils.GetMediaType(fileExt),
			URLs: []*static.URL{
				{
					URL: fmt.Sprintf("%s%s/large/%s/%s", files, fileInfo.Type, fileInfo.ID, fileInfo.Name),
					Ext: fileExt,
				},
			},
			Quality: quality,
			Size:    int64(fileInfo.Size),
		},
	}, nil

}

func videoFileInfoToStream(fileInfo *File, fileURL string, idx int) (map[string]*static.Stream, error) {
	u, err := url.Parse(fileURL)
	if err != nil {
		return nil, err
	}

	// https://www.iwara.tv/main.c0260392c56cd1aad6a4.js contains the hash salt
	// if you search for X-Version it should be above
	h := sha1.New()
	io.WriteString(h, fmt.Sprintf("%s_%s_5nFp9kmbNnHdAFhaqMvt", fileInfo.ID, u.Query().Get("expires")))

	res, err := request.GetAsBytesWithHeaders(fileURL, map[string]string{
		"Referer":   site,
		"X-Version": fmt.Sprintf("%x", h.Sum(nil)),
	})
	if err != nil {
		return nil, err
	}

	videoSources := []VideoSource{}
	err = json.Unmarshal(res, &videoSources)
	if err != nil {
		return nil, err
	}

	fileExt := utils.GetFileExt(fileInfo.Name)

	out := map[string]*static.Stream{}
	for i, source := range videoSources {
		fileSize, _ := request.Size("https:"+source.Src.View, site)
		quality := source.Name
		if _, err := strconv.ParseInt(quality, 10, 64); err == nil {
			quality = quality + "p"
		}

		out[fmt.Sprint(idx+len(videoSources)-i-1)] = &static.Stream{
			Type: utils.GetMediaType(fileExt),
			URLs: []*static.URL{
				{
					URL: "https:" + source.Src.View,
					Ext: fileExt,
				},
			},
			Quality: quality,
			Size:    fileSize,
		}
	}

	return out, nil
}
