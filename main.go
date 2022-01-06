package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/julienschmidt/httprouter"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rodaine/table"
)

type Plugin struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Username    string   `json:"username"`
	Repo        string   `json:"repo"`
	Link        string   `json:"link"`
	Tags        []string `json:"tags"`
	Homepage    string   `json:"homepage"`
	Description string   `json:"description"`
	Branch      string   `json:"branch"`
	OpenIssues  int      `json:"openIssues"`
	Watcher     int      `json:"watcher"`
	Forks       int      `json:"forks"`
	Stars       int      `json:"stars"`
	Subscribers int      `json:"subscribers"`
	Network     int      `json:"network"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

type Result struct {
	Plugin Plugin
	Rank   int
}

type Db struct {
	Plugins map[string]Plugin `json:"plugins"`
}

type Data struct {
	Plugins []Plugin
	Tags    map[string]bool
}

func NewData() *Data {
	return &Data{}
}

type ByRank []Result

func (a ByRank) Len() int      { return len(a) }
func (a ByRank) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByRank) Less(i, j int) bool {
	if a[i].Rank == a[j].Rank {
		return a[i].Plugin.Stars > a[j].Plugin.Stars
	}
	return a[i].Rank < a[j].Rank
}

func getTags(plugins []Plugin) map[string]bool {
	tags := make(map[string]bool)
	for _, plugin := range plugins {
		for _, tag := range plugin.Tags {
			tags[tag] = true
		}
	}
	return tags
}

func fetchData(data *Data) error {
	log.Println("Fetching neovimcraft data ...")

	var db Db
	resp, err := http.Get("https://storage.googleapis.com/neovimcraft.com/db.json")
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&db)
	if err != nil {
		return err
	}
	log.Println("Fetched data successfully")

	plugins := make([]Plugin, 0, len(db.Plugins))

	for _, plugin := range db.Plugins {
		plugins = append(plugins, plugin)
	}

	data.Plugins = plugins
	data.Tags = getTags(plugins)

	return nil
}

func printResults(results []Result, debug bool) table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	if debug {
		tbl := table.New("Score", "Name", "Stars")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, result := range results {
			tbl.AddRow(
				result.Rank,
				result.Plugin.Id,
				result.Plugin.Stars,
			)
		}

		return tbl
	} else {
		tbl := table.New("Name", "Stars", "OpenIssues", "Updated", "Description")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, result := range results {
			tbl.AddRow(
				result.Plugin.Id,
				result.Plugin.Stars,
				result.Plugin.OpenIssues,
				result.Plugin.UpdatedAt,
				result.Plugin.Description,
			)
		}

		return tbl
	}
}

func getRank(needle string, haystack string) int {
	return fuzzy.RankMatch(needle, haystack)
}

func matchTags(search string, plugin Plugin) int {
	for _, tag := range plugin.Tags {
		if search == tag {
			return 0
		}
	}
	return -1
}

func outputResults(w http.ResponseWriter, results *[]Result) {
	sort.Sort(ByRank(*results))

	tbl := printResults(*results, false)
	tbl.WithWriter(w)
	tbl.Print()
}

func helpHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
        columnFmt := color.New(color.FgYellow).SprintfFunc()

        fmt.Fprintf(w, "nvim.sh - neovim plugin search from the terminal\n\n")

        tbl := table.New("api", "description")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt).WithWriter(w)
        tbl.AddRow("https://nvim.sh", "help")
        tbl.AddRow("https://nvim.sh/s", "return all plugins in directory")
        tbl.AddRow("https://nvim.sh/s/:search", "search for plugin within directory")
        tbl.AddRow("https://nvim.sh/t", "list all tags within directory")
        tbl.AddRow("https://nvim.sh/t/:search", "search for plugins that exactly match tag within directory")
        tbl.Print()

        fmt.Fprintf(w, "\npowered by: https://neovimcraft.com\n")
    }
}

func searchTagsHandler(data *Data) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		search := strings.ToLower(ps.ByName("search"))
		results := make([]Result, 0, 10)

		for _, plugin := range data.Plugins {
			ranking := matchTags(search, plugin)
			if ranking >= 0 {
				results = append(results, Result{Plugin: plugin, Rank: ranking})
			}
		}

		outputResults(w, &results)
	}
}

func allHandler(data *Data) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		results := make([]Result, 0, 10)

		for _, plugin := range data.Plugins {
			results = append(results, Result{Plugin: plugin, Rank: 0})
		}

		outputResults(w, &results)
	}
}

func searchHandler(data *Data) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		search := strings.ToLower(ps.ByName("search"))
		results := make([]Result, 0, 10)

		for _, plugin := range data.Plugins {
			rankName := getRank(search, plugin.Name)
			// rankDesc := getRank(search, plugin.Description)
			rankTags := matchTags(search, plugin)
			ranking := (rankName + rankTags) / 2

			if ranking >= 0 {
				results = append(results, Result{Plugin: plugin, Rank: ranking})
			}
		}

		outputResults(w, &results)
	}
}

func tagsHandler(data *Data) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		for key := range data.Tags {
			fmt.Fprintf(w, "%s\n", key)
		}
	}
}

func fetch(data *Data) {
	for {
		err := fetchData(data)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(1 * time.Hour)
	}
}

func main() {
	data := NewData()
	go fetch(data)

	router := httprouter.New()
	router.GET("/", helpHandler())
	router.GET("/s", allHandler(data))
	router.GET("/s/:search", searchHandler(data))
	router.GET("/t", tagsHandler(data))
	router.GET("/t/:search", searchTagsHandler(data))
	log.Fatal(http.ListenAndServe(":8080", router))
}
