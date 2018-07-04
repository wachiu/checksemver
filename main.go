package main

import (
	"context"
	"fmt"
	"sort"

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
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, "kubernetes", "kubernetes", opt)
	if err != nil {
		panic(err) // TODO:is this really a good way?
	}
	minVersion := semver.New("1.8.0")
	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	versionSlice := LatestVersions(allReleases, minVersion)

	fmt.Printf("latest versions of kubernetes/kubernetes: %s", versionSlice)
}
