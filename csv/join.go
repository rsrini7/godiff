package csv

import (
	"github.com/rsrini7/godiff/utils"
)

// A Join can be used to construct a process that will join two streams of CSV records by matching
// records from each stream on the specified key columns.
type Join struct {
	LeftKeys   []string // the names of the keys from the left stream
	RightKeys  []string // the names of the keys from the right stream
	Numeric    []string // the names of the keys in the left stream that are numeric keys
	LeftOuter  bool     // perform a left outer join - left rows are copied even if there is no matching right row
	RightOuter bool     // perform a right outer join - right rows are copied even if there is no matching left row
}

// A decorator for a reader that returns groups of consecutive records from the underlying reader
// that have the same key.
type groupReader struct {
	next   Record
	group  []Record
	reader Reader
	key    []string
	tokey  func(r Record) []string
	less   func(l, r []string) bool
}

// Fill up the group slice with the set of records in the underlying stream that have the same key
func (g *groupReader) fill() bool {
	if g.next == nil {
		g.next = <-g.reader.C()
	}
	if g.next == nil {
		return false
	} else {
		g.key = g.tokey(g.next)
	}
	g.group = []Record{g.next}
	for {
		g.next = <-g.reader.C()
		var k []string
		if g.next != nil {
			k = g.tokey(g.next)
		}
		if g.next == nil || g.less(k, g.key) || g.less(g.key, k) {
			return true
		} else {
			g.group = append(g.group, g.next)
		}
	}
}

// Answers true if get() will yield a new group
func (g *groupReader) hasNext() bool {
	if g.group != nil && len(g.group) > 0 {
		return true
	}
	return g.fill()
}

// Answers the next group of records from the stream.
func (g *groupReader) get() []Record {
	if !g.hasNext() {
		panic("illegal state: get() called when hasNext() is false")
	}
	r := g.group
	g.group = nil
	return r
}

// Construct a key comparison function for key values
func (p *Join) less() StringSliceComparator {
	return (&SortKeys{
		Keys:    p.LeftKeys,
		Numeric: p.Numeric,
	}).AsStringSliceComparator()
}

// split the headers into the set of all headers, the set of key headers, the set of left headers
// and the set of right headers
func (p *Join) headers(leftHeader []string, rightHeader []string) ([]string, []string, []string, []string) {
	i, a, _ := utils.Intersect(leftHeader, p.LeftKeys)
	_, b, _ := utils.Intersect(rightHeader, p.RightKeys)
	f := make([]string, len(i)+len(a)+len(b))

	copy(f, i)
	copy(f[len(i):], a)
	copy(f[len(i)+len(a):], b)

	return f, i, a, b
}

func (p *Join) run(left Reader, right Reader, builder WriterBuilder, errCh chan<- error) {
	errCh <- func() (err error) {
		defer left.Close()
		defer right.Close()

		less := p.less()

		leftBlank := NewRecordBuilder(left.Header())([]string{})
		rightBlank := NewRecordBuilder(right.Header())([]string{})

		outputHeader, keyHeader, leftHeader, rightHeader := p.headers(left.Header(), right.Header())
		writer := builder(outputHeader)
		defer writer.Close(err)

		leftG := &groupReader{reader: left, less: less, tokey: (&SortKeys{Keys: p.LeftKeys}).AsStringProjection()}
		rightG := &groupReader{reader: right, less: less, tokey: (&SortKeys{Keys: p.RightKeys}).AsStringProjection()}

		w := func(k []string, l, r Record) error {
			o := writer.Blank()
			for i, h := range keyHeader {
				o.Put(h, k[i])
			}
			for _, h := range leftHeader {
				o.Put(h, l.Get(h))
			}
			for _, h := range rightHeader {
				o.Put(h, r.Get(h))
			}
			return writer.Write(o)
		}

		for leftG.hasNext() && rightG.hasNext() {
			if less(leftG.key, rightG.key) {
				//
				// copy left to output
				//
				for _, r := range leftG.get() {
					if p.LeftOuter {
						if err := w(leftG.key, r, rightBlank); err != nil {
							return err
						}
					}
				}
			} else if less(rightG.key, leftG.key) {
				//
				// copy right to output
				//
				for _, r := range rightG.get() {
					if p.RightOuter {
						if err := w(rightG.key, leftBlank, r); err != nil {
							return err
						}
					}
				}
			} else {
				// copy join product to output
				rg := rightG.get()
				for _, l := range leftG.get() {
					for _, r := range rg {
						if err := w(leftG.key, l, r); err != nil {
							return err
						}
					}
				}
			}
		}
		for leftG.hasNext() {
			//
			// copy left to output
			//
			for _, r := range leftG.get() {
				if p.LeftOuter {
					if err := w(leftG.key, r, rightBlank); err != nil {
						return err
					}
				}
			}
		}
		for rightG.hasNext() {
			//
			// copy right to output
			//
			for _, r := range rightG.get() {
				if p.RightOuter {
					if err := w(rightG.key, leftBlank, r); err != nil {
						return err
					}
				}
			}
		}

		lerr := left.Error()
		if lerr != nil {
			return lerr
		} else {
			return right.Error()
		}
	}()

}

type joinProcess struct {
	join  *Join
	right Reader
}

// Binds the specified reader as the right-hand side of a join and returns
// a Process whose reader will be considered as the left-hand side of the join.
func (p *Join) WithRight(r Reader) Process {
	return &joinProcess{
		join:  p,
		right: r,
	}
}

func (j *joinProcess) Run(r Reader, builder WriterBuilder, errCh chan<- error) {
	j.join.run(r, j.right, builder, errCh)
}
