package cmd

import "testing"

func Test_GetAbsLocalCnfPath(t *testing.T) {
	tests := []struct {
		inputPath  string
		wantedPath string
	}{
		{"~/test", "/Users/hjiang/test"},
		{"/User/hjiang/heng", "/User/hjiang/heng"},
	}

	for _, tt := range tests {
		out, err := GetAbsLocalCnfPath(tt.inputPath)
		if err != nil {
			t.Error(err)
		}

		if out != tt.wantedPath {
			t.Errorf("get %v, want %v, mismatch!", out, tt.wantedPath)
		}
	}
}

func Test_copyTo(t *testing.T) {
	tests := []struct {
		source string
		dest   string
		suffix string
	}{
		{"../test/source", "../test/dest", ".txt"},
	}

	for _, tt := range tests {
		err := copyTo(tt.source, tt.dest, tt.suffix)
		if err != nil {
			t.Error(err)
		}
	}
}
