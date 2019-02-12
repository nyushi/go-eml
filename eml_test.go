package eml

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func strptr(s string) *string {
	return &s
}
func TestParse(t *testing.T) {
	testParseISO2022JP(t)
}

func testParseISO2022JP(t *testing.T) {
	for _, c := range []struct {
		Path    string
		Message *Message
	}{
		{
			Path: "testdata/body-iso-2022-jp-encoded",
			Message: &Message{
				RawHeaders: map[string][]string{
					"Return-Path":               []string{"<nyushi@example.com>"},
					"Message-Id":                []string{"<bfdd474e-8853-5f17-fa16-c29668e317c4@gmail.com>"},
					"User-Agent":                []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:60.0) Gecko/20100101 Thunderbird/60.5.0"},
					"Mime-Version":              []string{"1.0"},
					"Content-Type":              []string{"text/plain; charset=ISO-2022-JP; format=flowed; delsp=yes"},
					"Content-Transfer-Encoding": []string{"7bit"},
					"To":                        []string{"nyushi@example.com"},
					"From":                      []string{"Yushi Nakai <nyushi@example.com>"},
					"Subject":                   []string{"=?ISO-2022-JP?B?GyRCJDMkTiVhITwlayRPGyhCSVNPLTIwMjItSlAbJEIkRyQ5GyhC?="},
					"Date":                      []string{"Fri, 8 Feb 2019 23:37:49 +0900"},
					"Content-Language":          []string{"en-US"},
				},
				DecodedHeaders: map[string][]string{
					"From":                      {"Yushi Nakai <nyushi@example.com>"},
					"Message-Id":                {"<bfdd474e-8853-5f17-fa16-c29668e317c4@gmail.com>"},
					"Mime-Version":              {"1.0"},
					"Return-Path":               {"<nyushi@example.com>"},
					"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:60.0) Gecko/20100101 Thunderbird/60.5.0"},
					"Content-Language":          {"en-US"},
					"Content-Transfer-Encoding": {"7bit"},
					"To":                        {"nyushi@example.com"},
					"Subject":                   {"このメールはISO-2022-JPです"},
					"Date":                      {"Fri, 8 Feb 2019 23:37:49 +0900"},
					"Content-Type":              {"text/plain; charset=ISO-2022-JP; format=flowed; delsp=yes"},
				},
				RawBody:     "\r\n\x1b$B!!8cGZ$o$,$O$$$OG-$G$\"$k!#L>A0$O$^$@L5$$!#\x1b(B\r\n\x1b$B!!$I$3$G@8$l$?$+$H$s$H8+Ev$1$s$H$&$,$D$+$L!#2?$G$bGv0E$$$8$a$8$a$7$?=j$G\x1b(B \r\n\x1b$B%K%c!<%K%c!<5c$$$F$$$?;v$@$1$O5-21$7$F$$$k!#8cGZ$O$3$3$G;O$a$F?M4V$H$$$&\x1b(B \r\n\x1b$B$b$N$r8+$?!#$7$+$b$\"$H$GJ9$/$H$=$l$O=q@8$H$$$&?M4VCf$G0lHV`X0-$I$&$\"$/$J\x1b(B \r\n\x1b$B<oB2$G$\"$C$?$=$&$@!#$3$N=q@8$H$$$&$N$O;~!92f!9$rJa$D$+$^$($F<Q$K$F?)$&$H\x1b(B \r\n\x1b$B$$$&OC$G$\"$k!#$7$+$7$=$NEv;~$O2?$H$$$&9M$b$J$+$C$?$+$iJLCJ62$7$$$H$b;W$o\x1b(B \r\n\x1b$B$J$+$C$?!#$?$@H`$N>8$F$N$R$i$K:\\$;$i$l$F%9!<$H;}$A>e$2$i$l$?;~2?$@$+%U%o\x1b(B \r\n\x1b$B%U%o$7$?46$8$,$\"$C$?$P$+$j$G$\"$k!#>8$N>e$G>/$7Mn$A$D$$$F=q@8$N4i$r8+$?$N\x1b(B \r\n\x1b$B$,$$$o$f$k?M4V$H$$$&$b$N$N8+;O$_$O$8$a$G$\"$m$&!#$3$N;~L/$J$b$N$@$H;W$C$?\x1b(B \r\n\x1b$B46$8$,:#$G$b;D$C$F$$$k!#Bh0lLS$r$b$C$FAu>~$5$l$Y$-$O$:$N4i$,$D$k$D$k$7$F\x1b(B \r\n\x1b$B$^$k$GLt4L$d$+$s$@!#$=$N8e$4G-$K$b$@$$$V0)$\"$C$?$,$3$s$JJRNX$+$?$o$K$O0l\x1b(B \r\n\x1b$BEY$b=P2q$G$/$o$7$?;v$,$J$$!#$N$_$J$i$:4i$N??Cf$,$\"$^$j$KFM5/$7$F$$$k!#$=\x1b(B \r\n\x1b$B$&$7$F$=$N7j$NCf$+$i;~!9$W$&$W$&$H1l$1$`$j$r?a$/!#$I$&$b0v$`$;$]$/$F<B$K\x1b(B \r\n\x1b$B<e$C$?!#$3$l$,?M4V$N0{$`1lAp$?$P$3$H$$$&$b$N$G$\"$k;v$O$h$&$d$/$3$N:\"CN$C$?!#\x1b(B\r\n",
				DecodedBody: strptr("\r\n\u3000吾輩わがはいは猫である。名前はまだ無い。\r\n\u3000どこで生れたかとんと見当けんとうがつかぬ。何でも薄暗いじめじめした所で \r\nニャーニャー泣いていた事だけは記憶している。吾輩はここで始めて人間という \r\nものを見た。しかもあとで聞くとそれは書生という人間中で一番獰悪どうあくな \r\n種族であったそうだ。この書生というのは時々我々を捕つかまえて煮にて食うと \r\nいう話である。しかしその当時は何という考もなかったから別段恐しいとも思わ \r\nなかった。ただ彼の掌てのひらに載せられてスーと持ち上げられた時何だかフワ \r\nフワした感じがあったばかりである。掌の上で少し落ちついて書生の顔を見たの \r\nがいわゆる人間というものの見始みはじめであろう。この時妙なものだと思った \r\n感じが今でも残っている。第一毛をもって装飾されべきはずの顔がつるつるして \r\nまるで薬缶やかんだ。その後ご猫にもだいぶ逢あったがこんな片輪かたわには一 \r\n度も出会でくわした事がない。のみならず顔の真中があまりに突起している。そ \r\nうしてその穴の中から時々ぷうぷうと煙けむりを吹く。どうも咽むせぽくて実に \r\n弱った。これが人間の飲む煙草たばこというものである事はようやくこの頃知った。\r\n"),
			},
		},
		{
			Path: "testdata/body-utf8-encoded",
			Message: &Message{
				RawHeaders: map[string][]string{
					"Return-Path":               {"<nyushi@example.com>"},
					"Message-Id":                {"<74f9ebda-89bb-61e1-7954-f973c4e48d2e@gmail.com>"},
					"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:60.0) Gecko/20100101 Thunderbird/60.5.0"},
					"Mime-Version":              {"1.0"},
					"Content-Type":              {"text/plain; charset=utf-8; format=flowed"},
					"Content-Transfer-Encoding": {"8bit"},
					"To":                        {"nyushi@example.com"},
					"From":                      {"Yushi Nakai <nyushi@example.com>"},
					"Subject":                   {"=?UTF-8?B?44OG44K544OI44Oh44O844Or?="},
					"Date":                      {"Fri, 8 Feb 2019 23:30:34 +0900"},
					"Content-Language":          {"en-US"},
				},
				DecodedHeaders: map[string][]string{
					"From":                      {"Yushi Nakai <nyushi@example.com>"},
					"Message-Id":                {"<74f9ebda-89bb-61e1-7954-f973c4e48d2e@gmail.com>"},
					"Mime-Version":              {"1.0"},
					"Return-Path":               {"<nyushi@example.com>"},
					"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:60.0) Gecko/20100101 Thunderbird/60.5.0"},
					"Content-Language":          {"en-US"},
					"Content-Transfer-Encoding": {"8bit"},
					"To":                        {"nyushi@example.com"},
					"Subject":                   {"テストメール"},
					"Date":                      {"Fri, 8 Feb 2019 23:30:34 +0900"},
					"Content-Type":              {"text/plain; charset=utf-8; format=flowed"},
				},
				RawBody:     "テスト用のデータです。\r\n",
				DecodedBody: strptr("テスト用のデータです。\r\n"),
			},
		},
		{
			Path: "testdata/multipart",
			Message: &Message{
				RawHeaders: map[string][]string{
					"Content-Language": {"en-US"},
					"Content-Type":     {"multipart/mixed; boundary=\"------------661FE4891A157355A5B6C1E2\""},
					"Date":             {"Mon, 11 Feb 2019 11:51:57 +0900"},
					"From":             {"Yushi Nakai \u003cnyushi@example.com\u003e"},
					"Message-Id":       {"\u003c744ea9a4-95b8-007d-657b-a6d909ce61f0@gmail.com\u003e"},
					"Mime-Version":     {"1.0"},
					"Return-Path":      {"\u003cnyushi@example.com\u003e"},
					"Subject":          {"multipart"},
					"To":               {"nyushi@example.com"},
					"User-Agent":       {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:60.0) Gecko/20100101 Thunderbird/60.5.0"},
				},
				DecodedHeaders: map[string][]string{
					"Content-Language": {"en-US"},
					"Content-Type":     {"multipart/mixed; boundary=\"------------661FE4891A157355A5B6C1E2\""},
					"Date":             {"Mon, 11 Feb 2019 11:51:57 +0900"},
					"From":             {"Yushi Nakai \u003cnyushi@example.com\u003e"},
					"Message-Id":       {"\u003c744ea9a4-95b8-007d-657b-a6d909ce61f0@gmail.com\u003e"},
					"Mime-Version":     {"1.0"},
					"Return-Path":      {"\u003cnyushi@example.com\u003e"},
					"Subject":          {"multipart"},
					"To":               {"nyushi@example.com"},
					"User-Agent":       {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:60.0) Gecko/20100101 Thunderbird/60.5.0"},
				},
				RawBody: "This is a multi-part message in MIME format.\r\n--------------661FE4891A157355A5B6C1E2\r\nContent-Type: text/plain; charset=utf-8; format=flowed\r\nContent-Transfer-Encoding: 7bit\r\n\r\nthis is test mail\r\n\r\n--------------661FE4891A157355A5B6C1E2\r\nContent-Type: text/plain; charset=UTF-8; x-mac-type=\"0\"; x-mac-creator=\"0\";\r\n name=\"data\"\r\nContent-Transfer-Encoding: base64\r\nContent-Disposition: attachment;\r\n filename=\"data\"\r\n\r\ndGhpcyBpcyBkYXRhCg==\r\n--------------661FE4891A157355A5B6C1E2--\r\n",
				Parts: []*Message{
					{
						RawHeaders: map[string][]string{
							"Content-Transfer-Encoding": {"7bit"},
							"Content-Type":              {"text/plain; charset=utf-8; format=flowed"},
						},
						DecodedHeaders: map[string][]string{
							"Content-Transfer-Encoding": {"7bit"},
							"Content-Type":              {"text/plain; charset=utf-8; format=flowed"},
						},
						RawBody:     "this is test mail\r\n",
						DecodedBody: strptr("this is test mail\r\n"),
					},
					{
						RawHeaders: map[string][]string{
							"Content-Disposition":       {"attachment; filename=\"data\""},
							"Content-Transfer-Encoding": {"base64"},
							"Content-Type":              {"text/plain; charset=UTF-8; x-mac-type=\"0\"; x-mac-creator=\"0\"; name=\"data\""},
						},
						DecodedHeaders: map[string][]string{
							"Content-Disposition":       {"attachment; filename=\"data\""},
							"Content-Transfer-Encoding": {"base64"},
							"Content-Type":              {"text/plain; charset=UTF-8; x-mac-type=\"0\"; x-mac-creator=\"0\"; name=\"data\""},
						},
						RawBody:     "dGhpcyBpcyBkYXRhCg==",
						DecodedBody: strptr("this is data\n"),
					},
				},
			},
		},
	} {
		r, err := os.Open(c.Path)
		if err != nil {
			t.Fatalf("failed to read %s: %s", c.Path, err)
		}
		msg, err := Parse(r)
		if err != nil {
			t.Fatalf("error at parse: %s", err)
		}
		if err := messageEqual(msg, c.Message); err != nil {
			t.Errorf("msg not equals: %s", err)
		}
	}
}

func messageEqual(got, expected *Message) error {
	if !reflect.DeepEqual(expected.RawHeaders, got.RawHeaders) {
		g, _ := json.MarshalIndent(got.RawHeaders, "", "  ")
		e, _ := json.MarshalIndent(expected.RawHeaders, "", "  ")
		return fmt.Errorf("RawHeaders not match\ngot: %s\nexpected: %s", g, e)
	}
	if !reflect.DeepEqual(expected.DecodedHeaders, got.DecodedHeaders) {
		g, _ := json.MarshalIndent(got.DecodedHeaders, "", "  ")
		e, _ := json.MarshalIndent(expected.DecodedHeaders, "", "  ")
		return fmt.Errorf("DecodedHeaders not match\ngot: %s\nexpected: %s", g, e)
	}

	if expected.RawBody != got.RawBody {
		return fmt.Errorf("RawBody not match:\n  got=`%s`\n  expect=`%s`", got.RawBody, expected.RawBody)
	}

	if expected.DecodedBody != nil && got.DecodedBody != nil && *expected.DecodedBody != *got.DecodedBody {
		return fmt.Errorf("DecodedBody not match:\n  got=%s\n  expect=%s", *got.DecodedBody, *expected.DecodedBody)
	}

	if len(expected.Parts) != len(got.Parts) {
		return fmt.Errorf("Parts size is not match: got=%d, expected=%d", len(got.Parts), len(expected.Parts))
	}

	for i, e := range expected.Parts {
		if err := messageEqual(got.Parts[i], e); err != nil {
			return fmt.Errorf("parts not match at %d: %s", i, err)
		}
	}
	return nil
}
