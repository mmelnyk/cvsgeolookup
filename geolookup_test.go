package cvsgeolookup

import (
	"strings"
	"testing"
)

func TestBasicInit(t *testing.T) {
	if _, err := New(); err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}
}

func TestBasicOptions(t *testing.T) {
	opts := []Option{
		WithBeginName("startA"),
		WithEndName("endA"),
		WithLantitudeName("lan"),
		WithLongtitudeName("lon"),
		WithSkipName("skip"),
		WithSkipValue("skip-val"),
		WithCommaValue('+'),
		WithCommentValue('/'),
		WithMetrics(nil),
	}

	eng, err := New(opts...)

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	if eng.options.fieldNameBegin != "startA" {
		t.Fatalf("Expected startA set in option, but got %v", eng.options.fieldNameBegin)
	}

	if eng.options.fieldNameEnd != "endA" {
		t.Fatalf("Expected endA set in option, but got %v", eng.options.fieldNameEnd)
	}

	if eng.options.fieldNameLantitude != "lan" {
		t.Fatalf("Expected lan set in option, but got %v", eng.options.fieldNameLantitude)
	}

	if eng.options.fieldNameLongtitude != "lon" {
		t.Fatalf("Expected lon set in option, but got %v", eng.options.fieldNameLongtitude)
	}

	if eng.options.fieldNameSkip != "skip" {
		t.Fatalf("Expected skip set in option, but got %v", eng.options.fieldNameSkip)
	}

	if eng.options.skipValue != "skip-val" {
		t.Fatalf("Expected skip-val set in option, but got %v", eng.options.skipValue)
	}

	if eng.options.commaRune != '+' {
		t.Fatalf("Expected + set in option, but got %v", eng.options.commaRune)
	}

	if eng.options.commentRune != '/' {
		t.Fatalf("Expected / set in option, but got %v", eng.options.commentRune)
	}

	if eng.options.metrics != nil {
		t.Fatalf("Expected nil(metrics) set in option, but got %v", eng.options.metrics)
	}
}

func TestBasicLoad(t *testing.T) {
	in := `start,end,lantitude,longtitude
10.0.0.0,10.255.255.255,-1.0,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`
	eng, err := New()

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	err = eng.Load(strings.NewReader(in))

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	if len(eng.records) != 2 {
		t.Fatalf("Expected 2 records, but got %v", len(eng.records))
	}

}

func TestBadLoad(t *testing.T) {
	testcase := []struct {
		in  string
		exp error
	}{

		{in: `st,end,lantitude,longtitude
10.0.0.0,10.255.255.255,-1.0,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`,
			exp: ErrNoBeginField},
		{in: `start,stop,lantitude,longtitude
10.0.0.0,10.255.255.255,-1.0,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`,
			exp: ErrNoEndField},
		{in: `start,end,lan,longtitude
10.0.0.0,10.255.255.255,-1.0,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`,
			exp: ErrNoLantitudeField},
		{in: `start,end,lantitude,long
10.0.0.0,10.255.255.255,-1.0,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`,
			exp: ErrNoLongtitudeField},
		{in: `start,end,lantitude,longtitude
10.0.0.0,0.255.255.255,-1.0,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`,
			exp: ErrIncorrectSegment},
		{in: `start,end,lantitude,longtitude
1,10.255.255.255,-1.0,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`,
			exp: ErrWrongIPFormat},
		{in: `start,end,lantitude,longtitude
10.0.0.0,10,-1.0,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`,
			exp: ErrWrongIPFormat},
	}

	eng, err := New()

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	for i, v := range testcase {
		err = eng.Load(strings.NewReader(v.in))

		if err != v.exp {
			t.Fatalf("(case %v) Expected %v, but got %v", i, v.exp, err)
		}
	}
}

func TestBadLoadSpecial(t *testing.T) {
	testcase := []struct {
		in string
	}{
		{in: `start,end,lantitude,longtitude
10.0.0.0,10.255.255.255,test,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`},
		{in: `start,end,lantitude,longtitude
10.0.0.0,10.255.255.255,-1.0,-1.0
20.0.0.0,20.255.255.255,1.0,test
`},
		{in: `start,end,lantitude,longtitude
10.0.0.0,10.255.255.255,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`},
		{in: `start,end,lantitude,longtitude
10.0.0.0,10.255.255.255,-1.0,-1.0,-1.0
20.0.0.0,20.255.255.255,1.0,1.0
`},
	}

	eng, err := New()

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	err = eng.Load(nil)
	if err != ErrReadInterfaceRequired {
		t.Fatalf("Expected ErrReadInterfaceRequired, but got %v", err)
	}

	for i, v := range testcase {
		err = eng.Load(strings.NewReader(v.in))

		if err == nil {
			t.Fatalf("(case %v) Expected error, but got nil", i)
		}
	}
}

func TestBasicLookup(t *testing.T) {
	in := `start,end,lantitude,longtitude,skip
10.0.0.0,10.49.255.255,-1.0,-1.0,yes
10.50.0.0,10.99.255.255,-2.0,-2.0,yes
10.100.0.0,10.149.255.255,-3.0,-3.0,yes
10.150.0.0,10.199.255.255,-4.0,-4.0,yes
10.200.0.0,10.255.255.255,-5.0,-5.0,yes
20.0.0.0,20.255.255.255,1.0,1.0,no
`
	eng, err := New()

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	err = eng.Load(strings.NewReader(in))

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	lant, long, err := eng.Lookup("10.1.1.1")

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	if lant != -1 {
		t.Fatalf("Expected lantitude -1, but got %v", lant)
	}

	if long != -1 {
		t.Fatalf("Expected longtitude -1, but got %v", long)
	}

	lant, long, err = eng.Lookup("10.177.1.1")

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	if lant != -4 {
		t.Fatalf("Expected lantitude -1, but got %v", lant)
	}

	if long != -4 {
		t.Fatalf("Expected longtitude -1, but got %v", long)
	}

	lant, long, err = eng.Lookup("20.1.1.1")

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	if lant != 1 {
		t.Fatalf("Expected lantitude 1, but got %v", lant)
	}

	if long != 1 {
		t.Fatalf("Expected longtitude 1, but got %v", long)
	}

	_, _, err = eng.Lookup("30.1.1.1")

	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound, but got %v", err)
	}
}

func TestBadLookup(t *testing.T) {
	in := `start,end,lantitude,longtitude,skip
10.0.0.0,10.49.255.255,-1.0,-1.0,yes
10.50.0.0,10.99.255.255,-2.0,-2.0,yes
10.100.0.0,10.149.255.255,-3.0,-3.0,yes
10.150.0.0,10.199.255.255,-4.0,-4.0,yes
10.200.0.0,10.255.255.255,-5.0,-5.0,yes
20.0.0.0,20.255.255.255,1.0,1.0,no
`
	eng, err := New()

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	_, _, err = eng.Lookup("10.1.1.1")

	if err != ErrNotInitialized {
		t.Fatalf("Expected ErrNotInitialized, but got %v", err)
	}

	err = eng.Load(strings.NewReader(in))

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	_, _, err = eng.Lookup("10.0")

	if err != ErrWrongIPFormat {
		t.Fatalf("Expected ErrWrongIPFormat, but got %v", err)
	}

	_, _, err = eng.Lookup("test")

	if err != ErrWrongIPFormat {
		t.Fatalf("Expected ErrWrongIPFormat, but got %v", err)
	}

}

func TestSkipLoad(t *testing.T) {
	in := `start,end,lantitude,longtitude,skip
10.0.0.0,10.255.255.255,-1.0,-1.0,yes
20.0.0.0,20.255.255.255,1.0,1.0,no
`
	eng, err := New(
		WithSkipName("skip"),
		WithSkipValue("yes"),
	)

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	err = eng.Load(strings.NewReader(in))

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	_, _, err = eng.Lookup("10.1.1.1")

	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound, but got %v", err)
	}

	lant, long, err := eng.Lookup("20.1.1.1")

	if err != nil {
		t.Fatalf("Expected nil, but got %v", err)
	}

	if lant != 1 {
		t.Fatalf("Expected lantitude 1, but got %v", lant)
	}

	if long != 1 {
		t.Fatalf("Expected longtitude 1, but got %v", long)
	}

	_, _, err = eng.Lookup("30.1.1.1")

	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound, but got %v", err)
	}
}
