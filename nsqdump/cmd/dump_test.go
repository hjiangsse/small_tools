package cmd_test

import (
	"fmt"
	"nsqdump/cmd"
	"sort"
	"testing"
)

func Test_GetFileName(t *testing.T) {
	var tests = []struct {
		path     string
		topic    string
		filenum  uint64
		expected string
	}{
		{".", "test", 1, "test.diskqueue.000001.dat"},
		{"/usr/home", "test", 1, "/usr/home/test.diskqueue.000001.dat"},
		{"/usr/home", "test", 1000, "/usr/home/test.diskqueue.001000.dat"},
	}

	for _, tt := range tests {
		actual := cmd.GetDataFileName(tt.path, tt.topic, tt.filenum)
		if actual != tt.expected {
			t.Errorf("expected %s, actual %s", tt.expected, actual)
		}
	}
}

func Test_GetAllTopics(t *testing.T) {
	var tests = []struct {
		filepath string
		expected []string
	}{
		{"/home/hjiang/go/src/nsqdump/test/gettopics", []string{"test1", "test2", "test3"}},
	}

	for _, tt := range tests {
		actual, err := cmd.GetAllTopics(tt.filepath)
		if err != nil {
			t.Error(err)
		}

		sort.Slice(actual, func(i, j int) bool {
			return actual[i] < actual[j]
		})

		actualStr := fmt.Sprintf("%v", actual)
		wantedStr := fmt.Sprintf("%v", tt.expected)
		if actualStr != wantedStr {
			t.Errorf("expected: %v, actual %v", tt.expected, actual)
		}
	}
}

func Test_DumpSpecificTopic(t *testing.T) {
	var tests = []struct {
		datapath string
		outpath  string
		topic    string
	}{
		{"/home/hjiang/Nsq/nsq/build/", "/home/hjiang/Data1/", "test1"},
	}

	for _, tt := range tests {
		err := cmd.DumpSpecificTopic(tt.datapath, tt.outpath, tt.topic)
		if err != nil {
			t.Error(err)
		}
	}
}
