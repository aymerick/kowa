# kowa [![Build Status](https://secure.travis-ci.org/aymerick/kowa.svg?branch=master)](http://travis-ci.org/aymerick/kowa)

The static website manager.

![Kowa Logo](https://github.com/aymerick/kowa/blob/master/kowa.png?raw=true "Kowa")

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


### Database

You need a running mongodb server running on standard port.


### Conf file

It is tedious to pass flags to `kowa` commands, so let's create a `$HOME/.kowa/config.toml` config file, with [TomML](https://github.com/toml-lang/toml) syntax like that:

    upload_dir = "/path/to/kowa-client/upload"
    themes_dir = "/path/to/kowa-themes"
    serve_output = true

  - The `upload_dir` setting indicates where uploaded files are stored (ie. the `/upload` directory of `kowa-client`).
  - The `themes_dir` setting points to the `kowa-themes` directory that you previously cloned.
  - The `serve_output` setting activates serving of static sites for the `server` and `build` commands.


### Install

Fetch kowa:

    $ go get github.com/aymerick/kowa
    $ go get github.com/tools/godep

Build kowa:

    $ cd $GOPATH/src/github.com/aymerick/kowa
    $ make build

Add an administrator user with two sites:

    $ ./kowa add_user mike mike@asso.ninja Michelangelo TMNT pizzaword true
    $ ./kowa add_site site1 'My First Site' mike
    $ ./kowa add_site site2 'My Second Site' mike

Start server:

    $ ./kowa server

The server is now waiting for API requests on port `35830` and serves generated sites on port `48910`.


### Embedded data

Install `go-bindata` package:

    $ go get -u github.com/jteeuwen/go-bindata/...

When you modify translations files in `locales/` or mailers templates in `mailers/templates/`, you have to regenerate `core/bindata.go` file with:

    $ make gen


## Development workflow

When you change server code, you have to rebuild it with `make build` and restart it.

When you change a `SASS` file in a theme you don't have to rebuild the server, but you have to rebuild the theme, for example:

    $ cd kowa-themes/willy
    $ grunt build

Every time you make a change on a site thanks to the client app, the corresponding static site is rebuilt in the background.

You can still trigger a manual rebuild of a static site with this command:

    $ ./kowa build site1

If you modify the code that handles images, you can regenerate all derivatives for a given site with this command:

    $ ./kowa gen_derivatives site1


## Dependencies

### To add a dependcy:

    $ go get foo/bar

Edit your code to import foo/bar, then:

    $ godep save

### To update a dependency:

    $ go get -u foo/bar
    $ godep update foo/bar
    $ godep save


## Test

To launch tests, go to `kowa` root directory, then:

    $ make test
