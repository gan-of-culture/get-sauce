package downloader

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/rain/torrent"
	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/schollz/progressbar/v2"
)

// DefaultConfig for Session. Do not pass zero value Config to NewSession. Copy this struct and modify instead.
var DefaultConfig = torrent.Config{
	// Session
	Database:                               "~/torrent/session.db",
	DataDir:                                "~/torrent/data",
	DataDirIncludesTorrentID:               true,
	PortBegin:                              50000,
	PortEnd:                                60000,
	MaxOpenFiles:                           10240,
	PEXEnabled:                             true,
	ResumeWriteInterval:                    30 * time.Second,
	PrivatePeerIDPrefix:                    "-RN" + "v" + "-",
	PrivateExtensionHandshakeClientVersion: "Rain " + "v",
	BlocklistUpdateInterval:                24 * time.Hour,
	BlocklistUpdateTimeout:                 10 * time.Minute,
	BlocklistEnabledForTrackers:            true,
	BlocklistEnabledForOutgoingConnections: true,
	BlocklistEnabledForIncomingConnections: true,
	BlocklistMaxResponseSize:               100 << 20,
	TorrentAddHTTPTimeout:                  30 * time.Second,
	MaxMetadataSize:                        10 << 20,
	MaxTorrentSize:                         10 << 20,
	MaxPieces:                              64 << 10,
	DNSResolveTimeout:                      5 * time.Second,

	// RPC Server
	RPCEnabled:         true,
	RPCHost:            "127.0.0.1",
	RPCPort:            7246,
	RPCShutdownTimeout: 5 * time.Second,

	// Tracker
	TrackerNumWant:              200,
	TrackerStopTimeout:          5 * time.Second,
	TrackerMinAnnounceInterval:  time.Minute,
	TrackerHTTPTimeout:          10 * time.Second,
	TrackerHTTPPrivateUserAgent: "Rain/" + "v",
	TrackerHTTPMaxResponseSize:  2 << 20,
	TrackerHTTPVerifyTLS:        true,

	// DHT node
	DHTEnabled:             true,
	DHTHost:                "0.0.0.0",
	DHTPort:                7246,
	DHTAnnounceInterval:    30 * time.Minute,
	DHTMinAnnounceInterval: time.Minute,
	DHTBootstrapNodes: []string{
		"router.bittorrent.com:6881",
		"dht.transmissionbt.com:6881",
		"router.utorrent.com:6881",
		"dht.libtorrent.org:25401",
		"dht.aelitis.com:6881",
	},

	// Peer
	UnchokedPeers:                3,
	OptimisticUnchokedPeers:      1,
	MaxRequestsIn:                250,
	MaxRequestsOut:               250,
	DefaultRequestsOut:           50,
	RequestTimeout:               20 * time.Second,
	EndgameMaxDuplicateDownloads: 20,
	MaxPeerDial:                  80,
	MaxPeerAccept:                20,
	ParallelMetadataDownloads:    2,
	PeerConnectTimeout:           5 * time.Second,
	PeerHandshakeTimeout:         10 * time.Second,
	PieceReadTimeout:             30 * time.Second,
	MaxPeerAddresses:             2000,
	AllowedFastSet:               10,

	// IO
	ReadCacheBlockSize: 128 << 10,
	ReadCacheSize:      256 << 20,
	ReadCacheTTL:       1 * time.Minute,
	ParallelReads:      1,
	ParallelWrites:     1,
	WriteCacheSize:     1 << 30,

	// Webseed settings
	WebseedDialTimeout:             10 * time.Second,
	WebseedTLSHandshakeTimeout:     10 * time.Second,
	WebseedResponseHeaderTimeout:   10 * time.Second,
	WebseedResponseBodyReadTimeout: 10 * time.Second,
	WebseedRetryInterval:           time.Minute,
	WebseedVerifyTLS:               true,
	WebseedMaxSources:              10,
	WebseedMaxDownloads:            4,
}

// Download data
func Download(data static.Data) error {

	var wg sync.WaitGroup

	if config.SelectStream == "" {
		config.SelectStream = "0"
	}

	scraperPath, err := filepath.Abs("../")
	if err == nil {
		// set static paths for torrent downloader
		DefaultConfig.Database = filepath.Join(scraperPath, DefaultConfig.Database)
		DefaultConfig.DataDir = filepath.Join(scraperPath, DefaultConfig.DataDir)
	}

	// select stream to download
	var stream static.Stream
	var ok bool
	if stream, ok = data.Streams[config.SelectStream]; !ok {
		return errors.New(fmt.Sprintf("Stream %s not found", config.SelectStream))
	}

	var saveErr error

	for _, URL := range stream.URLs {
		wg.Add(1)
		go func(URL static.URL, title string) {
			defer wg.Done()
			err := save(URL, title, config.FakeHeaders)
			if err != nil {
				saveErr = err
			}
		}(URL, data.Title)
		if saveErr != nil {
			return saveErr
		}
	}
	wg.Wait()

	return nil
}

func save(url static.URL, fileName string, headers map[string]string) error {
	if config.OutputName != "" {
		fileName = config.OutputName
	}

	var filePath string
	if config.OutputPath != "" {
		filePath = config.OutputPath
	}

	if !strings.HasSuffix(url.URL, ".torrent") {
		file, err := os.Create(filePath + fileName + "." + url.Ext)
		if err != nil {
			return err
		}

		_, err = writeFile(url.URL, file, headers)
		if err != nil {
			return err
		}
	} else {
		resp, err := request.Request(http.MethodGet, url.URL, map[string]string{
			"Accept-Encoding": "gzip, deflate, br",
		})
		if err != nil {
			return nil
		}
		defer resp.Body.Close()

		if filePath != "" {
			DefaultConfig.DataDir = filePath
		}

		session, err := torrent.NewSession(DefaultConfig)
		if err != nil {
			return err
		}
		defer session.Close()
		torrent, err := session.AddTorrent(resp.Body, &torrent.AddTorrentOptions{StopAfterDownload: true})
		if err != nil {
			return err
		}

		select {
		case <-torrent.NotifyComplete():
		case <-torrent.NotifyStop():
			log.Fatal(torrent.Stats().Error)
		default:
			stats := torrent.Stats()
			fmt.Println(fmt.Sprintf("Downloading %s - time left: %b - downloading with %bB/s - uploading with %bB/s", fileName, stats.ETA, stats.Speed.Download, stats.Speed.Upload))
			time.Sleep(1 * time.Second)
		}

	}

	return nil
}

func writeFile(url string, file *os.File, headers map[string]string) (int64, error) {
	res, err := request.Request(http.MethodGet, url, headers)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	bar := progressbar.NewOptions(
		int(res.ContentLength),
		progressbar.OptionSetBytes(int(res.ContentLength)),
		progressbar.OptionSetDescription(fmt.Sprintf("Downloading %s ...", file.Name())),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetRenderBlankState(true),
		//progressbar.OptionShowCount(),
	)
	writer := io.MultiWriter(file, bar)
	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	written, copyErr := io.Copy(writer, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return written, fmt.Errorf("file copy error: %s", copyErr)
	}
	return written, nil
}
