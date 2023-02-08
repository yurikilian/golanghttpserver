package matcher

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_patternMatches(t *testing.T) {
	type args struct {
		path    string
		pattern string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Match /transactions/:id/product/:productId with /transactions/:id/product/:productId",
			args: args{
				path:    "/transactions/1/product/1",
				pattern: "/transactions/:id/product/:productId",
			},
			want: true,
		},
		{
			name: "Not match /transactions/1/product => /transactions/:id/product/:productId",
			args: args{
				path:    "/transactions/1/product",
				pattern: "/transactions/:id/product/:productId",
			},
			want: false,
		},
		{
			name: "Not match transactions/1/product/1 => /transactions/:id/product/:productId",
			args: args{
				path:    "transactions/1/product/1",
				pattern: "/transactions/:id/product/:productId",
			},
			want: false,
		},
		{
			name: "Match /transactions/1/product => /transactions/:id/product",
			args: args{
				path:    "/transactions/1/product",
				pattern: "/transactions/:id/product",
			},
			want: true,
		},
		{
			name: "Match /transactions => /transactions ",
			args: args{
				path:    "/transactions",
				pattern: "/transactions",
			},
			want: true,
		},
		{
			name: "Match / => / ",
			args: args{
				path:    "/",
				pattern: "/",
			},
			want: true,
		},
		{
			name: "Match /1/product/1  with /:id/product/:productId",
			args: args{
				path:    "/1/product/1",
				pattern: "/:id/product/:productId",
			},
			want: true,
		},
		{
			name: "Match categories/1/products?name=test",
			args: args{
				path:    "/categories/1/products?name=test",
				pattern: "/categories/:id/products",
			},
			want: true,
		},
		{
			name: "Match /products?name=test/1/categories",
			args: args{
				path:    "/products?name=test/1/categories",
				pattern: "/products/:id/categories",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, MatchPath(strings.Split(tt.args.path, "/"), tt.args.pattern), "patternMatches(%v, %v)", tt.args.path, tt.args.pattern)
		})
	}
}
