package sdees

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type GithubCommitsJSON []struct {
	Sha    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			Sha string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		URL          string `json:"url"`
		CommentCount int    `json:"comment_count"`
	} `json:"commit"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Author      struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Committer struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"committer"`
	Parents []struct {
		Sha     string `json:"sha"`
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
	} `json:"parents"`
}

type GithubJson struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}

func CheckNewVersion(dir string, version string, lastcommit string, osType string) {
	logger.Debug("Current executable path: %s", dir)
	if version == "dev" {
		updateDevVersion(dir, version, lastcommit, osType)
	} else {
		updateDownloadVersion(dir, version, lastcommit, osType)
	}
}

func updateDevVersion(dir string, version string, lastcommit string, osType string) {
	logger.Debug("Updating dev version of sdees")
	url := "https://api.github.com/repos/schollz/sdees/commits"
	r, err := http.Get(url)
	if err != nil {
		logger.Debug("Couldn't call Github API for getting new date")
		return
	}
	defer r.Body.Close()
	var j GithubCommitsJSON
	err = json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		logger.Debug("Couldn't decode Github API")
		return
	}

	if len(j) == 0 {
		logger.Debug("No data form Github!")
		return
	}
	currentCommit, err := ParseDate(strings.Replace(lastcommit, "'", "", -1))
	if err != nil {
		logger.Debug("Couldn't parse Github API Commit date")
		return
	}
	logger.Debug("Github: %s, Current: %s", j[0].Commit.Author.Date.String(), currentCommit.String())
	if currentCommit.Sub(j[0].Commit.Author.Date).Hours() < 0 {
		fmt.Println("New version of sdees available! Run\n\n\tgo get -u github.com/schollz/sdees\n\nto download.")
	}
}

func updateDownloadVersion(dir string, version string, lastcommit string, osType string) {
	newVersion, versionName := checkGithub(version)
	if !newVersion {
		logger.Debug("Current version is up to date: %s / %s", version, versionName)
		return
	}
	var yesnoall string
	fmt.Printf("\nVersion %s is available. Download? (y/n) ", versionName)
	fmt.Scanln(&yesnoall)
	if yesnoall == "n" {
		return
	}
	downloadVersion := versionName
	downloadName := "sdees_" + osType + ".zip"
	os.Remove("sdees_" + osType + ".zip")
	fmt.Printf("\nDownloading %s/%s...", downloadVersion, downloadName)
	DownloadFile(downloadName, "https://github.com/schollz/sdees/releases/download/"+downloadVersion+"/"+downloadName)

	logger.Debug("Unzipping new version")
	Unzip(downloadName, "./")

	logger.Debug("Cleaning...")
	os.Remove(downloadName)
	fmt.Printf("\n\nsdees Version %s installed!\n", versionName)
	os.Exit(0)
}

func checkGithub(version string) (bool, string) {
	url := "https://api.github.com/repos/schollz/sdees/releases/latest"
	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	var j GithubJson
	err = json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		log.Fatal(err)
	}
	newVersion := j.TagName
	versions := strings.Split(newVersion, ".")
	if len(versions) != 3 {
		return false, ""
	}
	majorMinorWeb := []int{}
	for i := 0; i < 3; i++ {
		i, _ := strconv.Atoi(versions[i])
		majorMinorWeb = append(majorMinorWeb, i)
	}

	versions = strings.Split(version, ".")
	if len(versions) != 3 {
		return false, ""
	}
	majorMinor := []int{}
	for i := 0; i < 3; i++ {
		i, _ := strconv.Atoi(versions[i])
		majorMinor = append(majorMinor, i)
	}

	newVersionAvailable := false
	for i := range majorMinor {
		if majorMinorWeb[i] > majorMinor[i] {
			newVersionAvailable = true
			break
		}
	}

	return newVersionAvailable, newVersion
}
