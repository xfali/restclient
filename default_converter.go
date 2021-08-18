// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"reflect"
	"strings"
)

type BaseConverter struct {
	SupportType []MediaType
}

func (c *BaseConverter) SupportMediaType() []MediaType {
	return c.SupportType
}

func (c *BaseConverter) CanHandler(mediaType MediaType) bool {
	for i := range c.SupportType {
		if c.SupportType[i].Includes(mediaType) {
			return true
		}
	}

	return false
}

type ByteConverter struct {
	BaseConverter
}

type ByteEncoder struct {
	w io.Writer
}
type ByteDecoder struct {
	r io.Reader
}

func NewByteConverter() *ByteConverter {
	return &ByteConverter{
		BaseConverter{[]MediaType{
			ParseMediaType(MediaTypeAll),
			ParseMediaType(MediaTypeOctetStream),
		}},
	}
}

func (c *ByteConverter) CreateEncoder(w io.Writer) Encoder {
	return &ByteEncoder{w: w}
}
func (c *ByteConverter) CreateDecoder(r io.Reader) Decoder {
	return &ByteDecoder{r: r}
}

func (c *ByteEncoder) Encode(i interface{}) (int64, error) {
	if s, ok := i.([]byte); ok {
		n, err := c.w.Write(s)
		return int64(n), err
	}
	return 0, errors.New("ByteConverter not support Serialize ")
}

func (c *ByteConverter) CanEncode(o interface{}, mediaType MediaType) bool {
	if !mediaType.IsWildcard() && !c.CanHandler(mediaType) {
		return false
	}
	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Slice {
		return false
	}
	if t.Elem().Kind() == reflect.Uint8 {
		return true
	}

	return false
}

func (c *ByteDecoder) Decode(result interface{}) (int64, error) {
	buf := bytes.NewBuffer(nil)
	n, err := io.Copy(buf, c.r)
	if err != nil {
		return n, err
	}

	//在CanDeserialize中已经明确了result的类型
	v := reflect.ValueOf(result)
	v = v.Elem()
	v.SetBytes(buf.Bytes())
	return n, io.EOF
}

func (c *ByteConverter) CanDecode(o interface{}, mediaType MediaType) bool {
	if !mediaType.IsWildcard() && !c.CanHandler(mediaType) {
		return false
	}

	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Ptr {
		return false
	}
	t = t.Elem()
	if t.Kind() != reflect.Slice {
		return false
	}

	if t.Elem().Kind() == reflect.Uint8 {
		return true
	}
	return false
}

type StringConverter struct {
	BaseConverter
}

type StringEncoder struct {
	w io.Writer
}
type StringDecoder struct {
	r io.Reader
}

func (c *StringConverter) CreateEncoder(w io.Writer) Encoder {
	return &StringEncoder{w: w}
}
func (c *StringConverter) CreateDecoder(r io.Reader) Decoder {
	return &StringDecoder{r: r}
}

func NewStringConverter() *StringConverter {
	return &StringConverter{
		BaseConverter{[]MediaType{
			ParseMediaType(MediaTypeTextPlain),
			ParseMediaType(MediaTypeAll),
		}},
	}
}

func (c *StringEncoder) Encode(i interface{}) (int64, error) {
	if s, ok := i.(string); ok {
		n, err := io.WriteString(c.w, s)
		return int64(n), err
	}
	return 0, errors.New("StringConverter not support Serialize ")
}

func (c *StringConverter) CanEncode(o interface{}, mediaType MediaType) bool {
	if !mediaType.IsWildcard() && !c.CanHandler(mediaType) {
		return false
	}
	t := reflect.TypeOf(o)
	return t.Kind() == reflect.String
}

func (c *StringDecoder) Decode(result interface{}) (int64, error) {
	buf := &strings.Builder{}

	n, err := io.Copy(buf, c.r)
	if err != nil {
		return n, err
	}

	//在CanDeserialize中已经明确了result的类型
	v := reflect.ValueOf(result)
	v = v.Elem()
	v.SetString(buf.String())
	return n, io.EOF
}

func (c *StringConverter) CanDecode(o interface{}, mediaType MediaType) bool {
	if !mediaType.IsWildcard() && !c.CanHandler(mediaType) {
		return false
	}

	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Ptr {
		return false
	}
	t = t.Elem()
	return t.Kind() == reflect.String
}

type JsonConverter struct {
	BaseConverter
}

type JsonEncoder struct {
	e *json.Encoder
}

type JsonDecoder struct {
	last int64
	d    *json.Decoder
}

func (c *JsonConverter) CreateEncoder(w io.Writer) Encoder {
	return &JsonEncoder{
		e: json.NewEncoder(w),
	}
}
func (c *JsonConverter) CreateDecoder(r io.Reader) Decoder {
	return &JsonDecoder{d: json.NewDecoder(r)}
}

func NewJsonConverter() *JsonConverter {
	return &JsonConverter{
		BaseConverter{[]MediaType{
			ParseMediaType(MediaTypeJson),
			BuildMediaType("application", "*json"),
		}},
	}
}

func (c *JsonEncoder) Encode(i interface{}) (int64, error) {
	err := c.e.Encode(i)
	return 0, err
}

func (c *JsonConverter) CanEncode(o interface{}, mediaType MediaType) bool {
	if !mediaType.IsWildcard() && !c.CanHandler(mediaType) {
		return false
	}
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Interface, reflect.Struct, reflect.Map:
		return true
	case reflect.Slice:
		return t.Elem().Kind() != reflect.Uint8
	default:
		return false
	}
	return true
}

func (c *JsonDecoder) Decode(result interface{}) (int64, error) {
	err := c.d.Decode(result)
	n := c.last
	c.last = c.d.InputOffset()
	return c.last - n, err
}

func (c *JsonConverter) CanDecode(o interface{}, mediaType MediaType) bool {
	if !mediaType.IsWildcard() && !c.CanHandler(mediaType) {
		return false
	}

	t := reflect.TypeOf(o)
	//must be ptr
	if t.Kind() != reflect.Ptr {
		return false
	} else {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Interface, reflect.Struct, reflect.Map:
		return true
	case reflect.Slice:
		return t.Elem().Kind() != reflect.Uint8
	default:
		return false
	}
}

type XmlConverter struct {
	BaseConverter
}

type XmlEncoder struct {
	e *xml.Encoder
}

type XmlDecoder struct {
	last int64
	d    *xml.Decoder
}

func (c *XmlConverter) CreateEncoder(w io.Writer) Encoder {
	return &XmlEncoder{
		e: xml.NewEncoder(w),
	}
}
func (c *XmlConverter) CreateDecoder(r io.Reader) Decoder {
	return &XmlDecoder{d: xml.NewDecoder(r)}
}

func NewXmlConverter() *XmlConverter {
	return &XmlConverter{
		BaseConverter{[]MediaType{
			ParseMediaType(MediaTypeXml),
			BuildMediaType("application", "*xml"),
		}},
	}
}

func (c *XmlEncoder) Encode(i interface{}) (int64, error) {
	err := c.e.Encode(i)
	return 0, err
}

func (c *XmlConverter) CanEncode(o interface{}, mediaType MediaType) bool {
	if !mediaType.IsWildcard() && !c.CanHandler(mediaType) {
		return false
	}
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Interface, reflect.Struct:
		return true
	default:
		return false
	}

	return true
}

func (c *XmlDecoder) Decode(result interface{}) (int64, error) {
	err := c.d.Decode(result)
	n := c.last
	c.last = c.d.InputOffset()
	return c.last - n, err
}

func (c *XmlConverter) CanDecode(o interface{}, mediaType MediaType) bool {
	if !mediaType.IsWildcard() && !c.CanHandler(mediaType) {
		return false
	}
	t := reflect.TypeOf(o)
	//must be ptr
	if t.Kind() != reflect.Ptr {
		return false
	} else {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Interface, reflect.Struct:
		return true
	default:
		return false
	}
}

type YamlConverter struct {
	BaseConverter
}

type YamlEncoder struct {
	e *yaml.Encoder
}

type YamlDecoder struct {
	last int64
	d    *yaml.Decoder
}

func (c *YamlConverter) CreateEncoder(w io.Writer) Encoder {
	return &YamlEncoder{
		e: yaml.NewEncoder(w),
	}
}
func (c *YamlConverter) CreateDecoder(r io.Reader) Decoder {
	return &YamlDecoder{d: yaml.NewDecoder(r)}
}

func NewYamlConverter() *YamlConverter {
	return &YamlConverter{
		BaseConverter{[]MediaType{
			ParseMediaType(MediaTypeYaml),
			BuildMediaType("application", "*yaml"),
		}},
	}
}

func (c *YamlEncoder) Encode(i interface{}) (int64, error) {
	err := c.e.Encode(i)
	return 0, err
}

func (c *YamlConverter) CanEncode(o interface{}, mediaType MediaType) bool {
	if !mediaType.IsWildcard() && !c.CanHandler(mediaType) {
		return false
	}
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Interface, reflect.Struct, reflect.Map:
		return true
	case reflect.Slice:
		return t.Elem().Kind() != reflect.Uint8
	default:
		return false
	}
	return true
}

func (c *YamlDecoder) Decode(result interface{}) (int64, error) {
	err := c.d.Decode(result)
	//n := c.last
	//c.last = c.d.InputOffset()
	return 0, err
}

func (c *YamlConverter) CanDecode(o interface{}, mediaType MediaType) bool {
	if !mediaType.IsWildcard() && !c.CanHandler(mediaType) {
		return false
	}

	t := reflect.TypeOf(o)
	//must be ptr
	if t.Kind() != reflect.Ptr {
		return false
	} else {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Interface, reflect.Struct, reflect.Map:
		return true
	case reflect.Slice:
		return t.Elem().Kind() != reflect.Uint8
	default:
		return false
	}
}

func chooseEncoder(converters []Converter, o interface{}, mediaType MediaType) (Converter, error) {
	l := len(converters)
	for l > 0 {
		l--
		if converters[l].CanEncode(o, mediaType) {
			return converters[l], nil
		}
	}
	return nil, errors.New("Cannot Serialize Object ")
}

func chooseDecoder(converters []Converter, ret interface{}, mediaType MediaType) (Converter, error) {
	l := len(converters)
	for l > 0 {
		l--
		if converters[l].CanDecode(ret, mediaType) {
			return converters[l], nil
		}
	}
	return nil, errors.New("Cannot Deserialize Object ")
}

func getDefaultMediaType(converter Converter) MediaType {
	return converter.SupportMediaType()[0]
}
