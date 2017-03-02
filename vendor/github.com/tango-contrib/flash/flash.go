// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flash

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/session"
)

var (
	FlashName      = "tango_flash"
	FlashSeperator = "TANGOFLASH"

	_ Flasher = &Flash{}
)

type Data map[string]interface{}

// FlashData is a tools to maintain data when using across request.
type Flash struct {
	readed  Data
	flushed Data
	session *session.Session
	*Options
	saved bool
}

func (f *Flash) setFlash(sess *session.Session, readed Data, opt *Options) {
	f.readed = readed
	f.flushed = make(Data)
	f.session = sess
	f.Options = opt
}

func (f *Flash) FlushData() Data {
	return f.flushed
}

func (f *Flash) Merge() {
	for k, v := range f.readed {
		f.flushed[k] = v
	}
}

func (f *Flash) Data() Data {
	return f.readed
}

func (f *Flash) Get(key string) interface{} {
	return f.readed[key]
}

func (f *Flash) Set(key string, value interface{}) {
	f.readed[key] = value
	f.flushed[key] = value
}

func (f *Flash) Add(kvs Data) {
	for k, v := range kvs {
		f.Set(k, v)
	}
}

func (f *Flash) Save() {
	if f.saved {
		return
	}

	for key, _ := range f.readed {
		f.session.Del(f.Options.FlashName + f.Options.FlashSeperator + key)
	}

	var keys = make([]string, 0)
	for k, v := range f.flushed {
		f.session.Set(f.Options.FlashName+f.Options.FlashSeperator+k, v)
		keys = append(keys, k)
	}
	f.session.Set(f.Options.FlashName, keys)
	f.saved = true
}

type Flasher interface {
	setFlash(*session.Session, Data, *Options)
	FlushData() Data
	Save()
}

type Options struct {
	FlashName      string
	FlashSeperator string
}

func prepareOptions(opts []Options) Options {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}

	if len(opt.FlashName) == 0 {
		opt.FlashName = FlashName
	}
	if len(opt.FlashSeperator) == 0 {
		opt.FlashSeperator = FlashSeperator
	}
	return opt
}

// Flash return a FlashData handler.
func Flashes(sessions *session.Sessions, opts ...Options) tango.HandlerFunc {
	opt := prepareOptions(opts)
	return func(ctx *tango.Context) {
		var flasher Flasher
		var ok bool
		if action := ctx.Action(); action != nil {
			if flasher, ok = action.(Flasher); ok {
				sess := sessions.Session(ctx.Req(), ctx.ResponseWriter)
				fd := make(Data)
				if keys, has := sess.Get(opt.FlashName).([]string); has {
					for _, key := range keys {
						fd[key] = sess.Get(opt.FlashName + opt.FlashSeperator + key)
					}
				}
				flasher.setFlash(sess, fd, &opt)
			}
		}

		ctx.Next()

		if ok {
			flasher.Save()
		}
	}
}
