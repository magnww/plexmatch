package main

import (
	"flag"
	"testing"
)

func TestMatchSession(t *testing.T) {
	flag.Parse()
	ParseRegex()

	se := match("[orion origin] Isekai Nonbiri Nouka [02] [1080p] [H265 AAC] [CHS]", seRegexList)
	se2 := match("[BeanSub][Vinland_Saga_S2][05][GB][1080P][x264_AAC]", seRegexList)
	se3 := match("[VCB-Studio] Kyokou Suiri [Ma10p_1080p]", seRegexList)
	se4 := match("The.Last.of.Us.S01E04.2160p.HMAX.WEB-DL.DDP5.1.Atmos.DV.MKV.x265-SMURF[rartv]", seRegexList)
	se5 := match("[ANi] 虛構推理 第二季 - 05 [1080P][Baha][WEB-DL][AAC AVC][CHT]", seRegexList)
	if se != 0 {
		t.Errorf("session test 1 failed")
	}
	if se2 != 2 {
		t.Errorf("session test 2 failed")
	}
	if se3 != 0 {
		t.Errorf("session test 3 failed")
	}
	if se4 != 1 {
		t.Errorf("session test 4 failed")
	}
	if se5 != 2 {
		t.Errorf("session test 5 failed")
	}
}

func TestMatchEpisode(t *testing.T) {
	flag.Parse()
	ParseRegex()

	ep := match("[orion origin] Isekai Nonbiri Nouka [06] [1080p] [H265 AAC] [CHS].mp4", epRegexList)
	ep2 := match("[BeanSub][Vinland_Saga_S2][05][GB][1080P][x264_AAC].mp4", epRegexList)
	ep3 := match("[VCB-Studio] Kyokou Suiri [02][Ma10p_1080p][x265_flac].mkv", epRegexList)
	ep4 := match("The.Last.of.Us.S01E04.Please.Hold.on.to.My.Hand.2160p.HMAX.WEB-DL.DDP5.1.Atmos.DV.H.265-SMURF.mkv", epRegexList)
	if ep != 6 {
		t.Errorf("episode test 1 failed")
	}
	if ep2 != 5 {
		t.Errorf("episode test 2 failed")
	}
	if ep3 != 2 {
		t.Errorf("episode test 3 failed")
	}
	if ep4 != 4 {
		t.Errorf("episode test 4 failed")
	}
}

func TestMatchTitle(t *testing.T) {
	flag.Parse()
	ParseRegex()

	ti := removeNumber("[ANi] 虛構推理 第二季 - 05 [1080P][Baha][WEB-DL][AAC AVC][CHT]", epRegexList)
	if ti != "[ANi] 虛構推理 第二季 -  [1080P][Baha][WEB-DL][AAC AVC][CHT]" {
		t.Errorf("title test 1 failed")
	}
	ti2 := removeNumber("[Nekomoe kissaten][NieR Automata Ver1.1a][03][1080p][CHS]", epRegexList)
	if ti2 != "[Nekomoe kissaten][NieR Automata Ver1.1a][1080p][CHS]" {
		t.Errorf("title test 2 failed")
	}
}
