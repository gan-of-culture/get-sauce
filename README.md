# get-sauce <!-- omit in toc -->

[![Build](https://github.com/gan-of-culture/get-sauce/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/gan-of-culture/get-sauce/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gan-of-culture/get-sauce)](https://goreportcard.com/report/github.com/gan-of-culture/get-sauce)

Is a command line program to download Hentai videos and images from multiple websites. Most sites have very intrusive advertisements, disable downloading of content, limit download speeds or only provide lesser video resolutions, but this tool will always give you the opportunity to download the best available quality from each site. It's also possible to input URLs to a category or multiple URLs and get-sauce will download it all.

- [Installation](#installation)
- [Getting started](#getting-started)
  - [Download example](#download-example)
  - [Multiple inputs](#multiple-inputs)
  - [Captions](#captions)
- [Options](#options)
- [Supported sites](#supported-sites)
  - [Site requirements](#site-requirements)
- [Credit](#credit)
- [Donate](#donate)
- [Licencse](#licencse)

## Installation

The following dependencies are required and must be installed separately.

- **[FFmpeg](https://www.ffmpeg.org)**

> **Note**: FFmpeg does not affect the download, only affects the final file merge.

If you don't want to build the app yourself checkout the [releases page](https://github.com/gan-of-culture/get-sauce/releases).

Otherwise you can use [Golang](https://golang.org/dl/) to install or update it:

```console
go install github.com/gan-of-culture/get-sauce@latest
```

Or if you got Golang already installed you can download the source code and run the program from the source code this:

```console
go run main.go ...
```

## Getting started

Usage:

```console
get-sauce [OPTIONS] URL [URLs...]
```

### Download example

```console
get-sauce https://hentaimama.io/episodes/torokase-orgasm-animation-episode-1/
```

> Note: wrap the URL(s) in quotation marks if it contains special characters.  
> ```get-sauce "https://..."```

The ```-i``` option displays all available quality of video without downloading.

```console
get-sauce -i https://hentaimama.io/episodes/torokase-orgasm-animation-episode-1/

 Site:      https://hentaimama.io/
 Title:     Torokase Orgasm The Animation Episode 1
 Type:      video
 Streams:   # All available qualities
     [0]  -------------------
     Type:            video
     Info:            Mirror 1
     Quality:         1280x720
     Parts:           562
     Size:            ~ 427.6 MB
     # download with: get-sauce -s 0 ...


     [1]  -------------------
     Type:            video
     Info:            Mirror 1
     Quality:         842x480
     Parts:           562
     Size:            ~ 123.8 MB
     # download with: get-sauce -s 1 ...


     [2]  -------------------
     Type:            video
     Info:            Mirror 1
     Quality:         640x360
     Parts:           562
     Size:            ~ 120.7 MB
     # download with: get-sauce -s 2 ...


     [3]  -------------------
     Type:            video
     Info:            Mirror 2
     Quality:         unknown
     Size:            186.3 MB
     # download with: get-sauce -s 3 ...


     [4]  -------------------
     Type:            video
     Info:            Mirror 3
     Quality:         unknown
     Size:            186.3 MB
     # download with: get-sauce -s 4 ...
```

The default stream is 0 and it also offers the best available quality. If you want to download a different quality or from a mirrored server you can manually choose a stream with ```-s```.

### Multiple inputs

You can also download multiple URLs at once:

```console
get-sauce -i https://nhentai.net/g/364616/ https://nhentai.net/g/364591/

 Site:      https://nhentai.net
 Title:     Matsuri tte Ii na
 Type:      image
 Streams:   # All available qualities
     [0]  -------------------
     Type:            image
     Quality:         unknown
     Parts:           31
     Size:            0 B
     # download with: get-sauce -s 0 ...


 Site:      https://nhentai.net
 Title:     ASSTROLOGIC
 Type:      image
 Streams:   # All available qualities
     [0]  -------------------
     Type:            image
     Quality:         unknown
     Parts:           36
     Size:            0 B
     # download with: get-sauce -s 0 ...
```

The URLs will be downloaded one by one.

### Captions

For most of the sites the captions (CC, OC or Subtitles) are hard encoded into the video file and can't be downloaded separately. If it is not encoded into the video and a caption file was found, then you can download it using the option ```-c```.

To see if the captions can be downloaded separately use the option ```-i```. There will be extra information displayed if the option ```-c``` can be used.

```console
get-sauce -i https://hentai-moon.com/videos/285/isekai-harem-monogatari-ep-1/

 Site:      https://hentai-moon.com
 Title:     Isekai Harem Monogatari - Ep.1
 Type:      video
 Captions:  # All available languages

     [0]  -------------------
     Language:            English
     # download with: get-sauce -c 0 ...


 Streams:   # All available qualities
     [0]  -------------------
     Type:            video
     Quality:         unknown
     Size:            75.6 MB
     # download with: get-sauce -s 0 ...
```

## Options

```console

 -a             Amount of files (image boards only)

 get-sauce -a 5000 http...

------------------------------------------------------------------------------------------------

 -c             Download caption if separate to a extra file

 get-sauce -c 0 http...

------------------------------------------------------------------------------------------------

 -h             UserHeaders for the http requests. To bypass Cloudflare or DDOS-GUARD protection

 get-sauce -h "cf_clearance=..." http...


------------------------------------------------------------------------------------------------
 
 -i             Show info for the provided URL

 get-sauce -i http...

------------------------------------------------------------------------------------------------

 -j             Show extracted data as json

 get-sauce -j http...

------------------------------------------------------------------------------------------------

 -k             Keep video, audio and subtitles. Don't merge using ffmpeg

 get-sauce -k http...

------------------------------------------------------------------------------------------------
 -o             Output name of the file

 get-sauce -o myfilename http...

------------------------------------------------------------------------------------------------

 -O             Output path of the files

 get-sauce -O C://Users//User//Downloads// http...

------------------------------------------------------------------------------------------------

 -p             Enter pages like 1,2,3-4,6,7,8-9 for doujins

 get-sauce -p 1,2,3-4 http...

------------------------------------------------------------------------------------------------

 -q             Quiet mode - show minimal information 

 get-sauce -q http...

------------------------------------------------------------------------------------------------

 -s             Select a specific stream | 0 is default and has the best quality

 get-sauce -s 0 http...

------------------------------------------------------------------------------------------------

 -t             Truncate file if it already exists

 get-sauce -t http...

------------------------------------------------------------------------------------------------

 -T             Set a custom timeout (in minutes) for the http.client

 get-sauce -T 15 http...

------------------------------------------------------------------------------------------------
 -w             Number of download workers

 get-sauce -w 4 http...

```

## Supported sites

The following links will direct you to adult content. Please keep that in mind!

| Site                                                                                      | Images             | Videos           | Requirements         |
| ------------------------------------------------------------------------------------------|:------------------:|:----------------:|:--------------------:|
| [9hentai.to/ru](https://9hentai.to/)                                                      | :heavy_check_mark: |        ?         |
| [animeidhentai.com (1080p, 720p, 480p, 360p)](https://animeidhentai.com)                  |         ?          |:heavy_check_mark:|
| [booruproject (ex. rule34, gelbooru)](https://booru.org/top)                              | :heavy_check_mark: |:heavy_check_mark:|
| [booru.io](https://booru.io/)                                                             | :heavy_check_mark: |        ?         |
| [comicporn.xxx](https://comicporn.xxx)                                                    | :heavy_check_mark: |        ?         |
| [danbooru.donmai.us](https://danbooru.donmai.us)                                          | :heavy_check_mark: |        ?         |:car:|
| [doujin.sexy](https://doujin.sexy)                                                        | :heavy_check_mark: |        ?         |
| [e-hentai.org](http://e-hentai.org/)                                                      | :heavy_check_mark: |        ?         |
| [exhentai.org](http://exhentai.org/)                                                      | :heavy_check_mark: |        ?         |:closed_lock_with_key:|
| [haho.moe (1080p, 720p, 480p, 360p)](https://haho.moe)                                    |         ?          |:heavy_check_mark:|
| [hanime.tv (720p, 480p, 360p)](https://hanime.tv)                                         |         ?          |:heavy_check_mark:|
| [hentai.tv (1080p, 720p, 480p, 360p)](https://hentai.tv/)                                 |         ?          |:heavy_check_mark:|
| [hentai-moon.com (720p, 480p)](https://hentai-moon.com)                                   |         ?          |:heavy_check_mark:|
| [hentai2read.com](https://hentai2read.com)                                                | :heavy_check_mark: |        ?         |
| [hentai2w.com(720p, 480p, 360p)](https://hentai2w.com)                                    |         ?          |:heavy_check_mark:|
| [hentaicloud.com(720p)](https://www.hentaicloud.com)                                      |        :x:         |:heavy_check_mark:|
| [hentaidude.com(720p, 480, 360p)](https://hentaidude.com/)                                |         ?          |:heavy_check_mark:|
| [hentaiera.com](https://hentaiera.com)                                                    | :heavy_check_mark: |         ?        |
| [hentaienvy.com](https://hentaienvy.com)                                                  | :heavy_check_mark: |         ?        |
| [www.hentai-foundry.com](https://www.hentai-foundry.com/)                                 | :heavy_check_mark: |         ?        |
| [hentaifox.com](https://hentaifox.com)                                                    | :heavy_check_mark: |         ?        |
| [hentaihaven.co (1080p, 720p, 480p, 360p)](https://hentaihaven.co)                        |         ?          |:heavy_check_mark:|
| [hentaihaven.com (1080p, 720p, 480p, 360p)](https://hentaihaven.com)                      |         ?          |:heavy_check_mark:|
| [hentaihaven.red (1080p, 720p, 480p, 360p)](https://hentaihaven.red)                      |         ?          |:heavy_check_mark:|
| [hentaihaven.xxx (1080p, 720p, 480p, 360p)](https://hentaihaven.xxx)                      |         ?          |:heavy_check_mark:|
| [hentaimama.io(1080p, 720p)](https://hentaimama.io)                                       |         ?          |:heavy_check_mark:|
| [hentaipulse.com(720p, 420p)](https://hentaipulse.com)                                    |         ?          |:heavy_check_mark:|
| [hentairox.com](https://hentairox.com)                                                    | :heavy_check_mark: |         ?        |
| [hentaistream.moe (2160p, 1080p, 480p)](https://hentaistream.moe/)                        |         ?          |:heavy_check_mark:|
| [hentaistream.tv (1080p, 720p, 480p, 360p)](https://hentaistream.tv)                      |         ?          |:heavy_check_mark:|
| [hentaistream.xxx (1080p, 720, 480p, 360p)](https://hentaistream.xxx/)                    |         ?          |:heavy_check_mark:|
| [hentaivideos.net (1080p, 720p, 480p, 360p)](https://hentaivideos.net/)                   |         ?          |:heavy_check_mark:|
| [hentaiworld.tv (1080p, 720p, 480p)](https://hentaiworld.tv/)                             |         ?          |:heavy_check_mark:|
| [hentaiyes.com (1080p, 720p, 480p, 360p)](https://hentaiyes.com/)                         |         ?          |:heavy_check_mark:|
| [hentaizap.com](https://hentaizap.com)                                                    | :heavy_check_mark: |        ?         |
| [hitomi.la](https://hitomi.la/)                                                           | :heavy_check_mark: |        ?         |
| [imhentai.com](https://imhentai.xxx)                                                      | :heavy_check_mark: |        ?         |
| [iwara.tv](https://iwara.tv/)                                                             | :heavy_check_mark: |:heavy_check_mark:|
| [konachan.com](https://konachan.com/post?tags=)                                           | :heavy_check_mark: |        ?         |
| [latesthentai.com (1080p, 720p, 480p, 360p)](https://latesthentai.com/)                   |         ?          |:heavy_check_mark:|
| [miohentai.com (1080p, 720p, 480p)](https://miohentai.com/)                               | :heavy_check_mark: |        ?         |
| [nhentai.net](https://nhentai.net)                                                        | :heavy_check_mark: |        ?         |:cookie:|
| [ohentai.org (1080p, 720p, 480p)](https://ohentai.org/)                                   |         ?          |:heavy_check_mark:|
| [oppai.stream (2160p, 1080p, 720p)](https://oppai.stream/)                                |         ?          |:heavy_check_mark:|
| [pururin.to](https://pururin.to)                                                          | :heavy_check_mark: |        ?         |
| [rule34.paheal.net](https://rule34.paheal.net)                                            | :heavy_check_mark: |:heavy_check_mark:|
| [rule34video.com (2160p, 1080p, 720p, 480p, 360p)](https://rule34video.com/)              | :heavy_check_mark: |:heavy_check_mark:|
| [simply-hentai.com](https://www.simply-hentai.com)                                        | :heavy_check_mark: |        ?         |
| [thehentaiworld.com](https://thehentaiworld.com)                                          | :heavy_check_mark: |:heavy_check_mark:|
| [uncensoredhentai.xxx (1080p, 720p, 480p, 360p)](https://uncensoredhentai.xxx/)           |         ?          |:heavy_check_mark:|
| [yandere.re](https://yande.re/post)                                                       | :heavy_check_mark: |        ?         |

You can still try to use the universal downloader, if your site is not listed.

### Site requirements

🍪
--> you need to supply cookies for this extractor to work. Most likely it is used to bypass some kind of Cloudflare or DDOS-Guard protection
> Note: separate the different header values using a newline

```console
get-sauce -h "cookie: cf_clearance=k2TGEnkzhz_PtHs09vMryROlD4O3UZhrDFrU4svgjdM-1665105987-0-150; csrftoken=bLiwSENr0mqSZZ27wan1xdjLazVFoXnnABJu7DtrhbNRUacpbEZhV0Eggc5lD8m5
user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36" http...
```

🔐
--> you need to supply login credentials for this extractor to work

```console
get-sauce -un "MyUserName" -up "MyUserPassword" http...
```

🚗
--> requires [geckodriver](https://github.com/mozilla/geckodriver) to workaround DDOS protection

## Credit

- Thanks to [lux](https://github.com/iawia002/lux) for the great template

## Donate

You won't gain extra benefits from it. Although it's very much appriciated.

```console
XMR 4AFThbPDiig6tEZdRL4NnvDfqPETiuewDgpCJKkSs11BGCVqoydRUHkZr5cotGMx395V7c2swDxi5Xjhbztiqyod7P31szF
```

## Licencse

Pretty sure [MIT](LICENSE) is the way to go
