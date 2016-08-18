package sstrg

import (
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"sync"
	"time"
)

//SessionsTable is a thread-safe wrapper under session history
type SessionsTable struct {
	sync.RWMutex
	//TODO: replace it with slice-map pair (iterating over map is very slow and we do it very often)
	H map[string]*SessionHistory
}

//SessionHistory session statistics
type SessionHistory struct {
	sync.RWMutex
	Started  time.Time
	Ended    time.Time
	Requests []*RequestData
	Active   bool
}

//RequestData Information about session
type RequestData struct {
	Path        string
	Referer     string
	Method      string
	Header      http.Header
	Cookies     []*http.Cookie
	ContentType int
	Time        time.Time
}

//Content types for resources
const (
	ImageContentType       = iota
	FontContentType        = iota
	JSContentType          = iota
	CSSContentType         = iota
	OtherStaticContentType = iota
	//Unknown dynamic tyoe
	DynamicContentType = iota
)

//Regular expressions to match for content type
var regExps = []*regexp.Regexp{
	regexp.MustCompile(`\.(gif|jpg|jpeg|tiff|png|svg|ico)$`),
	regexp.MustCompile(`\.(ttf|otf|eot|woff|woff2)$`),
	regexp.MustCompile(`\.js$`),
	regexp.MustCompile(`\.css$`),
	regexp.MustCompile(`\.(xml|json|zip|gz|pdf|ico|doc|docx|xls|ppt|txt)$`),
}

//Replacements for resource paths
var resourcePathMasks = []string{
	"/image.jpg",
	"/font.ttf",
	"/script.js",
	"/style.css",
	"/other.static",
}

//GetSessionKey generates session key form request headers
func GetSessionKey(r *http.Request) string {
	return r.Header.Get("X-Forwarded-For") // + "|" + r.Header.Get("User-Agent")
}

//GetRequestData RequestData factory method
func GetRequestData(r *http.Request, useDataHeader bool) *RequestData {
	var request = new(RequestData)

	request.Path = r.URL.Path
	request.Method = r.Method
	request.Cookies = r.Cookies()
	request.Header = r.Header
	//Set request referer
	var ref, err = url.Parse(r.Referer())

	if err == nil {
		request.Referer = ref.Path
	} else {
		request.Referer = r.Referer()
	}

	//Set time
	if useDataHeader {
		t, err := time.Parse(time.RFC1123, r.Header.Get("Date"))

		if err == nil {
			request.Time = t.UTC()
		} else {
			request.Time = time.Now().UTC()
		}
	} else {
		request.Time = time.Now().UTC()
	}

	request.Path, request.ContentType = GetContentType(request.Path)

	return request
}

//GetContentType get content type by resource path
func GetContentType(path string) (string, int) {
	//Set content type
	for requestType, regExp := range regExps {
		if regExp.MatchString(path) {
			path = resourcePathMasks[requestType]
			return path, requestType
		}
	}

	return path, DynamicContentType
}

//By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *RequestData) bool

// planetSorter joins a By function and a slice of Planets to be sorted.
type requestSorter struct {
	requests []*RequestData
	by       func(p1, p2 *RequestData) bool // Closure used in the Less method.
}

//Len is part of sort.Interface.
func (s *requestSorter) Len() int {
	return len(s.requests)
}

//Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *requestSorter) Less(i, j int) bool {
	return s.by(s.requests[i], s.requests[j])
}

//Swap is part of sort.Interface.
func (s *requestSorter) Swap(i, j int) {
	s.requests[i], s.requests[j] = s.requests[j], s.requests[i]
}

//Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(requests []*RequestData) {
	rs := &requestSorter{
		requests: requests,
		by:       by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(rs)
}

//SortRequestsByTime sorts requets by time
func SortRequestsByTime(requests []*RequestData) {
	//Sort requests by timestamp
	timeOrder := func(r1, r2 *RequestData) bool {
		return r2.Time.Before(r1.Time)
	}

	By(timeOrder).Sort(requests)
}
