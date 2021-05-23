# go-hentai-scraper

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/gan-of-culture/go-hentai-scraper/Go)

Is a commandline tool to download hentai made in golang by me. This scraper is not using official APIs since some of them have limitations. Unfortunatly this also means that for some sites it might take longer to download content because the connection will be restricted.  
You can see what site is supported to what extent [here](#supported-sites).

## Setup guide

If you don't want to build the app yourself checkout the [releases page](https://github.com/gan-of-culture/go-hentai-scraper/releases).

Otherwise you can clone this Repository yourself and download it. Then just do:

```bash
go build
```

After that you should be provided with an **executable**.

Or you can just do

```bash
go run main.go ...
```

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

 -i             Show info for the provided URL

 go-hentai-scraper -i http...

--------------------------------------------------------------------------------

 -s             Select a specific stream

 go-hentai-scraper -s 0 http...

--------------------------------------------------------------------------------

 -t             Number of download threads (works if you download multiple files)

 go-hentai-scraper -t 4 http...

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
| [animeidhentai.com (1080p, 720p, 480p, 360p)](https://animeidhentai.com) |         ?          |:heavy_check_mark:|
| [booruproject (ex. rule34, gelbooru)](https://booru.org/top)             | :heavy_check_mark: |:heavy_check_mark:|
| [booru.io](https://booru.io/)                                            | :heavy_check_mark: |        ?         |
| [damn.stream](https://www.damn.stream)                                   |         ?          |:heavy_check_mark:|
| [danbooru.donmai.us](https://danbooru.donmai.us)                         | :heavy_check_mark: |        ?         |
| [e-hentai.org](http://e-hentai.org/)                                     | :heavy_check_mark: |        ?         |
| [exhentai.org*](http://exhentai.org/)                                    | :heavy_check_mark: |        ?         |
| [hanime.tv(1080p, 720p, 480p, 360p)](https://hanime.tv)                  | :heavy_check_mark: |:heavy_check_mark:|
| [hentai2w.com(720p, 480p, 360p)](https://hentai2w.com)                   |         ?          |:heavy_check_mark:|
| [hentaicloud.com(720p)](https://www.hentaicloud.com)                     |        :x:         |:heavy_check_mark:|
| [hentaidude.com(720p, 480, 360p)](https://hentaidude.com/)               |         ?          |:heavy_check_mark:|
| [hentaihaven.red (1080p, 720p, 480p, 360p)](https://hentaihaven.red)     |         ?          |:heavy_check_mark:|
| [hentaihaven.xxx (1080p, 720p, 480p, 360p)](https://hentaihaven.xxx)     |         ?          |:heavy_check_mark:|
| [hentaimama.io(1080p, 720p)](https://hentaimama.io)                      |         ?          |:heavy_check_mark:|
| [hentais.tube (720p, 480p, 360p)](https://www.hentais.tube/)             |         ?          |:heavy_check_mark:|
| [hentaistream.moe (2160p, 1080p, 480p)](https://hentaistream.moe/)       |         ?          |:heavy_check_mark:|
| [hentaistream.xxx (1080p, 720, 480p, 360p)](https://hentaistream.xxx/)   |         ?          |:heavy_check_mark:|
| [hentaiworld.tv (1080p, 720p, 480p)](https://hentaiworld.tv/)            |         ?          |:heavy_check_mark:|
| [hentai.tv (1080p, 720p, 480p, 360p)](https://hentai.tv/)                |         ?          |:heavy_check_mark:|
| [hentaiyes.com (1080p, 720p, 480p, 360p)](https://hentaiyes.com/)        |         ?          |:heavy_check_mark:|
| [hitomi.la](https://hitomi.la/)                                          | :heavy_check_mark: |        ?         |
| [konachan.com](https://konachan.com/post?tags=)                          | :heavy_check_mark: |        ?         |
| [miohentai.com (1080p, 720p, 480p)](https://miohentai.com/)              | :heavy_check_mark: |:heavy_check_mark:|
| [muchohentai.com (1080p, 720p, 480p, 360p)](https://muchohentai.com/)    |         ?          |:heavy_check_mark:|
| [nhentai.net](https://nhentai.net)                                       | :heavy_check_mark: |        ?         |
| [pururin.io](https://pururin.io)                                         | :heavy_check_mark: |        ?         |
| [rule34.paheal.net](https://rule34.paheal.net)                           | :heavy_check_mark: |:heavy_check_mark:|
| [yandere.re](https://yande.re/post)                                      | :heavy_check_mark: |        ?         |

*you need a login for this site. You can supply it via the parameters -un and -up

If your site is not listed - you can still try to use the universal downloader.

## TODO's

- Implement concurrency for big single file downloads
- Clean up coding and add more sites
- Speed improvements

## Donate

Donating is completly free and you gain no extra benefits from it. Although it's very much appriciated.

```bash
XMR 4AFThbPDiig6tEZdRL4NnvDfqPETiuewDgpCJKkSs11BGCVqoydRUHkZr5cotGMx395V7c2swDxi5Xjhbztiqyod7P31szF
```

## Licencse

Pretty sure [MIT](LICENSE) is the way to go
