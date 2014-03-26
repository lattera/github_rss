/*
(BSD 2-clause license)

Copyright (c) 2014, Shawn Webb
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

   * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
    "fmt"
    rss "github.com/jteeuwen/go-pkg-rss"
    "os"
    "time"
    "encoding/json"
    "strings"
)

type Project struct {
    Name string
    Branch string
}

type Config struct {
    Polltime int
    Projects []Project
}

func main() {
    config := new(Config)
    sleeper := make(chan int, 1)

    file, err := os.Open("config.json")
    if err != nil {
        fmt.Printf("Could not open the config file: %s\n", err.Error())
        return
    }

    defer file.Close()

    decoder := json.NewDecoder(file)
    decoder.Decode(&config)

    for _, project := range config.Projects {
        go PollFeed("https://github.com/" + project.Name + "/commits/" + project.Branch + ".atom", config.Polltime)
    }

    <-sleeper
}

func PollFeed(uri string, timeout int) {
    feed := rss.New(timeout, true, chanHandler, itemHandler)
    for {
        if err := feed.Fetch(uri, nil); err != nil {
            fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
                return
        }

        <-time.After(time.Duration(feed.SecondsTillUpdate()) * time.Second)
    }
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
    return
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
    for _, item := range newitems {
        fmt.Printf("%s\t%s\t%s\t%s\t%s\n", item.PubDate, getProjectName(feed.Url), getProjectBranch(feed.Url), item.Title, (*item.Links[0]).Href)
    }
}

func getProjectName(uri string) string {
    name := uri[len("https://github.com/"):]
    name = name[0:strings.LastIndex(name, "/") - len("/commits")]
    return name
}

func getProjectBranch(uri string) string {
    name := getProjectName(uri)
    return strings.Replace(uri[len("https://github.com/") + len(name) + len ("/commits/"):], ".atom", "", 1)
}
