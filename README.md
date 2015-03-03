kowa
====

Static website manager


## Development setup

First, you need a running mongodb server running on standard port.

Checkout the project:

    $ git clone git@github.com:aymerick/kowa.git

Build kowa:

    $ cd kowa
    $ go build

Setup the database:

    $ ./kowa setup

Add a user with two sites:

    $ ./kowa add_user mike mike@asso.ninja Michelangelo TMNT pizzaword

    $ ./kowa add_site site1 'My First Site' mike
    $ ./kowa add_site site2 'My Second Site' mike

Setup theme depencencies:

    $ npm install -g grunt-cli
    $ npm install -g bower
    $ gem install bundle

Build the theme:

    $ cd themes/willy

    $ npm install
    $ bower install
    $ bundle install

    $ grunt build

Start server (the `-s` switch activates serving of static sites):

    $ cd ../..
    $ ./kowa server -s

The server is now waiting for API requests on port `35830` and serves generated sites on port `48910`.

Start client:

    $ cd client
    $ ember server --proxy http://127.0.0.1:35830

The admin app is now running here: <http://127.0.0.1:4200> and you can login with credentials: `mike` / `pizzaword`.

If you want to disable live reload (when testing image upload for example), starts the client with the `--live-reload` flag like that:

    $ ember server --proxy http://127.0.0.1:35830 --live-reload=false


## Test

To launch tests, go to `kowa` root directory, then:

    $ go test ./... -v

    $ cd client
    $ ember test


## Sublime Text

Search Where: -*/bower_components/*,-*/node_modules/*,-*/client/dist/*,-*/client/tmp/*,-*.min.js,-*.min.css,-*.css.map,-*.min.map,-*.svg,-*hugo/test_site/public/*
