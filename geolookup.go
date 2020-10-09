package cvsgeolookup

import (
	"encoding/binary"
	"encoding/csv"
	"io"
	"net"
	"sort"
	"strconv"
)

type Metrics interface {
}

type nometrics struct{}

type record struct {
	begin      uint32  // first ip of segment
	end        uint32  // last ip of segment
	lantitude  float32 // lantitude
	longtitude float32 // longtitude
}

type engine struct {
	records    []record
	beginindex []uint64
	endindex   []uint64

	options options
}

// New returns new lookup engine
func New(opts ...Option) (*engine, error) {

	options := options{
		fieldNameBegin:      "start",
		fieldNameEnd:        "end",
		fieldNameLantitude:  "lantitude",
		fieldNameLongtitude: "longtitude",
		commaRune:           ',',
		commentRune:         '#',
		metrics:             &nometrics{},
	}

	for _, o := range opts {
		o.apply(&options)
	}

	return &engine{
		options: options,
	}, nil
}

func (e *engine) Load(r io.Reader) error {
	if r == nil {
		return ErrReadInterfaceRequired
	}

	data := csv.NewReader(r)
	data.Comma = e.options.commaRune
	data.Comment = e.options.commentRune
	data.ReuseRecord = true

	// Read header and find right fields' positions
	indexBegin := -1
	indexEnd := -1
	indexLantitude := -1
	indexLongtitude := -1
	indexSkip := -1

	if header, err := data.Read(); err == nil {
		for index, field := range header {
			switch field {
			case e.options.fieldNameBegin:
				indexBegin = index
			case e.options.fieldNameEnd:
				indexEnd = index
			case e.options.fieldNameLantitude:
				indexLantitude = index
			case e.options.fieldNameLongtitude:
				indexLongtitude = index
			case e.options.fieldNameSkip:
				indexSkip = index
			}
		}
	} else {
		return err
	}

	if indexBegin == -1 {
		return ErrNoBeginField
	}

	if indexEnd == -1 {
		return ErrNoEndField
	}

	if indexLantitude == -1 {
		return ErrNoLantitudeField
	}

	if indexLongtitude == -1 {
		return ErrNoLongtitudeField
	}

	// Build records
	records := make([]record, 0)

	for {
		line, err := data.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if indexSkip >= 0 {
			if line[indexSkip] == e.options.skipValue {
				continue
			}
		}

		var record record

		if record.begin, err = e.parseIP(line[indexBegin]); err != nil {
			return err
		}

		if record.end, err = e.parseIP(line[indexEnd]); err != nil {
			return err
		}

		if record.begin > record.end {
			return ErrIncorrectSegment
		}

		if record.lantitude, err = e.parseFloat(line[indexLantitude]); err != nil {
			return err
		}

		if record.longtitude, err = e.parseFloat(line[indexLongtitude]); err != nil {
			return err
		}

		records = append(records, record)
	}

	// Just make sure all segments are in monotone order
	sort.Slice(records[:], func(i, j int) bool {
		return records[i].begin < records[j].begin
	})

	// Build indexes
	beginindex := make([]uint64, 256)
	endindex := make([]uint64, 256)

	count := uint64(len(records))
	for i := uint64(0); i < count; i++ {
		end := (records[i].begin >> 24)
		begin := (records[count-i-1].begin >> 24)
		beginindex[begin] = count - i - 1
		endindex[end] = i
	}

	e.records = records
	e.beginindex = beginindex
	e.endindex = endindex

	return nil
}

func (e *engine) Lookup(ip string) (lantitude float32, longtitude float32, err error) {
	err = nil
	lantitude = 0
	longtitude = 0

	if e.records == nil {
		err = ErrNotInitialized
		return
	}

	look, err := e.parseIP(ip)

	if err != nil {
		return
	}

	begin := e.beginindex[look>>24]
	end := e.endindex[look>>24]

	for {
		current := (begin + end) / 2

		if e.records[current].begin <= look && e.records[current].end >= look {
			//found!
			lantitude = e.records[current].lantitude
			longtitude = e.records[current].longtitude
			return
		}

		if begin == end {
			break
		}

		if e.records[current].begin > look {
			// move to left side
			end = current
			continue
		}

		if e.records[current].end < look {
			// move to right side
			if begin != current {
				begin = current
			} else {
				begin = current + 1
			}
			continue
		}

		if begin > end {
			break
		}
	}

	err = ErrNotFound

	return
}

func (e *engine) parseIP(val string) (uint32, error) {
	ip := net.ParseIP(val)
	if ip == nil {
		return 0, ErrWrongIPFormat
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

func (e *engine) parseFloat(val string) (float32, error) {
	float, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}

	// Some precision corrections
	if float > 0 {
		float += 0.0000000001
	} else if float < 0 {
		float -= 0.0000000001
	}

	return float32(float), nil
}
