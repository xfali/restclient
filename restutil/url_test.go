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
