# Google PlayStore
A service to retrieve structured info from the _Google Play Store_.

## Installation
This project requires Go1.1+.

    $ go get github.com/airamrguez/go-playstore

## Examples

There are mainly two actions:
 - **Look Up**:
 - **Search**:

```Go
import (
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
