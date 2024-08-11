package mpegdash

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type AdaptationSet struct {
	Text               string `xml:",chardata"`
	ID                 string `xml:"id,attr"`
	ContentType        string `xml:"contentType,attr"`
	StartWithSAP       string `xml:"startWithSAP,attr"`
	SegmentAlignment   string `xml:"segmentAlignment,attr"`
	BitstreamSwitching string `xml:"bitstreamSwitching,attr"`
	FrameRate          string `xml:"frameRate,attr"`
	MaxWidth           string `xml:"maxWidth,attr"`
	MaxHeight          string `xml:"maxHeight,attr"`
	Par                string `xml:"par,attr"`
	Lang               string `xml:"lang,attr"`
	Representation     struct {
		Text              string `xml:",chardata"`
		ID                string `xml:"id,attr"`
		MimeType          string `xml:"mimeType,attr"`
		Codecs            string `xml:"codecs,attr"`
		Bandwidth         string `xml:"bandwidth,attr"`
		Width             string `xml:"width,attr"`
		Height            string `xml:"height,attr"`
		Sar               string `xml:"sar,attr"`
		AudioSamplingRate string `xml:"audioSamplingRate,attr"`
		SegmentTemplate   struct {
			Text            string `xml:",chardata"`
			Timescale       string `xml:"timescale,attr"`
			Initialization  string `xml:"initialization,attr"`
			Media           string `xml:"media,attr"`
			StartNumber     string `xml:"startNumber,attr"`
			SegmentTimeline struct {
				Text string `xml:",chardata"`
				S    []struct {
					Text string `xml:",chardata"`
					T    string `xml:"t,attr"`
					D    string `xml:"d,attr"`
					R    string `xml:"r,attr"`
				} `xml:"S"`
			} `xml:"SegmentTimeline"`
		} `xml:"SegmentTemplate"`
		AudioChannelConfiguration struct {
			Text        string `xml:",chardata"`
			SchemeIdUri string `xml:"schemeIdUri,attr"`
			Value       string `xml:"value,attr"`
		} `xml:"AudioChannelConfiguration"`
	} `xml:"Representation"`
}

type MPD struct {
	XMLName                   xml.Name `xml:"MPD"`
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xsi,attr"`
	Xmlns                     string   `xml:"xmlns,attr"`
	Xlink                     string   `xml:"xlink,attr"`
	SchemaLocation            string   `xml:"schemaLocation,attr"`
	Profiles                  string   `xml:"profiles,attr"`
	Type                      string   `xml:"type,attr"`
	MediaPresentationDuration string   `xml:"mediaPresentationDuration,attr"`
	MaxSegmentDuration        string   `xml:"maxSegmentDuration,attr"`
	MinBufferTime             string   `xml:"minBufferTime,attr"`
	ProgramInformation        string   `xml:"ProgramInformation"`
	ServiceDescription        struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
	} `xml:"ServiceDescription"`
	Period struct {
		Text          string `xml:",chardata"`
		ID            string `xml:"id,attr"`
		Start         string `xml:"start,attr"`
		AdaptationSet []AdaptationSet
	} `xml:"Period"`
}

var reIdentifier = regexp.MustCompile(`\$[^\$]+\$`)
var reFormatTag = regexp.MustCompile(`%[ 0\-+#]\d+d`)

func ExtractDASHManifest(URL string, headers map[string]string) (map[string]*static.Stream, error) {

	manifest, err := request.GetWithHeaders(URL, headers)
	if err != nil {
		return nil, err
	}

	return ParseDASHManifest(&manifest, URL)

}

// ParseDASHManifest from XML content
func ParseDASHManifest(xmlString *string, URL string) (map[string]*static.Stream, error) {
	mpd := MPD{}
	err := xml.Unmarshal([]byte(*xmlString), &mpd)
	if err != nil {
		return nil, err
	}

	out := map[string]*static.Stream{}
	for idx, aSet := range mpd.Period.AdaptationSet {
		out[fmt.Sprint(idx)], err = parseAdaptionSet(&aSet, URL)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func replaceIdentifier(URI string, vars map[string]string) (string, error) {
	for _, iName := range reIdentifier.FindAllString(URI, -1) {
		iNameTrimmed := strings.Trim(iName, "$")
		if formatTag := reFormatTag.FindString(iNameTrimmed); formatTag != "" {
			iNameTrimmed = strings.Replace(iNameTrimmed, formatTag, "", 1)
			number, err := strconv.Atoi(vars[iNameTrimmed])
			if err != nil {
				return "", err
			}
			vars[iNameTrimmed] = fmt.Sprintf(formatTag, number)
		}

		URI = strings.ReplaceAll(URI, iName, vars[iNameTrimmed])
	}
	return URI, nil
}

func parseAdaptionSet(aSet *AdaptationSet, URL string) (*static.Stream, error) {

	initialSegmentURL, err := replaceIdentifier(aSet.Representation.SegmentTemplate.Initialization, map[string]string{"RepresentationID": aSet.Representation.ID})
	if err != nil {
		return nil, err
	}

	URLs := []*static.URL{
		{
			URL: initialSegmentURL,
			Ext: utils.GetFileExt(initialSegmentURL),
		},
	}
	numberOfSegments := 0
	for _, segment := range aSet.Representation.SegmentTemplate.SegmentTimeline.S {
		repetitions, _ := strconv.Atoi(segment.R)
		if repetitions > 0 {
			numberOfSegments += repetitions
		}
		numberOfSegments += 1
	}
	startNumber, _ := strconv.Atoi(aSet.Representation.SegmentTemplate.StartNumber)
	for i := startNumber; i <= numberOfSegments; i++ {
		segmentURI, err := replaceIdentifier(aSet.Representation.SegmentTemplate.Media, map[string]string{"RepresentationID": aSet.Representation.ID, "Number": fmt.Sprint(i)})
		if err != nil {
			return nil, err
		}

		URLs = append(URLs, &static.URL{
			URL: segmentURI,
			Ext: utils.GetFileExt(segmentURI),
		})
	}

	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	for _, URL := range URLs {
		URL.URL = strings.ReplaceAll(URL.URL, `\`, "/")
		if !strings.Contains(URL.URL, "http") {
			segmentURL, err := baseURL.Parse(URL.URL)
			if err != nil {
				return nil, err
			}
			URL.URL = segmentURL.String()
		}
	}

	stream := &static.Stream{
		Type: static.DataType(aSet.ContentType),
		URLs: URLs,
		Info: aSet.Representation.Codecs,
		Ext:  utils.GetLastItemString(strings.Split(aSet.Representation.MimeType, "/")),
	}

	if stream.Type == static.DataTypeVideo {
		stream.Quality = fmt.Sprintf("%sx%s", aSet.Representation.Width, aSet.Representation.Height)
	}

	if aSet.Representation.Codecs == "mp4a.40.2" {
		stream.Ext = "aac"
	}

	return stream, nil
}
