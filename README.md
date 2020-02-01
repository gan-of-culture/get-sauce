# go-hentai-scraper

Is a commandline tool to download hentai made in go by me. You can see what site is supported in what extent down below.

**[Notice]** This repo is still heavily work in progress - so stuff can change

This will be part of another project I have planned so stay tuned.

## Setup guide

For now you need to clone the repo youself and download it. Then just do:

```bash
go build
```

After that you should be provided with a **executable**. 

Or you can just do 

```bash
go run main.go ...
```

I'll do proper releases later *probably~*

## Options

```console

 -o             Output name of the file

 main -o myfilename http...

--------------------------------------------------------------------------------

 -O             Output path of the files

 main -O c://Users//User//Downloads http...

--------------------------------------------------------------------------------

 -r             Restrict content -> don't download senisble content

 main -r http...

--------------------------------------------------------------------------------

 -i             Show info for the provided URL

 main -i http...

--------------------------------------------------------------------------------

 -s             Select a specific stream

 main -s 0 http...


```

## Supported sites

| Site                                                             | Images             | Videos           |
| -----------------------------------------------------------------|:------------------:|:----------------:|
| [danbooru](https://danbooru.donmai.us)                           | :heavy_check_mark: |        ?         |
| [hanime](https://hanime.tv)                                      | :heavy_check_mark: |       :x:        |
| [nhentai](https://nhentai.net)                                   | :heavy_check_mark: |        ?         |
| [rule34](https://rule34.paheal.net)                              | :heavy_check_mark: |:heavy_check_mark:|
| [underhentai](https://underhentai.net)                           |         ?          |:heavy_check_mark:|
| [e-hentai](http://e-hentai.org/)                                 | :heavy_check_mark: |        ?         |
| [exhentai](https://exhentai.org)                                 |        :x:         |        ?         |

*Note that I currently didn't find a way to download 1080p hentai from hanime.tv and implenting 720p isn't that high on my list - I'll try to find a way in the future*  

*Note exhentai is currently not supported, because it requires a user login and I don't want my user banned - I'll probably add a way to extract the data with your own user and a manual how to get one*

## Teach me

I am fairly new to programming in golang. Any tipps and changes that help me develop myself and this repo are highly appriciated.

## TODO's

- Clean up coding and add more sites
- Really just follow this repo to see what I am up to

## Licencse

Pretty sure [MIT](LICENSE) is the way to go