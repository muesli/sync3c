sync3c
=======

A little tool to sync/mirror media from https://media.ccc.de

## Installation

Make sure you have a working Go environment. Follow the [Go install instructions](http://golang.org/doc/install.html).

First of all you need to checkout the source code:

    go get github.com/muesli/sync3c

If you want to build it manually:

    cd $GOPATH/src/github.com/muesli/sync3c
    go get -u -v
    go build

## Example usage

#### Download all talks from all conferences, best available quality & original language only
```
$ $GOPATH/bin/sync3c -destination="/my/downloads"
Conference: 33C3: works for me, acronym: 33c3 (URL: https://api.media.ccc.de/public/conferences/101)
Event: Bonsai Kitten waren mir lieber - Rechte Falschmeldungen in sozialen Netzwerken - Auf der Hoaxmap werden seit vergangenem Febru...
        Found video (video/mp4): 34 minutes, 1920x1080 (HD: true, 313MiB) https://api.media.ccc.de/public/recordings/13601
        Found audio (audio/mpeg): 33 minutes (HD: false, 31MiB) https://api.media.ccc.de/public/recordings/13797
        Found audio (audio/opus): 33 minutes (HD: false, 24MiB) https://api.media.ccc.de/public/recordings/13834
Downloading: http://cdn.media.ccc.de/congress/2016/h264-hd/33c3-8288-deu-Bonsai_Kitten_waren_mir_lieber_-_Rechte_Falschmeldungen_in_sozialen_Netzwerken.mp4
Rechte Falschmeldungen-in-sozialen-Netzwerken (...).mp4  313.75 MiB / 313.75 MiB [#################################################] 100.00%
```

#### Download all talks from a specific conference, best available quality, including all translations
```
$ $GOPATH/bin/sync3c -acronym eh16 -destination="/my/downloads"
Conference: Easterhegg 2016, acronym: eh16 (URL: https://api.media.ccc.de/public/conferences/81)
...
```

### Download all talks, best available quality & original language only
```
$ $GOPATH/bin/sync3c -ignoreTranslations=false -destination="/my/downloads"
...
```

Enjoy the great content!
