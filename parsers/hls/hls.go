package hls

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

// ParseMaster into static.Stream to prefill the structure
// returns a pre filled structure where URLs[0].URL is the media stream URI
func ParseMaster(master *string) ([]*static.Stream, error) {
	re := regexp.MustCompile(`#EXT-X-STREAM-INF:([^\n]*)\n([^\n]+)`) // 1=PARAMS 2=MEDIAURI
	matchedStreams := re.FindAllStringSubmatch(*master, -1)
	if len(matchedStreams) < 1 {
		return nil, fmt.Errorf("unable to parse any stream in m3u master: %s", *master)
	}

	out := []*static.Stream{}
	for _, stream := range matchedStreams {
		s := &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: strings.TrimSpace(stream[2]),
					Ext: "",
				},
			},
		}

		re = regexp.MustCompile(`[A-Z\-]+=(?:"[^"]*"|[^,]*)`) // PARAMETERNAME=VALUE
		for _, streamParam := range re.FindAllString(stream[1], -1) {

			splitParam := strings.Split(streamParam, "=")
			splitParam[1] = strings.Trim(splitParam[1], `",`)
			switch splitParam[0] {
			case "RESOLUTION":
				s.Quality = splitParam[1]
			case "CODECS":
				s.Info = splitParam[1]
			}
		}

		out = append(out, s)
	}

	// AUDIO
	re = regexp.MustCompile(`#EXT-X-MEDIA:([^\n]*)\n`) // 1=PARAMS
	matchedAudioStream := re.FindStringSubmatch(*master)
	if len(matchedAudioStream) < 2 {
		return out, nil
	}

	params := map[string]string{}
	for param := range strings.SplitSeq(matchedAudioStream[1], ",") {
		splitParam := strings.Split(param, "=")
		params[splitParam[0]] = strings.Trim(splitParam[1], `"`)
	}

	out = append(out, &static.Stream{
		Type: static.DataTypeAudio,
		URLs: []*static.URL{
			{
				URL: params["URI"],
			},
		},
		Info: params["LANGUAGE"],
	})

	return out, nil
}

// ParseMediaStream into URLs and if found it's key
func ParseMediaStream(mediaStr *string, URL string) ([]*static.URL, []byte, error) {
	re := regexp.MustCompile(`\s[^#]+\s`) // 1=segment URI
	matchedSegmentURLs := re.FindAllString(*mediaStr, -1)
	if len(matchedSegmentURLs) == 0 {
		fmt.Println(*mediaStr)
		return nil, nil, errors.New("no segements found")
	}

	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, nil, err
	}

	segments := []*static.URL{}
	segmentURI := ""
	for _, v := range matchedSegmentURLs {
		segmentURI = strings.TrimSpace(v)
		if !strings.Contains(segmentURI, "http") {
			segmentURL, err := baseURL.Parse(segmentURI)
			if err != nil {
				return nil, nil, err
			}
			segmentURI = segmentURL.String()
		}
		segments = append(segments, &static.URL{
			URL: segmentURI,
			Ext: utils.GetFileExt(segmentURI),
		})
	}

	re = regexp.MustCompile(`#EXT-X-KEY:METHOD=([^,]*),URI="([^"]*)`) //1=HASH e.g. AES-128 2=KEYURI
	matchedEncryptionMeta := re.FindStringSubmatch(*mediaStr)
	if len(matchedEncryptionMeta) != 3 {
		return segments, nil, nil
	}

	keyURL := matchedEncryptionMeta[2]
	if !strings.HasPrefix(matchedEncryptionMeta[2], "http") {
		keyURI, err := baseURL.Parse(matchedEncryptionMeta[2])
		if err != nil {
			return nil, nil, err
		}
		keyURL = keyURI.String()
	}

	key, err := request.GetAsBytesWithHeaders(keyURL, map[string]string{
		"Referer": URL,
	})
	if err != nil {
		return nil, nil, err
	}

	return segments, key, nil
}

// Extract contents of a file/URL into the internal stream structure.
// If the playlist contains multiple streams then each stream will be represented as a unique stream internally
func Extract(URL string, headers map[string]string) (map[string]*static.Stream, error) {

	master, err := request.GetWithHeaders(URL, headers)
	if err != nil {
		return nil, err
	}

	mediaStreams, err := ParseMaster(&master)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	canBeSortedBySize := true
	// complete mediaURLs
	for _, stream := range mediaStreams {
		if strings.HasPrefix(stream.URLs[0].URL, "https://") {
			continue
		}

		mediaURL, err := baseURL.Parse(stream.URLs[0].URL)
		if err != nil {
			return nil, err
		}
		stream.URLs[0].URL = mediaURL.String()
	}

	for _, stream := range mediaStreams {
		mediaStr, err := request.GetWithHeaders(stream.URLs[0].URL, headers)
		if err != nil {
			return nil, err
		}

		stream.URLs, stream.Key, err = ParseMediaStream(&mediaStr, stream.URLs[0].URL)
		if err != nil {
			return nil, err
		}

		// approximate stream size - the first part URLS[0] seems to be larger than the others
		// thats why a part from the middle is used. Still some site don't expose content-length
		// in the header so stream.Size will be 0
		parts := len(stream.URLs)
		stream.Size, err = request.Size(stream.URLs[parts/2].URL, headers["Referer"])
		if stream.Size == 0 || err != nil {
			canBeSortedBySize = false
			continue
		}

		stream.Size = stream.Size * int64(parts)
	}

	// convert streams from slice to map on func exit
	streams := make(map[string]*static.Stream, len(mediaStreams))
	defer func() {
		for idx, stream := range mediaStreams {
			streams[fmt.Sprint(idx)] = stream
		}
	}()

	// SORT streams

	if canBeSortedBySize {

		sort.Slice(mediaStreams, func(i, j int) bool {
			return mediaStreams[i].Size > mediaStreams[j].Size
		})
		return streams, nil
	}

	sort.Slice(mediaStreams, func(i, j int) bool {
		resVal := 0
		resValI := 0
		for v := range strings.SplitSeq(mediaStreams[i].Quality, "x") {
			resVal, _ = strconv.Atoi(v)
			resValI += resVal
		}

		resValJ := 0
		for v := range strings.SplitSeq(mediaStreams[j].Quality, "x") {
			resVal, _ = strconv.Atoi(v)
			resValJ += resVal
		}

		return resValI > resValJ
	})

	return streams, nil
}
