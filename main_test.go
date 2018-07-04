package main

import (
	"reflect"
	"testing"

	"github.com/coreos/go-semver/semver"
)

func stringToVersionSlice(stringSlice []string) []*semver.Version {
	versionSlice := make([]*semver.Version, len(stringSlice))
	for i, versionString := range stringSlice {
		versionSlice[i] = semver.New(versionString)
	}
	return versionSlice
}

func versionToStringSlice(versionSlice []*semver.Version) []string {
	stringSlice := make([]string, len(versionSlice))
	for i, version := range versionSlice {
		stringSlice[i] = version.String()
	}
	return stringSlice
}

func TestLatestVersions(t *testing.T) {
	testCases := []struct {
		name           string
		versionSlice   []string
		expectedResult []string
		minVersion     *semver.Version
	}{
		{
			name:           "case 1",
			versionSlice:   []string{"1.8.11", "1.9.6", "1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.10.1", "1.9.6", "1.8.11"},
			minVersion:     semver.New("1.8.0"),
		},
		{
			name:           "case 2",
			versionSlice:   []string{"1.8.11", "1.9.6", "1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.10.1", "1.9.6"},
			minVersion:     semver.New("1.8.12"),
		},
		{
			name:           "case 3",
			versionSlice:   []string{"1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.10.1"},
			minVersion:     semver.New("1.10.0"),
		},
		{
			name:           "case 4",
			versionSlice:   []string{"2.2.1", "2.2.0"},
			expectedResult: []string{"2.2.1"},
			minVersion:     semver.New("2.2.1"),
		},
		{
			name:           "redundant versions",
			versionSlice:   []string{"2.2.1", "2.2.1", "2.2.1", "2.2.0", "2.2.0"},
			expectedResult: []string{"2.2.1"},
			minVersion:     semver.New("2.2.1"),
		},
		{
			name:           "all less than min ver",
			versionSlice:   []string{"1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{},
			minVersion:     semver.New("1.11.0"),
		},
		{
			name:           "empty version slice",
			versionSlice:   []string{},
			expectedResult: []string{},
			minVersion:     semver.New("0.0.1"),
		},
		{
			name:           "pre-release versions",
			versionSlice:   []string{"1.10.1-rc", "1.9.5-alpha", "1.8.10-beta", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.10.1-rc", "1.9.5", "1.8.10-beta", "1.7.14"},
			minVersion:     semver.New("1.7.0"),
		},
	}

	test := func(versionData []string, expectedResult []string, minVersion *semver.Version, tt *testing.T) {
		stringSlice := versionToStringSlice(LatestVersions(stringToVersionSlice(versionData), minVersion))
		if !reflect.DeepEqual(stringSlice, expectedResult) {
			tt.Errorf("Received %s, expected %s", stringSlice, expectedResult)
		}
	}

	for _, testValues := range testCases {
		t.Run(testValues.name, func(t *testing.T) {
			test(testValues.versionSlice, testValues.expectedResult, testValues.minVersion, t)
		})
	}
}
