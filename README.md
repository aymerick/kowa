kowa
====

## Development

Start server:

    $ go build
    $ ./kowa

Start client:

    $ cd client
    $ ember server --proxy http://127.0.0.1:35830


## Hugo website generation

### Build and dev theme

    $ cd hugo/themes/willy
    $ grunt

### Test site

    $ cd hugo/test_site
    $ hugo server -w
