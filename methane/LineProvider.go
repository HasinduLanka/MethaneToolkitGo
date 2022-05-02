package methane

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LineProvider struct {
	ch chan string
}
type ReusableLineProvider struct {
	GetInstance func() *LineProvider
	Discard     func()
}
type LinePipe struct {
	Proc func(*LineProvider) *LineProvider
}

func NewReusableLineProvider_Static(S string) *ReusableLineProvider {
	return &ReusableLineProvider{
		GetInstance: func() *LineProvider {
			LP := &LineProvider{make(chan string, 1)}
			go func() {
				LP.ch <- S
				close(LP.ch)
			}()
			return LP
		},
	}
}

// On demand line generator from start to finish (inclusive). Zeros are padded if output is shorted than minLength.
func NewReusableLineProvider_IntRange(start int, finish int, minLength int) *ReusableLineProvider {
	minLengthInteger := strconv.Itoa(minLength)

	return &ReusableLineProvider{
		GetInstance: func() *LineProvider {
			LP := &LineProvider{make(chan string, 512)}
			go func() {
				for i := start; i <= finish; i++ {
					LP.ch <- padNumberWithZero(i, minLengthInteger)
				}
				close(LP.ch)
			}()
			return LP
		},
	}
}

func NewReusableLineProvider_FromString(S string) *ReusableLineProvider {
	return &ReusableLineProvider{
		GetInstance: func() *LineProvider {
			return NewLineProvider_FromString(S)
		},
	}
}

func NewLineProvider_FromString(S string) *LineProvider {
	LP := &LineProvider{make(chan string)}

	go func() {
		for _, line := range strings.Split(S, "\n") {
			LP.ch <- line
		}

		close(LP.ch)
	}()

	return LP
}

func (lp *LineProvider) CacheToMemmory() *ReusableLineProvider {
	var Cache []string = []string{}

	for line := range lp.ch {
		Cache = append(Cache, line)
	}

	return &ReusableLineProvider{
		GetInstance: func() *LineProvider {
			o := &LineProvider{make(chan string)}

			go func() {
				for _, line := range Cache {
					o.ch <- line
				}
			}()
			return o
		},
	}
}

func NewLinePipe_TrimAndFilterOutEmpty() LinePipe {
	return LinePipe{
		func(lp *LineProvider) *LineProvider {
			olp := LineProvider{make(chan string)}

			go func() {
				for line := range lp.ch {
					line = strings.TrimSpace(line)
					if len(line) != 0 {
						olp.ch <- line
					}
				}
				close(olp.ch)
			}()

			return &olp
		},
	}
}

func NewLinePipe_Prefix(prefix string) LinePipe {
	return LinePipe{
		func(lp *LineProvider) *LineProvider {
			olp := LineProvider{make(chan string)}

			go func() {
				for line := range lp.ch {
					olp.ch <- (prefix + line)
				}
				close(olp.ch)
			}()

			return &olp
		},
	}
}

func NewLinePipe_FilterOutShellStyleComments() LinePipe {
	rgx := regexp.MustCompile("^ *#")
	return NewLinePipe_FilterMatches(rgx, false)
}

// FilterIn ? Keep matches : Skip matches
func NewLinePipe_FilterMatches(rgx *regexp.Regexp, FilterIn bool) LinePipe {
	return LinePipe{
		func(lp *LineProvider) *LineProvider {
			olp := LineProvider{make(chan string)}

			if FilterIn {
				go func() {
					for line := range lp.ch {
						if rgx.MatchString(line) {
							olp.ch <- line
						}
					}
					close(olp.ch)
				}()
			} else {
				go func() {
					for line := range lp.ch {
						if !rgx.MatchString(line) {
							olp.ch <- line
						}
					}
					close(olp.ch)
				}()
			}

			return &olp
		},
	}
}

func (rlp *ReusableLineProvider) PrintAll() {
	for line := range rlp.GetInstance().ch {
		Print(line)
	}
}

var _uniqey int = 0

func UniqueStamp() string {
	B, _ := time.Now().MarshalText()
	S := url.PathEscape(string(B))

	_uniqey++
	return S + "-" + strconv.Itoa(_uniqey)
}

func NewCacheFilename() string {
	return WSCache + ".cached-" + UniqueStamp()
}

func padNumberWithZero(value int, minLengthInteger string) string {
	return fmt.Sprintf("%0"+minLengthInteger+"d", value)
}
