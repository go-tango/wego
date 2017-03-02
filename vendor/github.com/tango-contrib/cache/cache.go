// Copyright 2013 Beego Authors
// Copyright 2014 The Macaron Authors
// Copyright 2015 The Tango Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package cache is a middleware that provides the cache management of Tango.
package cache

import (
	"fmt"

	"github.com/lunny/tango"
)

// Adapter is the interface that operates the cache data.
type Adapter interface {
	// Put puts value into cache with key and expire time.
	Put(key string, val interface{}, timeout int64) error
	// Get gets cached value by given key.
	Get(key string) interface{}
	// Delete deletes cached value by given key.
	Delete(key string) error
	// Incr increases cached int-type value by given key as a counter.
	Incr(key string) error
	// Decr decreases cached int-type value by given key as a counter.
	Decr(key string) error
	// IsExist returns true if cached value exists.
	IsExist(key string) bool
	// Flush deletes all cached data.
	Flush() error
	// StartAndGC starts GC routine based on config string settings.
	StartAndGC(opt Options) error
}

// Options represents a struct for specifying configuration options for the cache middleware.
type Options struct {
	// Name of adapter. Default is "memory".
	Adapter string
	// Adapter configuration, it's corresponding to adapter.
	AdapterConfig string
	// GC interval time in seconds. Default is 60.
	Interval int
	// Occupy entire database. Default is false.
	OccupyMode bool
	// Configuration section name. Default is "cache".
	Section string
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}
	if len(opt.Section) == 0 {
		opt.Section = "cache"
	}
	if len(opt.Adapter) == 0 {
		opt.Adapter = "memory"
	}
	if opt.Interval == 0 {
		opt.Interval = 60
	}
	if len(opt.AdapterConfig) == 0 {
		opt.AdapterConfig = "data/caches"
	}
	return opt
}

// NewAdapter creates and returns a new cacheAdapter by given adapter name and configuration.
// It panics when given adapter isn't registered and starts GC automatically.
func NewAdapter(name string, opt Options) (Adapter, error) {
	adapter, ok := adapters[name]
	if !ok {
		return nil, fmt.Errorf("cache: unknown adapter '%s'(forgot to import?)", name)
	}
	return adapter, adapter.StartAndGC(opt)
}

var adapters = make(map[string]Adapter)

// Register registers a adapter.
func Register(name string, adapter Adapter) {
	if adapter == nil {
		panic("cache: cannot register adapter with nil value")
	}
	if _, dup := adapters[name]; dup {
		panic(fmt.Errorf("cache: cannot register adapter '%s' twice", name))
	}
	adapters[name] = adapter
}

// Cacher provides tango action handler to get proper handler
type Cacher interface {
	SetCaches(*Caches)
}

// Caches
type Caches struct {
	options Options
	adapter Adapter
}

// Cache maintains cache adapter and tango handler interface
type Cache struct {
	*Caches
}

var _ Cacher = new(Cache)

// set Caches
func (c *Cache) SetCaches(cs *Caches) {
	c.Caches = cs
}

// Put puts value into cache with key and expire time.
func (c *Caches) Put(key string, val interface{}, timeout int64) error {
	return c.adapter.Put(key, val, timeout)
}

// Get gets cached value by given key.
func (c *Caches) Get(key string) interface{} {
	return c.adapter.Get(key)
}

// Delete deletes cached value by given key.
func (c *Caches) Delete(key string) error {
	return c.adapter.Delete(key)
}

// Incr increases cached int-type value by given key as a counter.
func (c *Caches) Incr(key string) error {
	return c.adapter.Incr(key)
}

// Decr decreases cached int-type value by given key as a counter.
func (c *Caches) Decr(key string) error {
	return c.adapter.Decr(key)
}

// IsExist returns true if cached value exists.
func (c *Caches) IsExist(key string) bool {
	return c.adapter.IsExist(key)
}

// Flush deletes all cached data.
func (c *Caches) Flush() error {
	return c.adapter.Flush()
}

// Options return cache option.
func (c *Caches) Option() Options {
	return c.options
}

// Handle implement tango.Handle
func (c *Caches) Handle(ctx *tango.Context) {
	if action := ctx.Action(); ctx != nil {
		if s, ok := action.(Cacher); ok {
			s.SetCaches(c)
		}
	}

	ctx.Next()
}

// New is a middleware that maps a cache.Cache service into the tango handler chain.
// An single variadic cache.Options struct can be optionally provided to configure.
func New(options ...Options) *Caches {
	opt := prepareOptions(options)
	adapter, err := NewAdapter(opt.Adapter, opt)
	if err != nil {
		panic(err)
	}
	return &Caches{
		options: opt,
		adapter: adapter,
	}
}
