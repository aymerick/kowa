# kowa [![Build Status](https://secure.travis-ci.org/aymerick/kowa.svg?branch=master)](http://travis-ci.org/aymerick/kowa)

The static website manager.

**WARNING: This is a work in progress, tests are missing, documentation is missing... a lot of stuff is missing, and it has NOT been deployed in production yet.**

Build a typical showcase website thanks to a modern admin web app. You website is generated statically everytime you make a change to it.

A static website is easy to deploy, cost effective and really fast.

With kowa, create a website for your organisation with:

  - a customizable `homepage`
  - a `contact` page
  - an `activities` page to clearly explain what your organisation is about
  - a `posts` page: post your latests news, and publish them automatically on social networks *(not implemented yet)*
  - an `events` page: inform your audience about your ongoing events
  - a `members` page featuring your organisation team
  - and as many custom pages as you want

All these features are optionnals: only take what your need.

The server is written in Go and the client is an Ember.js application that you can find at <https://github.com/aymerick/kowa-client>.


## Development setup

### Client

Follow instructions at: <https://github.com/aymerick/kowa-client>


### Themes

Fetch all themes:

    $ git clone --recursive git@github.com:aymerick/kowa-themes.git


### Server

You need a running mongodb server running on standard port.

Fetch kowa:

    $ go get git@github.com:aymerick/kowa

Build kowa:

    $ cd $GOPATH/src/github.com/aymerick/kowa
    $ make build

Setup the database:

    $ ./kowa setup -t `/path/to/kowa-themes` -u `/path/to/kowa-client/public/upload`

  - The `-t` flag points to the `kowa-themes` directory you previously cloned.
  - The `-u` flag is mandatory and indicates where uploaded files are stored (ie. the `/public/upload` dir of `kowa-client`).

Add a user with two sites:

    $ ./kowa add_user mike mike@asso.ninja Michelangelo TMNT pizzaword

    $ ./kowa add_site site1 'My First Site' mike
    $ ./kowa add_site site2 'My Second Site' mike

Start server:

    $ cd ../..
    $ ./kowa server -s -t `/path/to/kowa-themes` -u `/path/to/kowa-client/public/upload`

The `-s` flag activates serving of static sites.

The server is now waiting for API requests on port `35830` and serves generated sites on port `48910`.


## Configuration file

If you want to get rid of passing flags to `kowa` commands, just create a `$HOME/.kowa/config.toml` config file, with [TomML](https://github.com/toml-lang/toml) syntax like that:

    upload_dir = "/path/to/kowa-client/public/upload"
    themes_dir = "/path/to/kowa-themes"
    serve_output = true

Now, you can start the server without flags:

    $ ./kowa server


## Development workflow

When you change server code, you have to rebuild it with `go build` and restart it.

When you change a `SASS` file in a theme you don't have to rebuild the server, but you have to rebuild the theme, for example:

    $ cd kowa-themes/willy
    $ grunt build

Every time you make a change on a site thanks to the client app, the corresponding static site is rebuilt in the background.

You can still trigger a manual rebuild of a static site with this command:

    $ ./kowa build site1 -t `/path/to/kowa-themes` -u `/path/to/kowa-client/public/upload`

If you modify the code that handles images, you can regenerate all derivatives for a given site with this command:

    $ ./kowa gen_derivatives site1 -u `/path/to/kowa-client/public/upload`


## Test

To launch tests, go to `kowa` root directory, then:

    $ go test ./... -v
