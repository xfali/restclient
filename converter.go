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

package restclient

import (
	"io"
)

type Encoder interface {
	Encode(o interface{}) (int64, error)
}

type Decoder interface {
	Decode(o interface{}) (int64, error)
}

type Converter interface {
	CreateEncoder(io.Writer) Encoder
	CreateDecoder(io.Reader) Decoder

	CanEncode(o interface{}, mediaType MediaType) bool
	CanDecode(o interface{}, mediaType MediaType) bool
	SupportMediaType() []MediaType
}
