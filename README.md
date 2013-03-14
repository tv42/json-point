json-point -- Perform JSON Pointer queries from the command line
================================================================

This is a command-line tool to easily extract parts of JSON documents.
The query language is JSON Pointer.

JSON is always read from standard input; to read from a file, redirect
input with `json-point <FILE ARGS..`

Examples:

To extract a field:

    $ echo '{"foo":"bar", "quux":"thud"}' | json-point /foo
    "bar"
    $ echo '{"foo":{"quux":"thud"}}' | json-point /foo
    {"quux":"thud"}
    $ echo '{"foo":{"quux":"thud"}}' | json-point /foo/quux
    "thud"

You can extract multiple fields at once:

    $ echo '{"foo":"bar", "quux":"thud"}' | json-point -pretty /quux /foo
    "thud"
    "bar"

If a field is not found, an empty line is printed in its place, and
exit status will be 1:

    $ echo '{"foo":"bar", "quux":"thud"}' | json-point /nope
    
    # exits with status 1

You can process multiple input documents, newline-separated or not, at
once. Multiple queries can be run on each document, one at a time:

    $ echo '{"foo":"bar", "quux":"thud"}{"foo": "xyzzy"}' | json-point /foo
    "bar"
    "xyzzy"
    $ echo '{"foo":"bar", "quux":"thud"}{"foo": "xyzzy"}' | json-point /foo /quux
    "bar"
    "thud"
    "xyzzy"
    
    # exits with status 1

You can list all the possible queries that could be performed, given
this input. Note that the result includes an empty line; that's a
query that matches the whole document, and is executed as `json-point
''`.

    $ echo '{"foo":"bar", "quux":"thud"}' | json-point -list
    
    /foo
    /quux

To extract the value of a JSON string, without the quotes, use the
flag `-pretty`:

    $ echo '{"foo":"bar", "quux":"thud"}' | json-point -pretty /foo
    bar

The `-pretty` flag is ignored if the match doesn't have a pretty
format; it does nothing for JSON objects or arrays:

    $ echo '{"foo":{"quux":"thud"}}' | json-point -pretty /foo
    {"quux":"thud"}
    $ echo '{"foo":["bar"]}' | json-point -pretty /foo
    ["bar"]

As the result is unquoted, a single output entry may span multiple
lines.

    $ echo '{"foo":"bar\nbaz"}' | json-point -pretty /foo
    bar
    baz

In general, `-pretty` is meant for human consumption only, but you
could use it to e.g. extract JSON fields into shell variables:

    $ FOO="$(echo '{"foo":"bar"}' | json-point -pretty /foo)"
    $ printf 'The value is %s.\n' "$FOO"
    The value is bar.

To install `json-point`, you need a working Go installation. Please
see http://golang.org/ for that. Then run:

    go get github.com/tv42/json-point

This program uses a JSON Pointer library by Dustin Sallings:
https://github.com/dustin/go-jsonpointer
