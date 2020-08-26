package scale_rules

import (
	"testing"
)

func TestRule_result(t *testing.T) {
	type args struct {
		value float64
	}
	tests := []struct {
		name    string
		r       Rule
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Smaller 200 = 1",
			r:    "<200?1:1.1",
			args: args{
				value: 10,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Greater 200 = 1.1",
			r:    "<200?1:1.1",
			args: args{
				value: 1000,
			},
			want:    1.1,
			wantErr: false,
		}, {
			name: "Smaller 200 = 2",
			r:    ">200?3.5:2.0",
			args: args{
				value: 10,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "Greater 200 = 3.5",
			r:    ">200?3.5:2.0",
			args: args{
				value: 1000,
			},
			want:    3.5,
			wantErr: false,
		}, {
			name: "Equal 15 = 3",
			r:    "=15?3:2.3",
			args: args{
				value: 15,
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "not Equal 15 = 2.3",
			r:    "=15?3:2.3",
			args: args{
				value: 1000,
			},
			want:    2.3,
			wantErr: false,
		}, {
			name: "NOT Equal 100 = 10",
			r:    "!100?10:5",
			args: args{
				value: 15,
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "Equal 100 = 5",
			r:    "!100?10:5",
			args: args{
				value: 100,
			},
			want:    5,
			wantErr: false,
		}, {
			name: "Between 100-200 = 15",
			r:    "@100-200?15:89",
			args: args{
				value: 100,
			},
			want:    15,
			wantErr: false,
		},
		{
			name: "NOT Between 100-200 = 89",
			r:    "@100-200?15:89",
			args: args{
				value: 250,
			},
			want:    89,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err, _ := tt.r.result(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("result() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("result() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRules_resultForVal(t *testing.T) {
	ruls := NewScaleRules()
	rules, _ := ruls.Parse("[\"<200?1:1.1\",\">200?3.5:2.0\",\"=15?3:2.3\",\"!100?10:5\",\"@100-200?15:89\"]")
	tr := true

	type args struct {
		value   float64
		combine CombineFunction
		only *bool
	}
	tests := []struct {
		name string
		r    Rules
		args args
		want float64
	}{
		{
			name: "CombineLastOne not between 100-200",
			r:    *rules,
			args: args{
				value:   10,
				combine: CombineLastOne,
			},
			want: 89,
		}, {
			name: "CombineSum (10) = 104.3",
			r:    *rules,
			args: args{
				value:   10,
				combine: CombineSum,
			},
			want: 104.3,
		},{
			name: "CombineSum (10) = 104.3",
			r:    *rules,
			args: args{
				value:   10,
				combine: CombineSum,
				only: &tr,
			},
			want: 11,
		}, {
			name: "CombineFirstOne smaller 200 = 1",
			r:    *rules,
			args: args{
				value:   10,
				combine: CombineFirstOne,
			},
			want: 1,
		}, {
			name: "CombineIntersections = 1",
			r:    *rules,
			args: args{
				value:   10,
				combine: CombineIntersections,
			},
			want: 1,
		}, {
			name: "CombineAvg = 20.86",
			r:    *rules,
			args: args{
				value:   10,
				combine: CombineAvg,
			},
			want: 20.86,
		}, {
			name: "CombineAvg = 20.86",
			r:    *rules,
			args: args{
				value:   10,
				combine: CombineAvg,
				only: &tr,
			},
			want: 5.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.resultForVal(tt.args.value, tt.args.combine, tt.args.only); got != tt.want {
				t.Errorf("resultForVal() = %v, want %v", got, tt.want)
			}
		})
	}
}
