package helper

import "testing"

// TestConvertTimeToZone for unit testing
func TestConvertTimeToZone(t *testing.T) {
	type args struct {
		timeFromGitlab string
		strZone        string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"success", args{"2019-06-16 13:18:36 UTC", "Asia/Jakarta"}, "16 Jun 19 20:18", false},
		{"success", args{"2019-06-16 13:18:36 UTC", "Asia/Singapore"}, "16 Jun 19 21:18", false},
		{"invalid_timezone", args{"2019-06-16 13:18:36 UTC", "Asia/anywhere"}, "16 Jun 19 13:18", true},
		{"invalid_datetime", args{"2019-06-16XXX 13:18:36 UTC", "Asia/anywhere"}, "2019-06-16XXX 13:18:36 UTC", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertTimeToZone(tt.args.timeFromGitlab, tt.args.strZone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertTimeToZone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertTimeToZone() = %v, want %v", got, tt.want)
			}
		})
	}
}
