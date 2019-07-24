package helper

import "testing"

func TestGroupDetection(t *testing.T){
	type args struct{
		url string
	}
	tests := []struct{
		name string
		args args
		want string
	}{
		{"success",args{"git@gitlab.warungpintar.co:back-end/warbot-go.git"}, "back-end"},
		{"success",args{"git@gitlab.warungpintar.co:wartech/warbot-go.git"}, "wartech"},
		{"success",args{"git@gitlab.warungpintar.co:erp-tech/warbot-go.git"}, "erp-tech"},
	}
	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			got := GroupDetection(tt.args.url)
			if got != tt.want {
				t.Errorf("GroupDetection() = %v, want %v", got,tt.want)
			}
		})
	}
}