package log4go

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : pattern.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/6 10:46
* 修改历史 : 1. [2022/4/6 10:46] 创建文件 by NST
*/

const (
	FORMAT_DEFAULT = "[%D %T %m] [%L] (%S) %M"
	FORMAT_SHORT   = "[%t %d] [%L] %M"
	FORMAT_ABBREV  = "[%L] %M"
)

type TimeSlice struct {
	LastUpdateSeconds    int64
	shortTime, shortDate string
	longTime, longDate   string
	lock                 sync.Mutex
}

var timeSliceCache = &TimeSlice{}

func (t *TimeSlice) GetTimeSlice(time time.Time) {
	secs := time.UnixNano() / 1e9
	cached := getCachedTimeSlice()
	if cached.LastUpdateSeconds != secs {
		month, day, year := time.Month(), time.Day(), time.Year()
		hour, minute, second := time.Hour(), time.Minute(), time.Second()
		//
		//ns := time.Nanosecond()
		//mil := ns / 1e6
		zone, _ := time.Zone()

		t.LastUpdateSeconds = secs
		t.shortTime = fmt.Sprintf("%02d:%02d", hour, minute)
		t.shortDate = fmt.Sprintf("%02d/%02d/%02d", day, month, year%100)
		t.longTime = fmt.Sprintf("%02d:%02d:%02d %s", hour, minute, second, zone)
		t.longDate = fmt.Sprintf("%04d/%02d/%02d", year, month, day)

		setCachedTimeSlice(*t)
	} else {
		*t = cached
		//t.LastUpdateSeconds = cached.LastUpdateSeconds
		//t.shortTime = cached.shortTime
		//t.shortDate = cached.shortDate
		//t.longTime = cached.longTime
		//t.longDate = cached.longDate
	}
}

func setCachedTimeSlice(f TimeSlice) {
	timeSliceCache.lock.Lock()
	defer timeSliceCache.lock.Unlock()
	timeSliceCache.LastUpdateSeconds = f.LastUpdateSeconds
	timeSliceCache.shortTime = f.shortTime
	timeSliceCache.shortDate = f.shortDate
	timeSliceCache.longTime = f.longTime
	timeSliceCache.longDate = f.longDate
}

func getCachedTimeSlice() TimeSlice {
	timeSliceCache.lock.Lock()
	defer timeSliceCache.lock.Unlock()
	return *timeSliceCache
}

type PatternConverter interface {
	FormatLogRecord(pattern string, rec LogRecord) string
}

type defaultPatternConverter struct {
}

func getDefaultPatternConvert() PatternConverter {
	return *DefaultPatternConverter
}

var DefaultPatternConverter = &defaultPatternConverter{}

// FormatLogRecord
// Known format codes:
// %T - Time (15:04:05 MST)
// %t - Time (15:04)
// %D - Date (2006/01/02)
// %d - Date (01/02/06)
// %o - UnixMilli (1649224315937)
// %m - Millisecond (607)  //
// %n - Nanosecond (794411287)
// %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
// %l - Label (TagName)
// %S - Source
// %M - Message
// Ignores unknown formats
// Recommended: "[%D %T] [%L] (%S) %M"
func (dc defaultPatternConverter) FormatLogRecord(pattern string, rec LogRecord) string {
	if len(pattern) == 0 {
		return ""
	}

	out := bytes.NewBuffer(make([]byte, 0, 1024))

	timeSlice := &TimeSlice{}
	timeSlice.GetTimeSlice(rec.Created)

	// Split the string into pieces by % signs
	pieces := bytes.Split([]byte(pattern), []byte{'%'})
	// Iterate over the pieces, replacing known formats
	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'T':
				out.WriteString(timeSlice.longTime)
			case 'o':
				out.WriteString(strconv.FormatInt(rec.Created.UnixMilli(), 10))
			case 'n':
				out.WriteString(LeftPad(strconv.Itoa(rec.Created.Nanosecond()), 9, '0'))
			case 'm':
				out.WriteString(LeftPad(strconv.Itoa(rec.Created.Nanosecond()/1e6), 3, '0'))
			case 't':
				out.WriteString(timeSlice.shortTime)
			case 'D':
				out.WriteString(timeSlice.longDate)
			case 'd':
				out.WriteString(timeSlice.shortDate)
			case 'l':
				out.WriteString(rec.TagName)
			case 'L':
				out.WriteString(levelStrings[rec.Level])
			case 'S':
				out.WriteString(rec.Source)
			case 's':
				slice := strings.Split(rec.Source, "/")
				out.WriteString(slice[len(slice)-1])
			case 'M':
				out.WriteString(rec.Message)
			}
			if len(piece) > 1 {
				out.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			out.Write(piece)
		}
	}
	out.WriteByte('\n')

	return out.String()
}
