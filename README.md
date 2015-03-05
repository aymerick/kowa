kowa
====

Static website manager

**WARNING: This is a work in progress, tests are missing, language selection is missing, documentation is missing... a lot of stuff is missing, and it has NOT been deployed in production yet.**

The server is written in Go, and the client is an Ember.js application inside `client` directory.


## Development setup

First, you need a running mongodb server running on standard port.

Then checkout the project:

    $ git clone git@github.com:aymerick/kowa.git

Build kowa:

    $ cd kowa
    $ go build -i

Setup the database:

    $ ./kowa setup

Add a user with two sites:

    $ ./kowa add_user mike mike@asso.ninja Michelangelo TMNT pizzaword

    $ ./kowa add_site site1 'My First Site' mike
    $ ./kowa add_site site2 'My Second Site' mike

Setup theme building depencencies:

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

The admin app is now running at <http://127.0.0.1:4200> and you can login with credentials: `mike` / `pizzaword`.

If you want to disable live reload (when testing image upload for example), set the `--live-reload`:

    $ ember server --proxy http://127.0.0.1:35830 --live-reload=false


## Development workflow

When you change client code, the app is rebuilt automatically, and the browser reloads it (unless `--live-reload=false` flag is set).

When you change server code, you have to rebuild it with `go build` and restart it.

When you change a `SASS` file in `willy` theme you don't have to rebuild the server, by you have to rebuild the theme:

    $ cd themes/willy
    $ grunt build

Every time you make a change on a site thanks to the admin app, the corresponding static site is rebuilt in the background.

You can still trigger a manual rebuild of a static site with this command:

    $ ./kowa build site1

If you modify the code that handle images, you can regenerate all derivatives for a given site with this command:

    $ ./kowa gen_derivatives site1


## Test

To launch tests, go to `kowa` root directory, then:

    $ go test ./... -v

    $ cd client
    $ ember test


## Sublime Text

Search Where: -*/bower_components/*,-*/node_modules/*,-*/client/dist/*,-*/client/tmp/*,-*.min.js,-*.min.css,-*.css.map,-*.min.map,-*.svg,-*hugo/test_site/public/*
