# go-hentai-scraper

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/gan-of-culture/go-hentai-scraper/Go)

Is a commandline tool to download hentai made in golang by me. This scraper is not using official APIs since some of them have limitations. Unfortunatly this also means that for somesites it might take longer to download content because the connection will be restricted.  
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

 -t             Number of download threads

 go-hentai-scraper -t 4 http...

--------------------------------------------------------------------------------

 -a              Amount of files (image boards only)

 go-hentai-scraper -a 5000 http...

--------------------------------------------------------------------------------

 -p              Enter pages like 1,2,3-4,6,7,8-9 for doujins

 go-hentai-scraper -p 1,2,3-4 http...
```

## Supported sites

| Site                                                             | Images             | Videos           |
| -----------------------------------------------------------------|:------------------:|:----------------:|
| [booruproject (ex. rule34, gelbooru)](https://booru.org/top)     | :heavy_check_mark: |:heavy_check_mark:|
| [booru](https://booru.io/)                                       | :heavy_check_mark: |        ?         |
| [danbooru](https://danbooru.donmai.us)                           | :heavy_check_mark: |        ?         |
| [e-hentai](http://e-hentai.org/)                                 | :heavy_check_mark: |        ?         |
| [exhentai*](http://exhentai.org/)                                | :heavy_check_mark: |        ?         |
| [hanime](https://hanime.tv)                                      | :heavy_check_mark: |       :x:        |
| [hentaimama](https://hentaimama.io)                              |         ?          |:heavy_check_mark:|
| [hentais (480p only)](https://www.hentais.tube/)                 |         ?          |:heavy_check_mark:|
| [hentaistream (2160p, 1080p, 480p)](https://hentaistream.moe/)   |         ?          |:heavy_check_mark:|
| [hentaiworld (1080p, 720p, 480p)](https://hentaiworld.tv/)       |         ?          |:heavy_check_mark:|
| [konachan](https://konachan.com/post?tags=)                      | :heavy_check_mark: |        ?         |
| [nhentai](https://nhentai.net)                                   | :heavy_check_mark: |        ?         |
| [rule34](https://rule34.paheal.net)                              | :heavy_check_mark: |:heavy_check_mark:|
| [yandere](https://yande.re/post)                                 | :heavy_check_mark: |        ?         |

* you need a login for this site. You can supply it via the parameters -un and -up

If your site is not listed - you can still try to use the universal downloader
This works really good for the pitures of hanime.tv. This also works for .torrent urls

## TODO's

- Clean up coding and add more sites
- Speed improvements

## Licencse

Pretty sure [MIT](LICENSE) is the way to go
