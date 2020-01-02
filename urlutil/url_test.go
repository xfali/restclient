/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package urlutil

import "testing"

func TestQueryUrl(t *testing.T) {
    t.Run("1", func(t *testing.T) {
        url := QueryUrl("x", map[string]string{
            "a": "1",
            "b": "2",
        })
        if url != "x?a=1&b=2" {
            t.Fatal(url)
        }

        t.Log(url)
    })

    t.Run("2", func(t *testing.T) {
        url := QueryUrl("x?", map[string]string{
            "a": "1",
            "b": "2",
        })
        if url != "x?a=1&b=2" {
            t.Fatal(url)
        }

        t.Log(url)
    })

    t.Run("3", func(t *testing.T) {
        url := QueryUrl("x?c=3", map[string]string{
            "a": "1",
            "b": "2",
        })
        if url != "x?c=3&a=1&b=2" {
            t.Fatal(url)
        }

        t.Log(url)
    })

    t.Run("4", func(t *testing.T) {
        url := QueryUrl("x?c=3&", map[string]string{
            "a": "1",
            "b": "2",
        })
        if url != "x?c=3&a=1&b=2" {
            t.Fatal(url)
        }

        t.Log(url)
    })
}

func TestPlaceholderUrl(t *testing.T) {
    t.Run("1", func(t *testing.T) {
        url := PlaceholderUrl("x", map[string]string{
            "a": "1",
            "b": "2",
        })
        if url != "x" {
            t.Fatal(url)
        }

        t.Log(url)
    })

    t.Run("2", func(t *testing.T) {
        url := PlaceholderUrl("x/${a}?b=${b}", map[string]string{
            "a": "1",
            "b": "2",
        })
        if url != "x/1?b=2" {
            t.Fatal(url)
        }

        t.Log(url)
    })
}
