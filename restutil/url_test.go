/*
 * Copyright 2022 Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package restutil

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

func TestEncodeQuery(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		url, err := Query("a", "1", "b", 2)
		if err != nil {
			t.Fatal(err)
		}
		if url != "a=1&b=2" {
			t.Fatal(url)
		}

		t.Log(url)
	})

	t.Run("2", func(t *testing.T) {
		url := EncodeQuery(map[string]interface{}{
			"a": "1",
			"b": 2,
		})
		if url != "a=1&b=2" {
			t.Fatal(url)
		}

		t.Log(url)
	})
}

func TestReplaceUrl(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		url := ReplaceUrl("x", "", "", map[string]interface{}{
			"a": "1",
			"b": 2,
		})
		if url != "x" {
			t.Fatal(url)
		}

		t.Log(url)
	})

	t.Run("2", func(t *testing.T) {
		url := ReplaceUrl("x/:a/tt?b=:b", ":", "", map[string]interface{}{
			"a": "1",
			"b": 2,
		})
		if url != "x/1/tt?b=2" {
			t.Fatal(url)
		}

		t.Log(url)
	})

	t.Run("3", func(t *testing.T) {
		url := ReplaceUrl("x/:a/tt/:a?b=:b", ":", "", map[string]interface{}{
			"a": "1",
			"b": 2,
		})
		if url != "x/1/tt/1?b=2" {
			t.Fatal(url)
		}

		t.Log(url)
	})
}

func TestUrlBuilder(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		b := NewUrlBuilder("x")
		b.PathVariable("a", "1")
		b.PathVariable("b", 2)
		url := b.Build()
		if url != "x" {
			t.Fatal(url)
		}

		t.Log(url)
	})

	t.Run("2", func(t *testing.T) {
		b := NewUrlBuilder("x/:a/tt?b=:b")
		b.PathVariable("a", "1")
		b.PathVariable("b", 2)
		url := b.Build()
		if url != "x/1/tt?b=2" {
			t.Fatal(url)
		}

		t.Log(url)
	})

	t.Run("3", func(t *testing.T) {
		b := NewUrlBuilder("x/:a/tt/:a?b=:b")
		b.PathVariable("a", "1")
		b.PathVariable("b", 2)
		url := b.Build()
		if url != "x/1/tt/1?b=2" {
			t.Fatal(url)
		}

		t.Log(url)
	})

	t.Run("4", func(t *testing.T) {
		b := NewUrlBuilder("x/:a/tt/:b")
		b.PathVariable("a", "1")
		b.PathVariable("b", 2)
		b.QueryVariable("c", 100)
		b.QueryVariable("d", 1.1)
		url := b.Build()
		if url != "x/1/tt/2?c=100&d=1.1" {
			t.Fatal(url)
		}

		t.Log(url)
	})

	t.Run("4", func(t *testing.T) {
		b := NewUrlBuilder("x/:a/tt/:b?")
		b.PathVariable("a", "1")
		b.PathVariable("b", 2)
		b.QueryVariable("c", 100)
		b.QueryVariable("d", 1.1)
		url := b.Build()
		if url != "x/1/tt/2?c=100&d=1.1" {
			t.Fatal(url)
		}

		t.Log(url)
	})
}
