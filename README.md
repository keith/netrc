# netrc

This is a simple implementation for reading a [.netrc](http://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-File.html)
file for a user's credentials for services

[![Build Status](https://travis-ci.org/Keithbsmiley/netrc.png)](https://travis-ci.org/Keithbsmiley/netrc)

## Installation

```
$ go get github.com/Keithbsmiley/netrc
```

The add it to your `import`

```
import (
  "github.com/Keithbsmiley/netrc"
)
```

### Development

At the moment this library still has some major feature omissions. Such
as setting new credentials in the user's netrc. I plan on adding these
soon and I will write a better usage guide then. For now most of the
methods are commented and you can view [`netrc_test`](https://github.com/Keithbsmiley/netrc/blob/master/netrc_test.go)
to see some of the method calls.

