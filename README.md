kowa
====

Association website generation


## Development

Start server:

    $ go build
    $ ./kowa server

Start client:

    $ cd client
    $ ember server --proxy http://127.0.0.1:35830

Start without live reload:

    $ ember server --proxy http://127.0.0.1:35830 --live-reload=false

Test:

    $ go test ./... -v


## Sublime Text

Search Where: -*/bower_components/*,-*/node_modules/*,-*/client/dist/*,-*/client/tmp/*,-*.min.js,-*.min.css,-*.css.map,-*.svg,-*hugo/test_site/public/*
