package http

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/algmib/kit"
	"testing"
)

var logger = kit.InitLogger(&kit.LogConfig{Level: kit.InfoLevel})
var logf = func() kit.CLogger {
	return kit.L(logger)
}

type sortConvertTestSuite struct {
	kit.Suite
}

func (s *sortConvertTestSuite) SetupSuite() {
	s.Suite.Init(logf)
}

func TestTagSuite(t *testing.T) {
	suite.Run(t, new(sortConvertTestSuite))
}

func (s *sortConvertTestSuite) Test_ParseSortBy() {
	tests := []struct {
		name       string
		sortString string
		want       []*kit.SortRequest
		wantErr    bool
	}{
		{
			name:       "Empty string",
			sortString: "",
			want:       nil,
		},
		{
			name:       "real example",
			sortString: "reportedAt desc",
			want: []*kit.SortRequest{
				{
					Field: "reportedAt",
					Desc:  true,
				},
			},
		},
		{
			name:       "All ok (without missings)",
			sortString: "field1,field2 desc",
			want: []*kit.SortRequest{
				{
					Field: "field1",
					Desc:  false,
				},
				{
					Field: "field2",
					Desc:  true,
				},
			},
		},
		{
			name:       "All ok (with missings)",
			sortString: "field1 asc first,field2 desc last,field3 asc",
			want: []*kit.SortRequest{
				{
					Field:     "field1",
					Desc:      false,
					NullsLast: false,
				},
				{
					Field:     "field2",
					Desc:      true,
					NullsLast: true,
				},
				{
					Field: "field3",
					Desc:  false,
				},
			},
		},
		{
			name:       "Whitespaces",
			sortString: " field1    asc  , field2 desc  ",
			wantErr:    true,
		},
		{
			name:       "1 field",
			sortString: "field1 asc",
			want: []*kit.SortRequest{
				{
					Field: "field1",
					Desc:  false,
				},
			},
		},
		{
			name:       "1 field only name",
			sortString: "field1",
			want: []*kit.SortRequest{
				{
					Field: "field1",
					Desc:  false,
				},
			},
		},
		{
			name:       "Illegal sort mode",
			sortString: "field1 asc,field2 illegal_mode",
			wantErr:    true,
		},
		{
			name:       "Illegal missing mode",
			sortString: "field1 asc,field2 desc illegal_mode",
			wantErr:    true,
		},
		{
			name:       "Illegal syntax 1",
			sortString: "field1 asc,field2=desc",
			wantErr:    true,
		},
		{
			name:       "Illegal syntax 2",
			sortString: "field1 asc,desc=field2",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		res, err := ParseSortBy(s.Ctx, tt.sortString)
		s.Equal(tt.want, res)
		s.Equal(tt.wantErr, err != nil)
	}
}
