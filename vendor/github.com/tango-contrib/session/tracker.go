// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Tracker provide and set sessionid
type Tracker interface {
	SetMaxAge(maxAge time.Duration)
	Get(req *http.Request) (Id, error)
	Set(req *http.Request, rw http.ResponseWriter, id Id)
	Clear(rw http.ResponseWriter)
}

// CookieTracker provide sessionid from cookie
type CookieTracker struct {
	Name     string
	MaxAge   time.Duration
	Lock     sync.Mutex
	Secure   bool
	RootPath string
	Domain   string
}

func NewCookieTracker(name string, maxAge time.Duration, secure bool, rootPath string) *CookieTracker {
	return &CookieTracker{
		Name:     name,
		MaxAge:   maxAge,
		Secure:   secure,
		RootPath: rootPath,
	}
}

func (tracker *CookieTracker) SetMaxAge(maxAge time.Duration) {
	tracker.MaxAge = maxAge
}

func (tracker *CookieTracker) Get(req *http.Request) (Id, error) {
	cookie, err := req.Cookie(tracker.Name)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil
		}
		return "", err
	}
	if cookie.Value == "" {
		return Id(""), nil
	}
	id, _ := url.QueryUnescape(cookie.Value)
	return Id(id), nil
}

func (tracker *CookieTracker) Set(req *http.Request, rw http.ResponseWriter, id Id) {
	sid := url.QueryEscape(string(id))
	tracker.Lock.Lock()
	defer tracker.Lock.Unlock()
	cookie, _ := req.Cookie(tracker.Name)
	if cookie == nil {
		cookie = &http.Cookie{
			Name:     tracker.Name,
			Value:    sid,
			Path:     tracker.RootPath,
			Domain:   tracker.Domain,
			HttpOnly: true,
			Secure:   tracker.Secure,
		}

		req.AddCookie(cookie)
	} else {
		cookie.Value = sid
	}
	http.SetCookie(rw, cookie)
}

func (tracker *CookieTracker) Clear(rw http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     tracker.Name,
		Path:     tracker.RootPath,
		Domain:   tracker.Domain,
		HttpOnly: true,
		Secure:   tracker.Secure,
		Expires:  time.Date(0, 1, 1, 0, 0, 0, 0, time.Local),
		MaxAge:   -1,
	}
	http.SetCookie(rw, &cookie)
}

var _ Tracker = NewCookieTracker("test", 0, false, "/")

// UrlTracker provide sessionid from url
type UrlTracker struct {
	Key         string
	ReplaceLink bool
}

func NewUrlTracker(key string, replaceLink bool) *UrlTracker {
	return &UrlTracker{key, replaceLink}
}

func (tracker *UrlTracker) Get(req *http.Request) (Id, error) {
	sessionId := req.URL.Query().Get(tracker.Key)
	if sessionId != "" {
		sessionId, _ = url.QueryUnescape(sessionId)
		return Id(sessionId), nil
	}

	return Id(""), nil
}

func (tracker *UrlTracker) Set(req *http.Request, rw http.ResponseWriter, id Id) {
	if tracker.ReplaceLink {

	}
}

func (tracker *UrlTracker) SetMaxAge(maxAge time.Duration) {

}

func (tracker *UrlTracker) Clear(rw http.ResponseWriter) {
}

var (
	_ Tracker = NewUrlTracker("id", false)
)

//for SWFUpload ...
func NewCookieUrlTracker(name string, maxAge time.Duration, secure bool, rootPath string) *CookieUrlTracker {
	return &CookieUrlTracker{
		CookieTracker: CookieTracker{
			Name:     name,
			MaxAge:   maxAge,
			Secure:   secure,
			RootPath: rootPath,
		},
	}
}

type CookieUrlTracker struct {
	CookieTracker
}

func (tracker *CookieUrlTracker) Get(req *http.Request) (Id, error) {
	sessionId := req.URL.Query().Get(tracker.Name)
	if sessionId != "" {
		sessionId, _ = url.QueryUnescape(sessionId)
		return Id(sessionId), nil
	}

	return tracker.CookieTracker.Get(req)
}

type HeaderTracker struct {
	Name string
}

func NewHeaderTracker(name string) *HeaderTracker {
	return &HeaderTracker{
		Name: name,
	}
}

func (tracker *HeaderTracker) SetMaxAge(maxAge time.Duration) {
}

func (tracker *HeaderTracker) Get(req *http.Request) (Id, error) {
	val := req.Header.Get(tracker.Name)
	return Id(val), nil
}

func (tracker *HeaderTracker) Set(req *http.Request, rw http.ResponseWriter, id Id) {
	rw.Header().Set(tracker.Name, string(id))
}

func (tracker *HeaderTracker) Clear(rw http.ResponseWriter) {
}

/*
type CompositeTracker struct {
	Trackers []Tracker
}

func NewCompositeTracker(trackers ...Tracker) *CompositeTracker {
	return &CompositeTracker{trackers}
}

func (trackers *CompositeTracker) Get(req *http.Request) (Id, error) {
	for _, tracker := range trackers.Trackers {
		if id, err := tracker.Get(req); err == nil {
			return id, nil
		}
	}
	return Id(""), nil
}

func (trackers *CompositeTracker) SetMaxAge(maxAge time.Duration) {
	for _, tracker := range trackers.Trackers {
		tracker.SetMaxAge(maxAge)
	}
}

func (trackers *CompositeTracker) Set(req *http.Request, rw http.ResponseWriter, id Id) {
	for _, tracker := range trackers.Trackers {
		tracker.Set(req, rw, id)
	}
}

func (trackers *CompositeTracker) Clear(rw http.ResponseWriter) {
	for _, tracker := range trackers.Trackers {
		tracker.Clear(rw)
	}
}*/
