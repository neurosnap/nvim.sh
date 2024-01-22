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
	Plugin Plugin `json:"plugin"`
	Rank   int    `json:"rank"`
}

type ResultCollection struct {
	Results *[]Result `json:"results"`
}

type TagCollection struct {
	Tags []string `json:"tags"`
}

type LinksCollection struct {
	Help          string `json:"help"`
	Plugins       string `json:"plugins"`
	SearchPlugins string `json:"search_plugins"`
	Tags          string `json:"tags"`
	TagsSearch    string `json:"tags_search"`
}

type HomeCollection struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Credits     []string        `json:"credits"`
	Links       LinksCollection `json:"links"`
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
	resp, err := http.Get("https://neovimcraft.com/db.json")
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

type Output int64

const (
	PLAIN Output = 0
	JSON         = 1
	HTML         = 2
)

func (o Output) String() string {
	switch o {
	case PLAIN:
		return "plain"
	case JSON:
		return "json"
	case HTML:
		return "html"
	}
	return "unknown"
}

func sanitizeOutputFormat(str string) Output {
	if str == "json" {
		return JSON
	}

	if str == "html" {
		return HTML
	}

	return PLAIN
}

func outputResults(w http.ResponseWriter, results *[]Result, o Output) {
	sort.Sort(ByRank(*results))

	if o == JSON {
		collection := ResultCollection{Results: results}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(collection)
		return
	}

	tbl := printResults(*results, false)
	tbl.WithWriter(w)
	tbl.Print()
}

func helpHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()
		format := sanitizeOutputFormat(r.URL.Query().Get("format"))
		name := "nvim.sh"
		desc := "neovim plugin search from the terminal"
		url := "https://nvim.sh"
		pluginsUrl := fmt.Sprintf("%s/s", url)
		searchPluginsUrl := fmt.Sprintf("%s/s/:search", url)
		tagsUrl := fmt.Sprintf("%s/t", url)
		tagsSearchUrl := fmt.Sprintf("%s/t/:search", url)
		credits := []string{"https://neovimcraft.com", "https://github.com/neurosnap/nvim.sh", "https://bower.sh"}

		if format == JSON {
			collection := HomeCollection{
				Title:       name,
				Description: desc,
				Credits:     credits,
				Links: LinksCollection{
					Help:          url,
					Plugins:       pluginsUrl,
					SearchPlugins: searchPluginsUrl,
					Tags:          tagsUrl,
					TagsSearch:    tagsSearchUrl,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(collection)
			return
		}

		fmt.Fprintf(w, "nvim.sh - neovim plugin search from the terminal\n\n")

		tbl := table.New("api", "description")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt).WithWriter(w)
		tbl.AddRow(url, "help")
		tbl.AddRow(pluginsUrl, "return all plugins in directory")
		tbl.AddRow(searchPluginsUrl, "search for plugin within directory")
		tbl.AddRow(tagsUrl, "list all tags within directory")
		tbl.AddRow(tagsSearchUrl, "search for plugins that exactly match tag within directory")
		tbl.Print()

		fmt.Fprintf(w, "\npowered by: %s\n", credits[0])
		fmt.Fprintf(w, "source: %s\n", credits[1])
		fmt.Fprintf(w, "created by: %s\n", credits[2])
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

		format := sanitizeOutputFormat(r.URL.Query().Get("format"))
		outputResults(w, &results, format)
	}
}

func allHandler(data *Data) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		results := make([]Result, 0, 10)

		for _, plugin := range data.Plugins {
			results = append(results, Result{Plugin: plugin, Rank: 0})
		}

		format := sanitizeOutputFormat(r.URL.Query().Get("format"))
		outputResults(w, &results, format)
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

		format := sanitizeOutputFormat(r.URL.Query().Get("format"))
		outputResults(w, &results, format)
	}
}

func tagsHandler(data *Data) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		format := sanitizeOutputFormat(r.URL.Query().Get("format"))
		if format == JSON {
			tags := make([]string, 0, len(data.Tags))
			for k := range data.Tags {
				tags = append(tags, k)
			}
			collection := TagCollection{Tags: tags}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(collection)
			return
		}

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
	log.Fatal(http.ListenAndServe("0.0.0.0:80", router))
}
