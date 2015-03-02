kowa
====

Static website management


## Development

You need a running mongodb server.

Start server:

    $ go build
    $ ./kowa server -s

Start client:

    $ cd client
    $ ember server --proxy http://127.0.0.1:35830

Start client without live reload:

    $ ember server --proxy http://127.0.0.1:35830 --live-reload=false

Browse generated site:

    <http://127.0.0.1:48910/>

Test:

    $ go test ./... -v

    $ cd client
    $ ember test


## Sublime Text

Search Where: -*/bower_components/*,-*/node_modules/*,-*/client/dist/*,-*/client/tmp/*,-*.min.js,-*.min.css,-*.css.map,-*.min.map,-*.svg,-*hugo/test_site/public/*
