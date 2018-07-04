package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	versionSlice := filterVersions(releases, minVersion)
	versionSlice = onlyMaxPatch(versionSlice)
	sort.Slice(versionSlice, func(i, j int) bool {
		if versionSlice[i].Compare(*versionSlice[j]) == 1 {
			return true
		}
		return false
	})
	return versionSlice
}

// filterVersions remove versions which are less than minVersion
func filterVersions(input []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var ret []*semver.Version
	for _, v := range input {
		if !v.LessThan(*minVersion) {
			ret = append(ret, v)
		}
	}
	return ret
}

// onlyMaxPatch keep only versions with the max patch
func onlyMaxPatch(input []*semver.Version) []*semver.Version {
	maxVerMap := make(map[string]*semver.Version)
	for _, v := range input {
		key := fmt.Sprintf("%d.%d", v.Major, v.Minor)
		curMax, ok := maxVerMap[key]
		if !ok {
			maxVerMap[key] = v
			continue
		}
		if curMax.LessThan(*v) {
			maxVerMap[key] = v
		}
	}
	var ret []*semver.Version
	for _, v := range maxVerMap {
		ret = append(ret, v)
	}
	return ret
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("missing file path in first argument")
	}

	filepath := args[1]

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("fail to open input file, err: %+v", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	var passHeader bool
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("fail to read input file, err: %+v", err)
		}
		if !passHeader {
			passHeader = true
			continue
		}

		// Github
		client := github.NewClient(nil)
		ctx := context.Background()
		opt := &github.ListOptions{PerPage: 100}
		repoInfo := strings.SplitN(record[0], "/", 2)
		var allRepoReleases []*github.RepositoryRelease
		for {
			releases, resp, err := client.Repositories.ListReleases(ctx, repoInfo[0], repoInfo[1], opt)
			if err != nil {
				log.Fatalf("fail to list releases from GitHub API, err: %+v", err)
			}
			allRepoReleases = append(allRepoReleases, releases...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
		minVersion := semver.New(record[1])
		var allReleases []*semver.Version
		for _, release := range allRepoReleases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			v, err := semver.NewVersion(versionString)
			if err != nil {
				log.Printf("Invalid version string, repo: %s, version string: %s, err: %+v", record[0], versionString, err)
				continue
			}
			allReleases = append(allReleases, v)
		}
		versionSlice := LatestVersions(allReleases, minVersion)

		fmt.Printf("latest versions of %s: %s\n", record[0], versionSlice)
	}

}
