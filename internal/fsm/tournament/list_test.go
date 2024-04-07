package tournament

import (
	"github.com/oke11o/go-telegram-bot/internal/model"
	"testing"
)

func TestListTournament_PrintTournaments(t *testing.T) {

	tests := []struct {
		name        string
		tournaments []model.Tournament

		want string
	}{

		{
			name: "Test PrintTournaments",
			tournaments: []model.Tournament{
				{Title: "Test", Date: "2021-10-10"},
				{Title: "Test2", Date: "2021-10-11"},
				{Title: "Test3", Date: "2021-10-12"},
			},
			want: "List of opened tournaments:\n- Test [2021-10-10]\n- Test2 [2021-10-11]\n- Test3 [2021-10-12]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ListTournament{}
			got := m.PrintTournaments(tt.tournaments)
			if got != tt.want {
				t.Errorf("PrintTournaments() = %v, want %v", got, tt.want)
			}
		})
	}
}
