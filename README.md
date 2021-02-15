# go-hentai-scraper

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/gan-of-culture/go-hentai-scraper/Go)

Is a commandline tool to download hentai made in golang by me. You can see what site is supported to what extent down below.

## Setup guide

You can clone this Repository yourself and download it. Then just do:

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

 go-hentai-scraper -O c://Users//User//Downloads http...

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
| [exhentai](https://exhentai.org)                                 |        :x:         |        ?         |
| [hanime](https://hanime.tv)                                      | :heavy_check_mark: |       :x:        |
| [hentais (480p only)](https://www.hentais.tube/)                 |         ?          |:heavy_check_mark:|
| [hentaiworld (1080p, 720p, 480p)](https://hentaiworld.tv/)       |         ?          |:heavy_check_mark:|
| [nhentai](https://nhentai.net)                                   | :heavy_check_mark: |        ?         |
| [rule34](https://rule34.paheal.net)                              | :heavy_check_mark: |:heavy_check_mark:|

*Note exhentai is currently not supported, because it requires a user login and I don't want my user banned - I'll probably add a way to extract the data with your own user and a manual how to get one*

If your site is not listed - you can still try to use the universal downloader
This works really good for the pitures of hanime.tv or reddit.com. This also works for .torrent urls

## TODO's

- Add site to download hentai 1080p only (video)
- Clean up coding and add more sites

## Licencse

Pretty sure [MIT](LICENSE) is the way to go
