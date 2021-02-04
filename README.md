# go-hentai-scraper

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/gan-of-culture/go-hentai-scraper/Go)

Is a commandline tool to download hentai made in go by me. You can see what site is supported in what extent down below.

## Setup guide

Download the right release under the releases tab.

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

 main -o myfilename http...

--------------------------------------------------------------------------------

 -O             Output path of the files

 main -O c://Users//User//Downloads http...

--------------------------------------------------------------------------------

 -r             Restrict content -> don't download senisble content (e-hentai.org only)

 main -r http...

--------------------------------------------------------------------------------

 -i             Show info for the provided URL

 main -i http...

--------------------------------------------------------------------------------

 -s             Select a specific stream

 main -s 0 http...

--------------------------------------------------------------------------------

 -a              Amount of files (booru.io only)

 main -a 5000 http...

--------------------------------------------------------------------------------

 -p              Enter pages like 1,2,3-4,6,7,8-9 for doujins

 main -p 1,2,3-4 http...
```

## Supported sites

| Site                                                             | Images             | Videos           |
| -----------------------------------------------------------------|:------------------:|:----------------:|
| [danbooru](https://danbooru.donmai.us)                           | :heavy_check_mark: |        ?         |
| [booru](https://booru.io/)                                       | :heavy_check_mark: |        ?         |
| [hanime](https://hanime.tv)                                      | :heavy_check_mark: |       :x:        |
| [nhentai](https://nhentai.net)                                   | :heavy_check_mark: |        ?         |
| [rule34](https://rule34.paheal.net)                              | :heavy_check_mark: |:heavy_check_mark:|
| [e-hentai](http://e-hentai.org/)                                 | :heavy_check_mark: |        ?         |
| [exhentai](https://exhentai.org)                                 |        :x:         |        ?         |

*Note exhentai is currently not supported, because it requires a user login and I don't want my user banned - I'll probably add a way to extract the data with your own user and a manual how to get one*

If your site is not listed - you can still try to use the universal downloader
This works really good for the pitures of hanime.tv or reddit.com

## Teach me

I am fairly new to programming in golang. Any tipps and changes that help me develop myself and this repo are highly appriciated.

## TODO's

- Add site to download hentai (video)
- Clean up coding and add more sites

## Licencse

Pretty sure [MIT](LICENSE) is the way to go
