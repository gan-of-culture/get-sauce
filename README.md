# go-hentai-scraper <!-- omit in toc --> 

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/gan-of-culture/go-hentai-scraper/Go)

This scraper is not using official APIs since some of them have limitations. Unfortunatly this also means that for some sites it might take longer to download content because the connection can be restricted. Also some sites only provide single downloads, downloads to lesser qualities or no downloads at all. This scraper will always give you the opportunity to download the best available quality from each site.  

- [Installation](#installation)
- [Getting started](#getting-started)
  - [Download example](#download-example)
  - [Multiple inputs](#multiple-inputs)
- [Options](#options)
- [Supported sites](#supported-sites)
- [Credit](#credit)
- [Donate](#donate)
- [Licencse](#licencse)

## Installation

If you don't want to build the app yourself checkout the [releases page](https://github.com/gan-of-culture/go-hentai-scraper/releases).

Otherwise you can clone this repository yourself and download it. Then use [Golang](https://golang.org/dl/) to build it yourself:

```console
go build
```

After that you should be provided with an **executable**.

Or if you got Golang already installed you can download the source code and replace the executable with this:

```console
go run main.go ...
```

## Getting started

Usage:

```console
go-hentai-scraper [OPTIONS] URL [URLs...]
```

### Download example

```console
go-hentai-scraper https://hanime.tv/videos/hentai/kuro-gal-ni-natta-kara-shin-yuu-to-shite-mita-season-1
```

> Note: wrap the URL(s) in quotation marks if it contains special characters.  
> ```go-hentai-scraper "https://..."```

The ```-i``` option displays all available quality of video without downloading.

```console
go-hentai-scraper -i https://hanime.tv/videos/hentai/netosis-1

 Site:      https://hanime.tv/
 Title:     Netosis 1
 Type:      video
 Streams:   # All available qualities
     [0]  -------------------
     Info:            netosis-1-1080p-v1x.mp4-v1x.m3u8
     Quality:         1080p; 1920 x 1080
     Parts:           55
     Size:            276.00 MB (276000000 Bytes)
     # download with: go-hentai-scraper -s 0 ...


     [1]  -------------------
     Info:            netosis-1-720p-v1x.mp4-v1x.m3u8
     Quality:         720p; 1280 x 720
     Parts:           55
     Size:            141.00 MB (141000000 Bytes)
     # download with: go-hentai-scraper -s 1 ...


     [2]  -------------------
     Info:            netosis-1-480p-v1x.mp4-v1x.m3u8
     Quality:         480p; 854 x 480
     Parts:           56
     Size:            69.00 MB (69000000 Bytes)
     # download with: go-hentai-scraper -s 2 ...


     [3]  -------------------
     Info:            netosis-1-360p-v1x.mp4-v1x.m3u8
     Quality:         360p; 640 x 360
     Parts:           56
     Size:            49.00 MB (49000000 Bytes)
     # download with: go-hentai-scraper -s 3 ...
```

The default stream is 0 and it also offers the best available quality. If you want to download a different quality or from a mirrored server you can manually choose a stream with ```-s```.

### Multiple inputs

You can also download multiple URLs at once:

```console
go-hentai-scraper -i https://nhentai.net/g/364616/ https://nhentai.net/g/364591/

 Site:      https://nhentai.net
 Title:     Matsuri tte Iina
 Type:      image
 Streams:   # All available qualities
     [0]  -------------------
     Quality:         unknown
     Parts:           31
     Size:            0.00 MB (0 Bytes)
     # download with: go-hentai-scraper -s 0 ...


 Site:      https://nhentai.net
 Title:     ASSTROLOGIC
 Type:      image
 Streams:   # All available qualities
     [0]  -------------------
     Quality:         unknown
     Parts:           36
     Size:            0.00 MB (0 Bytes)
     # download with: go-hentai-scraper -s 0 ...

```

The URLs will be downloaded one by one.

## Options

```console

 -o             Output name of the file

 go-hentai-scraper -o myfilename http...

--------------------------------------------------------------------------------

 -O             Output path of the files

 go-hentai-scraper -O C://Users//User//Downloads// http...

--------------------------------------------------------------------------------

 -r             Restrict content -> don't download senisble content (e-hentai.org only)

 go-hentai-scraper -r http...

--------------------------------------------------------------------------------

 -j             Show extracted data for the provided URL

 go-hentai-scraper -j http...

--------------------------------------------------------------------------------

 -i             Show info for the provided URL

 go-hentai-scraper -i http...

--------------------------------------------------------------------------------

 -s             Select a specific stream | 0 is default and the best quality of that stream

 go-hentai-scraper -s 0 http...

--------------------------------------------------------------------------------

 -w             Number of download workers

 go-hentai-scraper -w 4 http...

--------------------------------------------------------------------------------

 -a              Amount of files (image boards only)

 go-hentai-scraper -a 5000 http...

--------------------------------------------------------------------------------

 -p              Enter pages like 1,2,3-4,6,7,8-9 for doujins

 go-hentai-scraper -p 1,2,3-4 http...
```

## Supported sites

| Site                                                                     | Images             | Videos           |
| -------------------------------------------------------------------------|:------------------:|:----------------:|
| [9hentai.to/ru](https://9hentai.to/)                                     | :heavy_check_mark: |        ?         |
| [animeidhentai.com (1080p, 720p, 480p, 360p)](https://animeidhentai.com) |         ?          |:heavy_check_mark:|
| [booruproject (ex. rule34, gelbooru)](https://booru.org/top)             | :heavy_check_mark: |:heavy_check_mark:|
| [booru.io](https://booru.io/)                                            | :heavy_check_mark: |        ?         |
| [comicporn.xxx](https://comicporn.xxx)                                   | :heavy_check_mark: |        ?         |
| [damn.stream(720p, 480p, 360p)](https://www.damn.stream)                 |         ?          |:heavy_check_mark:|
| [danbooru.donmai.us](https://danbooru.donmai.us)                         | :heavy_check_mark: |        ?         |
| [doujin.sexy](https://doujin.sexy)                                       | :heavy_check_mark: |        ?         |
| [e-hentai.org](http://e-hentai.org/)                                     | :heavy_check_mark: |        ?         |
| [ecchi.iwara.tv](https://ecchi.iwara.tv/)                                | :heavy_check_mark: |:heavy_check_mark:|
| [exhentai.org*](http://exhentai.org/)                                    | :heavy_check_mark: |        ?         |
| [hanime.tv(1080p, 720p, 480p, 360p)](https://hanime.tv)                  | :heavy_check_mark: |:heavy_check_mark:|
| [hentai2read.com](https://hentai2read.com)                               | :heavy_check_mark: |        ?         |
| [hentai2w.com(720p, 480p, 360p)](https://hentai2w.com)                   |         ?          |:heavy_check_mark:|
| [hentaicloud.com(720p)](https://www.hentaicloud.com)                     |        :x:         |:heavy_check_mark:|
| [hentaidude.com(720p, 480, 360p)](https://hentaidude.com/)               |         ?          |:heavy_check_mark:|
| [hentaiera.com](https://hentaiera.com)                                   | :heavy_check_mark: |         ?        |
| [hentaifox.com](https://hentaifox.com)                                   | :heavy_check_mark: |         ?        |
| [hentaihaven.red (1080p, 720p, 480p, 360p)](https://hentaihaven.red)     |         ?          |:heavy_check_mark:|
| [hentaihaven.xxx (1080p, 720p, 480p, 360p)](https://hentaihaven.xxx)     |         ?          |:heavy_check_mark:|
| [hentaimama.io(1080p, 720p)](https://hentaimama.io)                      |         ?          |:heavy_check_mark:|
| [hentais.tube (720p, 480p, 360p)](https://www.hentais.tube/)             |         ?          |:heavy_check_mark:|
| [hentaistream.moe (2160p, 1080p, 480p)](https://hentaistream.moe/)       |         ?          |:heavy_check_mark:|
| [hentaistream.xxx (1080p, 720, 480p, 360p)](https://hentaistream.xxx/)   |         ?          |:heavy_check_mark:|
| [hentaiworld.tv (1080p, 720p, 480p)](https://hentaiworld.tv/)            |         ?          |:heavy_check_mark:|
| [hentai.tv (1080p, 720p, 480p, 360p)](https://hentai.tv/)                |         ?          |:heavy_check_mark:|
| [hentai.pro (1080p, 720p, 480p, 360p)](https://hentai.pro/)              |         ?          |:heavy_check_mark:|
| [hentaiyes.com (1080p, 720p, 480p, 360p)](https://hentaiyes.com/)        |         ?          |:heavy_check_mark:|
| [hitomi.la](https://hitomi.la/)                                          | :heavy_check_mark: |        ?         |
| [imhentai.com](https://imhentai.xxx)                                     | :heavy_check_mark: |        ?         |
| [konachan.com](https://konachan.com/post?tags=)                          | :heavy_check_mark: |        ?         |
| [miohentai.com (1080p, 720p, 480p)](https://miohentai.com/)              | :heavy_check_mark: |:heavy_check_mark:|
| [muchohentai.com (1080p, 720p, 480p, 360p)](https://muchohentai.com/)    |         ?          |:heavy_check_mark:|
| [nhentai.net](https://nhentai.net)                                       | :heavy_check_mark: |        ?         |
| [ohentai.org (1080p, 720p, 480p)](https://ohentai.org/)                  |         ?          |:heavy_check_mark:|
| [pururin.to](https://pururin.to)                                         | :heavy_check_mark: |        ?         |
| [rule34.paheal.net](https://rule34.paheal.net)                           | :heavy_check_mark: |:heavy_check_mark:|
| [simply-hentai.com](https://www.simply-hentai.com)                       | :heavy_check_mark: |        ?         |
| [thehentaiworld.com](https://thehentaiworld.com)                         | :heavy_check_mark: |:heavy_check_mark:|
| [yandere.re](https://yande.re/post)                                      | :heavy_check_mark: |        ?         |

*you need a login for this site. You can supply it via the parameters -un and -up

If your site is not listed - you can still try to use the universal downloader.

## Credit

- Thanks to [annie](https://github.com/iawia002/annie) for the great template

## Donate

You won't gain extra benefits from it. Although it's very much appriciated.

```console
XMR 4AFThbPDiig6tEZdRL4NnvDfqPETiuewDgpCJKkSs11BGCVqoydRUHkZr5cotGMx395V7c2swDxi5Xjhbztiqyod7P31szF
```

## Licencse

Pretty sure [MIT](LICENSE) is the way to go
