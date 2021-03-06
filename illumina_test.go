// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package illumina

import (
	"testing"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/feat"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

type tester struct {
	name string
	desc string
}

func (t tester) Name() string        { return t.name }
func (t tester) Description() string { return t.desc }

type alphaTester struct {
	name  string
	desc  string
	alpha alphabet.Alphabet
}

func (t alphaTester) Name() string                { return t.name }
func (t alphaTester) Description() string         { return t.desc }
func (t alphaTester) Alphabet() alphabet.Alphabet { return t.alpha }

var underscored = alphabet.MustComplement(alphabet.NewComplementor(
	"acgt_",
	feat.DNA,
	alphabet.MustPair(alphabet.NewPairing("acgtnxACGTNX-_", "tgcanxTGCANX-_")),
	'-', 'n',
	!alphabet.CaseSensitive,
))

func (s *S) TestParse(c *check.C) {
	for i, t := range []struct {
		in   Interface
		meta Metadata
	}{
		{
			tester{"HWUSI-EAS100R:6:73:941:1973#0/1", ""},
			Metadata{
				Type:       PreCasava,
				Instrument: "HWUSI-EAS100R",
				Run:        -1, // Pre-Casava
				Lane:       6,
				Tile:       73,
				Coordinate: Coordinate{941, 1973},
				Mate:       1,
				Multiplex:  Multiplex{Index: 0},
			},
		},
		{
			tester{"HWUSI-EAS100R:6:73:941:1973#ATCACG/1", ""},
			Metadata{
				Type:       PreCasava,
				Instrument: "HWUSI-EAS100R",
				Run:        -1, // Pre-Casava
				Lane:       6,
				Tile:       73,
				Coordinate: Coordinate{941, 1973},
				Mate:       1,
				Multiplex:  Multiplex{Index: -1, Tag: "ATCACG"},
			},
		},
		{
			tester{"EAS139:136:FC706VJ:2:2104:15343:197393", "1:Y:18:ATCACG"},
			Metadata{
				Type:        Casava,
				Instrument:  "EAS139",
				Run:         136,
				FlowCell:    "FC706VJ",
				Lane:        2,
				Tile:        2104,
				Coordinate:  Coordinate{15343, 197393},
				Mate:        1,
				BadRead:     true,
				ControlBits: 18,
				Multiplex:   Multiplex{Index: -1, Tag: "ATCACG"},
			},
		},
		{ // This test is for sequences that have passed through a pipeline that strips desc fields.
			tester{"EAS139:136:FC706VJ:2:2104:15343:197393", ""},
			Metadata{
				Type:        Casava,
				Instrument:  "EAS139",
				Run:         136,
				FlowCell:    "FC706VJ",
				Lane:        2,
				Tile:        2104,
				Coordinate:  Coordinate{15343, 197393},
				Mate:        0,
				BadRead:     false,
				ControlBits: -1,
				Multiplex:   Multiplex{Index: -1},
			},
		},
		{ // This test is for tags where unusual tag characters have been used.
			alphaTester{"FCC4T93ACXX:2:1101:1899:1998#ATTACTCG_GGCTCTGA/2", "", underscored},
			Metadata{
				Type:       PreCasava,
				Instrument: "FCC4T93ACXX",
				Run:        -1,
				Lane:       2,
				Tile:       1101,
				Coordinate: Coordinate{1899, 1998},
				Mate:       2,
				BadRead:    false,
				Multiplex:  Multiplex{Index: -1, Tag: "ATTACTCG_GGCTCTGA"},
			},
		},
	} {
		m, err := Parse(t.in)
		c.Check(err, check.Equals, nil, check.Commentf("Test %d", i))
		c.Check(m, check.Equals, t.meta, check.Commentf("Test %d", i))
	}
}
