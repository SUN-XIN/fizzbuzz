package main

import (
	"fmt"
	"sync"
	"time"

	locip "github.com/SUN-XIN/iplocation"
)

/////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////    Server   //////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////
type server struct {
	Password    string
	StartDate   int64
	HistoryCall []*fizzBuzz

	// response in gzip ?
	ForceGzip    bool
	ForceGzipNum int64

	Mutex *sync.Mutex
}

func NewServer(forceGzip bool, forceGzipNum int64) *server {
	return &server{
		Password:    ADMIN_PASSWORD,
		StartDate:   time.Now().Unix(),
		HistoryCall: make([]*fizzBuzz, 0, MAX_HIST),

		ForceGzip:    forceGzip,
		ForceGzipNum: forceGzipNum,

		Mutex: &sync.Mutex{},
	}
}

// Add a new client call info into server.HistoryCall
func (s *server) AddNewCall(fb *fizzBuzz) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.HistoryCall = append(s.HistoryCall, fb)
}

// Check if need to force gzip
func (s *server) CheckForceGzip(limit int) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if !s.ForceGzip {
		return false
	}

	return limit >= int(s.ForceGzipNum)
}

// Only admin can access
func (s *server) CheckAdmin(pass string) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	return s.Password == pass
}

// Return server's stat info
func (s *server) ShowInfo() string {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	res := fmt.Sprintf("Start at: %s, %d call in total", time.Unix(s.StartDate, 0), len(s.HistoryCall))
	for _, call := range s.HistoryCall {
		res = fmt.Sprintf("%s\n\t%s (%s %s) called at %s with params %+v, stat %s", res,
			call.locip.IP, call.locip.CountryCode, call.locip.CountryName,
			time.Unix(call.date, 0),
			call.request,
			call.status,
		)
	}

	return res
}

// Update server's configuration
func (s *server) UpdateConf(c *configuration) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.ForceGzip = c.ForceGzip
	s.ForceGzipNum = c.ForceGzipNum
	s.Password = c.Password
}

// all information about client's call
type fizzBuzz struct {
	request *clientRequest
	locip   *locip.IPStackResponse
	date    int64
	status  string // ok, failed
}

func NewFizzBuzz(cr *clientRequest, clientIP string) *fizzBuzz {
	return &fizzBuzz{
		request: cr,
		locip: &locip.IPStackResponse{
			IP: clientIP,
		},
		date:   time.Now().Unix(),
		status: "failed",
	}
}

func (fb *fizzBuzz) SetOK() {
	fb.status = "ok"
}

/////////////////////////////////////////////////////////////////////////////////////
//////////////////////    Client request/Response   /////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

type clientRequest struct {
	// required
	String1 string `json:"string1"`
	String2 string `json:"string2"`
	Int1    int    `json:"int1"`
	Int2    int    `json:"int2"`
	Limit   int    `json:"limit"`

	// optional
	ResponseGzip bool `json:"response_gzip"`
}

func (cr *clientRequest) validate() error {
	if cr.String1 == "" || cr.String2 == "" {
		return fmt.Errorf("string1/string2 must not be empty")
	}

	if cr.Int1 <= 0 || cr.Int2 <= 0 {
		return fmt.Errorf("int1/int2 must be positive")
	}

	if cr.Limit <= 0 {
		return fmt.Errorf("limit must be positive")
	}

	return nil
}

type clientResponse struct {
	Result string `json:"result"`
}

/////////////////////////////////////////////////////////////////////////////////////
///////////////////////////    CONFIGURATION   //////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////
type configuration struct {
	// TODO: must be more securely, ex: public_key, private_key
	Password string `json:"password"`

	ForceGzip    bool  `json:"force_gzip"`
	ForceGzipNum int64 `json:"force_gzip_num"`
}

func (c *configuration) validate() error {
	if c.Password == "" {
		// no password for ADMIN is possible ?
	}

	if c.ForceGzip && c.ForceGzipNum <= 0 {
		c.ForceGzipNum = DEFAULT_FORCE_GZIP_NUM
	}

	return nil
}
