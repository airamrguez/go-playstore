# Google PlayStore
A package to query the _Google Play Store_.

[![build status](https://secure.travis-ci.org/airamrguez/go-playstore)](http://travis-ci.org/airamrguez/go-playstore)
**NOTE**: This project is a work in progress.

## Installation
This project requires Go1.1+.

    $ go get github.com/airamrguez/go-playstore

## Examples

There are mainly two actions lookups and searchs.

### Lookup example

This example shows how to fetch data for the Candy Crush Saga in english, spanish
and french.

```Go
import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "github.com/airamrguez/go-playstore"
)

func main() {
    // For an app that doesn't use Google App Engine you can directly
    // use net/http package.
    httpGet := func(url *url.URL) (*http.Response, error) {
        return http.Get(url)
    }
    // English content is always fetched.
    app, err := playstore.MultiLookUp(httpGet, "com.king.candycrushsaga", ["es", "fr"])
    if err != nil {
        panic(err)
    }
    json, err := json.Marshal(app)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(json))
}
```

### Search example

Search all applications that match a term. In this case we search for all
applications matching candy.

```Go
import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "github.com/airamrguez/go-playstore"
)

func main() {
    httpGet := func(url *url.URL) (*http.Response, error) {
        fmt.Println(url.String())
        return http.Get(url.String())
    }
    apps, err := playstore.Search(httpGet, "candy", 40, "en")
    if err != nil {
        panic(err)
    }
    json, err := json.Marshal(apps)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(json))
}
```

## Changelog

 *  Look up with multilingual support.
 *  Search a term.

## TODO

- [ ] Add tests
- [ ] Add leadership
