package cmd

import "testing"

func Test_DecodeAndDumpNsqBackFile(t *testing.T) {
	var tests = []struct {
		sourcefile string
		destfile   string
	}{
		{"/home/hjiang/Nsq/nsq/build/test1.diskqueue.000000.dat", "/home/hjiang/Nsq/nsq/build/test1.diskqueue.000000.json"},
	}

	for _, tt := range tests {
		err := DecodeAndDumpNsqBackFile(tt.sourcefile, tt.destfile)
		if err != nil {
			t.Error(err)
		}
	}
}
