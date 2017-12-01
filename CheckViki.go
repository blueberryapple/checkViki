package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "log"
    "flag"
    "strings"
)

// instance variables
var vikiApi = "https://api.viki.io/v4/"
var appId = "100444a"

func compStr(s1, s2 string) bool {
    left := strings.ToLower(s1)
    right := strings.ToLower(s2)
    return strings.Contains(left, right)
}

func report(err error) {
    if (err != nil) {
        log.Fatal(err)
    }
}

func getJson(url string) []byte{
    // grabs url response
    resp, err := http.Get(url)
    report(err)

    defer resp.Body.Close()

    // retrieves body
    body, err := ioutil.ReadAll(resp.Body)
    report(err)

    return body
}

func getId(rawName string) string{
    name := strings.Replace(rawName, " ", "+", -1)
    lookUp := "search.json?c=" + name + "&"
    url := vikiApi + lookUp + "app=" + appId
    search := getJson(url)

    // sets up struct
    type Result struct {
        Id string
        Tt string
    }

    // parses result for series id
    var res []Result
    json.Unmarshal(search, &res)
    if len(res) != 0 && compStr(res[0].Tt, rawName) {
        return res[0].Id
    } else {
        fmt.Println("series not found. blame viki api")
        return ""
    }
}

func getCent(id string, ep int) int{
    // prep and grab json data
    lookup := "containers/" + id + "/episodes.json?"
    url := vikiApi + lookup + "app=" + appId
    eps := getJson(url)

    // define the structure of the json
    type Ep struct {
        Subtitle_completions struct {
            En int
        }
        Number int
    }

    type Resp struct {
        Response []Ep
    }

    // parse subtitle percents for episodes
    res := &Resp{}
    err := json.Unmarshal(eps, &res)
    report(err)

    // reverse index
    numEps :=  res.Response[0].Number
    i := numEps - ep
    if (i > len(res.Response) - 1 || i < 0) {
        fmt.Println("latest episode is:", numEps)
        fmt.Println("queried episode too old or does not exist yet")
        return -1
    }

    // retrieves percentage
    cent := res.Response[i].Subtitle_completions.En
    return cent
}

func main() {
    // argument parser
    var name = flag.String("series", "bong soon", "kdrama series to look up")
    var ep = flag.Int("episode", 16, "episode number")
    flag.Parse()

    // sets up variables    
    id := getId(*name)
    if (len(id) == 0) {
        return 
    }

    cent := getCent(id, *ep)
    fmt.Println(*name, "episode", *ep, "subbed at", cent)
}
