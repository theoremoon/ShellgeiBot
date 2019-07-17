package main

import (
	"strings"
	"testing"

	"github.com/ChimeraCoder/anaconda"
	"github.com/stretchr/testify/assert"
)

func Test_MakeTweetable(t *testing.T) {
	s := "hogehogehogehogehoge"
	e := "hogehogehogehogehoge"
	r := MakeTweetable(s)
	if r != e {
		t.Fatalf("expected %s but %s", e, r)
	}

	s = strings.Repeat("y", 280)
	e = s
	r = MakeTweetable(s)
	if r != e {
		t.Fatalf("expected %s but %s", e, r)
	}

	s = strings.Repeat("y", 281)
	e = strings.Repeat("y", 280)
	r = MakeTweetable(s)
	if r != e {
		t.Fatalf("expected %s but %s", e, r)
	}

	s = "あめんぼあかいなあいうえお　うきもにこえびもおよいでる　かきのきくりのきかきくけこ　きつつきこつこつかれけやき　ささげにすをかけさしすせそ　そそそそそそそそそそそそそ　たてちつてとたとたちつてと　とてとてたったととびたった　なめくじぬめってなにぬねの　なんどにぬめってなにねばる　はらひれほろろろ"
	e = "あめんぼあかいなあいうえお　うきもにこえびもおよいでる　かきのきくりのきかきくけこ　きつつきこつこつかれけやき　ささげにすをかけさしすせそ　そそそそそそそそそそそそそ　たてちつてとたとたちつてと　とてとてたったととびたった　なめくじぬめってなにぬねの　なんどにぬめってなにねばる　"
	r = MakeTweetable(s)
	if r != e {
		t.Fatalf("expected %s but %s", e, r)
	}

	s = "2.718281828459045235360287471352662497757247093699959574966967627724\\ 07663035354759457138217852516642742746639193200305992181741359662904\\2.718281828459045235360287471352662497757247093699959574966967627724\\ 07663035354759457138217852516642742746639193200305992181741359662904\\"
	e = s
	r = MakeTweetable(s)
	if r != e {
		t.Fatalf("expected %s but %s", e, r)
	}

}

func TestRemoveTags(t *testing.T) {
	tags := []string{
		"シェル芸",
		"ゆるシェル",
		"危険シェル芸",
	}
	type TestData struct {
		desc     string // テストの目的、理由
		expect   string
		text     string
		entities anaconda.Entities
		tags     []string
	}
	testDatas := []TestData{
		{
			desc:   "#シェル芸 タグだけが削除される",
			expect: "echo test \n#シェル芸2 #shellgei",
			text:   "echo test #シェル芸\n#シェル芸2 #shellgei",
			entities: anaconda.Entities{
				Hashtags: []struct {
					Indices []int
					Text    string
				}{
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{10, 15},
						Text:    "シェル芸",
					},
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{16, 22},
						Text:    "シェル芸2",
					},
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{23, 32},
						Text:    "shellgei",
					},
				},
			},
			tags: tags,
		},
		{
			desc:   "tagsに存在するものはすべて削除される。前後の空白は削除される。",
			expect: "echo シェル芸",
			text:   " echo シェル芸 #シェル芸 #ゆるシェル #危険シェル芸 ",
			entities: anaconda.Entities{
				Hashtags: []struct {
					Indices []int
					Text    string
				}{
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{11, 16},
						Text:    "シェル芸",
					},
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{17, 23},
						Text:    "ゆるシェル",
					},
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{24, 31},
						Text:    "危険シェル芸",
					},
				},
			},
			tags: tags,
		},
		{
			desc:   "削除対象のタグが存在しないときはそのまま返す",
			expect: "echo test #shellgei #シェルぎえ",
			text:   "echo test #shellgei #シェルぎえ",
			entities: anaconda.Entities{
				Hashtags: []struct {
					Indices []int
					Text    string
				}{
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{10, 19},
						Text:    "shellgei",
					},
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{20, 26},
						Text:    "シェルぎえ",
					},
				},
			},
			tags: tags,
		},
		{
			desc:   "シェル芸っぽいというだけのタグは消えない",
			expect: "echo test #シェル芸a #bシェル芸 # シェル芸",
			text:   "echo test #シェル芸a #bシェル芸 # シェル芸 #シェル芸",
			entities: anaconda.Entities{
				Hashtags: []struct {
					Indices []int
					Text    string
				}{
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{10, 16},
						Text:    "シェル芸a",
					},
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{17, 23},
						Text:    "bシェル芸",
					},
					struct {
						Indices []int
						Text    string
					}{
						Indices: []int{31, 36},
						Text:    "シェル芸",
					},
				},
			},
			tags: tags,
		},
	}
	for _, v := range testDatas {
		got := removeTags(v.text, v.entities, v.tags)
		assert.Equal(t, v.expect, got, v.desc)
	}
}
