package main

import (
	"strings"
	"testing"
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
