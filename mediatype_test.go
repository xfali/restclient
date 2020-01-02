// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package restclient

import "testing"

func TestParseMediaType(t *testing.T) {
    t.Run("*/*", func(t *testing.T) {
        s := ParseMediaType("*/*")
        o := ParseMediaType(MediaTypeJson)
        if !s.Includes(o) {
            t.Fatal(`*/* not match: `, MediaTypeJson)
        }
    })

    t.Run("application/*", func(t *testing.T) {
        s := ParseMediaType("application/*")
        o := ParseMediaType(MediaTypeJson)
        if !s.Includes(o) {
            t.Fatal(`application/* not match: `, MediaTypeJson)
        }
    })

    t.Run("application/*json", func(t *testing.T) {
        s := ParseMediaType("application/*json")
        o := ParseMediaType(MediaTypeJson)
        if !s.Includes(o) {
            t.Fatal(`application/*json not match: `, MediaTypeJson)
        }
    })

    t.Run("application/json", func(t *testing.T) {
        s := ParseMediaType("application/json")
        o := ParseMediaType(MediaTypeJson)
        if !s.Includes(o) {
            t.Fatal(`application/json not match: `, MediaTypeJson)
        }
    })

    t.Run("application/xml", func(t *testing.T) {
        s := ParseMediaType("application/xml")
        o := ParseMediaType(MediaTypeJson)
        if s.Includes(o) {
            t.Fatal(`application/xml not match: `, MediaTypeJson)
        }
    })
}

func TestParseMediaType2(t *testing.T) {
    t.Run("*/*", func(t *testing.T) {
        s := ParseMediaType("*/*")
        o := ParseMediaType(MediaTypeJsonUtf8)
        if !s.Includes(o) {
            t.Fatal(`*/* not match: `, MediaTypeJsonUtf8)
        }
    })

    t.Run("application/*", func(t *testing.T) {
        s := ParseMediaType("application/*")
        o := ParseMediaType(MediaTypeJson)
        if !s.Includes(o) {
            t.Fatal(`application/* not match: `, MediaTypeJsonUtf8)
        }
    })

    t.Run("application/*json", func(t *testing.T) {
        s := ParseMediaType("application/*json")
        o := ParseMediaType(MediaTypeJsonUtf8)
        if !s.Includes(o) {
            t.Fatal(`application/*json not match: `, MediaTypeJsonUtf8)
        }
    })

    t.Run("application/json", func(t *testing.T) {
        s := ParseMediaType("application/json")
        o := ParseMediaType(MediaTypeJsonUtf8)
        if s.Includes(o) {
            t.Fatal(`application/json not match: `, MediaTypeJsonUtf8)
        }
    })

    t.Run("application/xml", func(t *testing.T) {
        s := ParseMediaType("application/xml")
        o := ParseMediaType(MediaTypeJsonUtf8)
        if s.Includes(o) {
            t.Fatal(`application/xml not match: `, MediaTypeJsonUtf8)
        }
    })
}
