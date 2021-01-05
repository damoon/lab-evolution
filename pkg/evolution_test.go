package evolution

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSimpleGenExchange(t *testing.T) {
	type args struct {
		r      *rand.Rand
		mother Lifeform
		father Lifeform
	}
	tests := []struct {
		name string
		args args
		want Genome
	}{
		{
			name: "a few bits from father",
			args: args{
				r:      rand.New(rand.NewSource(0)),
				mother: Lifeform{genes: []byte{0, 0, 0}},
				father: Lifeform{genes: []byte{255, 255, 255}},
			},
			want: []byte{0, 0, 63},
		},
		{
			name: "a few bits from mother",
			args: args{
				r:      rand.New(rand.NewSource(4)),
				mother: Lifeform{genes: []byte{0, 0, 0}},
				father: Lifeform{genes: []byte{255, 255, 255}},
			},
			want: []byte{7, 255, 255},
		},
		{
			name: "a few of both",
			args: args{
				r:      rand.New(rand.NewSource(6)),
				mother: Lifeform{genes: []byte{0, 0, 0}},
				father: Lifeform{genes: []byte{255, 255, 255}},
			},
			want: []byte{0, 7, 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleGenExchange(tt.args.r, tt.args.mother, tt.args.father); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleGenExchange() = %v, want %v", got, tt.want)
			}
		})
	}
}
