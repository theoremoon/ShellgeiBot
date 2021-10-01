package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeTweetable(t *testing.T) {
	s := "hogehogehogehogehoge"
	e := "hogehogehogehogehoge"
	r := makeTweetable(s)
	if r != e {
		t.Fatalf("expected %s but %s", e, r)
	}

	s = strings.Repeat("y", 280)
	e = s
	r = makeTweetable(s)
	if r != e {
		t.Fatalf("expected %s but %s", e, r)
	}

	s = strings.Repeat("y", 281)
	e = strings.Repeat("y", 280)
	r = makeTweetable(s)
	if r != e {
		t.Fatalf("expected %s but %s", e, r)
	}

	s = "あめんぼあかいなあいうえお　うきもにこえびもおよいでる　かきのきくりのきかきくけこ　きつつきこつこつかれけやき　ささげにすをかけさしすせそ　そそそそそそそそそそそそそ　たてちつてとたとたちつてと　とてとてたったととびたった　なめくじぬめってなにぬねの　なんどにぬめってなにねばる　はらひれほろろろ"
	e = "あめんぼあかいなあいうえお　うきもにこえびもおよいでる　かきのきくりのきかきくけこ　きつつきこつこつかれけやき　ささげにすをかけさしすせそ　そそそそそそそそそそそそそ　たてちつてとたとたちつてと　とてとてたったととびたった　なめくじぬめってなにぬねの　なんどにぬめってなにねばる　"
	r = makeTweetable(s)
	if r != e {
		t.Fatalf("expected %s but %s", e, r)
	}

	s = "2.718281828459045235360287471352662497757247093699959574966967627724\\ 07663035354759457138217852516642742746639193200305992181741359662904\\2.718281828459045235360287471352662497757247093699959574966967627724\\ 07663035354759457138217852516642742746639193200305992181741359662904\\"
	e = s
	r = makeTweetable(s)
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
		desc        string // テストの目的、理由
		expect      string
		text        string
		hashtags    tweetEntitiesHashtags
		extHashtags tweetEntitiesHashtags // Empty
		tags        []string
	}
	testDatas := []TestData{
		{
			desc:   "#シェル芸 タグだけが削除される",
			expect: "echo test \n#シェル芸2 #shellgei",
			text:   "echo test #シェル芸\n#シェル芸2 #shellgei",
			hashtags: tweetEntitiesHashtags{
				{Indices: []int{10, 15}, Text: "シェル芸"},
				{Indices: []int{16, 22}, Text: "シェル芸2"},
				{Indices: []int{23, 32}, Text: "shellgei"},
			},
			tags: tags,
		},
		{
			desc:   "tagsに存在するものはすべて削除される。前後の空白は削除される。",
			expect: "echo シェル芸",
			text:   " echo シェル芸 #シェル芸 #ゆるシェル #危険シェル芸 ",
			hashtags: tweetEntitiesHashtags{
				{Indices: []int{11, 16}, Text: "シェル芸"},
				{Indices: []int{17, 23}, Text: "ゆるシェル"},
				{Indices: []int{24, 31}, Text: "危険シェル芸"},
			},
			tags: tags,
		},
		{
			desc:   "削除対象のタグが存在しないときはそのまま返す",
			expect: "echo test #shellgei #シェルぎえ",
			text:   "echo test #shellgei #シェルぎえ",
			hashtags: tweetEntitiesHashtags{
				{Indices: []int{10, 19}, Text: "shellgei"},
				{Indices: []int{20, 26}, Text: "シェルぎえ"},
			},
			tags: tags,
		},
		{
			desc:   "シェル芸っぽいというだけのタグは消えない",
			expect: "echo test #シェル芸a #bシェル芸 # シェル芸",
			text:   "echo test #シェル芸a #bシェル芸 # シェル芸 #シェル芸",
			hashtags: tweetEntitiesHashtags{
				{Indices: []int{10, 16}, Text: "シェル芸a"},
				{Indices: []int{17, 23}, Text: "bシェル芸"},
				{Indices: []int{31, 36}, Text: "シェル芸"},
			},
			tags: tags,
		},
		{
			desc:   "同じタグが付与されている場合もすべて削除される",
			expect: "echo test",
			text:   "echo test #シェル芸 #シェル芸 #シェル芸",
			hashtags: tweetEntitiesHashtags{
				{Indices: []int{10, 15}, Text: "シェル芸"},
				{Indices: []int{16, 22}, Text: "シェル芸"},
				{Indices: []int{22, 27}, Text: "シェル芸"},
			},
			tags: tags,
		},
		{
			desc:   "タグが20個 extなし",
			expect: "echo test  #シェル芸1 #シェル芸2 #シェル芸3 #シェル芸4 #シェル芸5 #シェル芸6 #シェル芸7 #シェル芸8 #シェル芸9 #シェル芸10 #シェル芸11 #シェル芸12 #シェル芸13 #シェル芸14 #シェル芸15 #シェル芸16 #シェル芸17 #シェル芸18 #シェル芸19 #シェル芸20",
			text:   "echo test #シェル芸 #シェル芸1 #シェル芸2 #シェル芸3 #シェル芸4 #シェル芸5 #シェル芸6 #シェル芸7 #シェル芸8 #シェル芸9 #シェル芸10 #シェル芸11 #シェル芸12 #シェル芸13 #シェル芸14 #シェル芸15 #シェル芸16 #シェル芸17 #シェル芸18 #シェル芸19 #シェル芸20",
			hashtags: tweetEntitiesHashtags{
				{Indices: []int{10, 15}, Text: "シェル芸"},
				{Indices: []int{16, 22}, Text: "シェル芸1"},
				{Indices: []int{17, 23}, Text: "シェル芸2"},
				{Indices: []int{24, 30}, Text: "シェル芸3"},
				{Indices: []int{31, 37}, Text: "シェル芸4"},
				{Indices: []int{38, 44}, Text: "シェル芸5"},
				{Indices: []int{45, 51}, Text: "シェル芸6"},
				{Indices: []int{52, 58}, Text: "シェル芸7"},
				{Indices: []int{59, 65}, Text: "シェル芸8"},
				{Indices: []int{66, 72}, Text: "シェル芸9"},
				{Indices: []int{73, 80}, Text: "シェル芸10"},
				{Indices: []int{81, 88}, Text: "シェル芸11"},
				{Indices: []int{89, 96}, Text: "シェル芸12"},
				{Indices: []int{97, 104}, Text: "シェル芸13"},
				{Indices: []int{105, 112}, Text: "シェル芸14"},
				{Indices: []int{113, 120}, Text: "シェル芸15"},
				{Indices: []int{121, 128}, Text: "シェル芸16"},
				{Indices: []int{129, 136}, Text: "シェル芸17"},
				{Indices: []int{137, 142}, Text: "シェル芸18"},
				{Indices: []int{143, 150}, Text: "シェル芸19"},
				{Indices: []int{151, 158}, Text: "シェル芸20"},
			},
			tags: tags,
		},
	}
	for _, v := range testDatas {
		got := removeTags(v.text, v.hashtags, v.extHashtags, v.tags)
		assert.Equal(t, v.expect, got, v.desc)
	}
}
