sync3c
======

[![Latest Release](https://img.shields.io/github/release/muesli/sync3c.svg)](https://github.com/muesli/sync3c/releases)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/sync3c)
[![Build Status](https://github.com/muesli/sync3c/workflows/build/badge.svg)](https://github.com/muesli/sync3c/actions)
[![Go ReportCard](http://goreportcard.com/badge/muesli/sync3c)](http://goreportcard.com/report/muesli/sync3c)

A little tool to sync/download media from https://media.ccc.de

Finds the best available quality of a talk and only downloads that version. It will download each talk in its original language, unless you specify your own preferred language (-language). It will fallback to the original language version if there's no translation available.

Per default it will download all talks from all conferences. Be careful, this will fetch and store several hundred GiB of videos over the network. It's probably a good idea to pass in the name (-name) of a conference you're interested in.

If you don't specify a download destination (-destination) everything will be stored in a folder named "downloads" in your current working directory.

It will create sub-directories for each conference within that destination and will name downloaded files following this schema: "{author} - {title} ({subtitle}) {translation}.ext". If that file already exists, it will skip the download. Hence you should delete the latest started download should you decide to kill the process mid-download.

## Installation

### Packages & Installers

- Arch Linux: [sync3c](https://aur.archlinux.org/packages/sync3c/)
- [Binaries](https://github.com/muesli/sync3c/releases) for Linux, macOS & Windows

### From Source

Make sure you have a working Go environment. Follow the [Go install instructions](http://golang.org/doc/install.html).

To install sync3c, simply run:

    go get github.com/muesli/sync3c

To compile it from source:

    cd $GOPATH/src/github.com/muesli/sync3c
    go get -u -v
    go build

## Usage

#### List all available conferences
```
$ $GOPATH/bin/sync3c list
Conference: 33c3 (33C3: works for me)
Conference: ds2015 (Datenspuren 2015)
Conference: eh16 (Easterhegg 2016)
Conference: 32c3 (32C3: gated communities)
...
```

#### Download all talks from all conferences, best available quality & original language only
```
$ $GOPATH/bin/sync3c -destination /my/downloads
Conference: 33c3 (33C3: works for me)
Event: Bonsai Kitten waren mir lieber - Rechte Falschmeldungen in sozialen Netzwerken - Auf der Hoaxmap werden seit vergangenem Febru...
        Found video (video/mp4): 34 minutes, 1920x1080 (HD: true, 313MiB) https://api.media.ccc.de/public/recordings/13601
        Found audio (audio/mpeg): 33 minutes (HD: false, 31MiB) https://api.media.ccc.de/public/recordings/13797
        Found audio (audio/opus): 33 minutes (HD: false, 24MiB) https://api.media.ccc.de/public/recordings/13834
Downloading: http://cdn.media.ccc.de/congress/2016/h264-hd/33c3-8288-deu-Bonsai_Kitten_waren_mir_lieber_-_Rechte_Falschmeldungen_in_sozialen_Netzwerken.mp4
Rechte Falschmeldungen-in-sozialen-Netzwerken (...).mp4  313.75 MiB / 313.75 MiB [#################################################] 100.00%
```

#### Download all talks from a specific conference, best available quality & original language only
```
$ $GOPATH/bin/sync3c -name eh16 -destination /my/downloads
Conference: eh16 (Easterhegg 2016)
...
```

#### Download all talks, best available quality for a specific language, with fallback to original language
```
$ $GOPATH/bin/sync3c -language eng -destination /my/downloads
...
```

#### Use a different media source
```
$ $GOPATH/bin/sync3c -source media.freifunk.net ...
...
```

Enjoy the great content!
