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

func NewByteConverter() *ByteConverter {
    return &ByteConverter{
        BaseConverter{[]MediaType{
            ParseMediaType(MediaTypeAll),
            ParseMediaType(MediaTypeOctetStream),
        }},
    }
}

func (c *ByteConverter) Serialize(i interface{}) (io.Reader, error) {
    if s, ok := i.([]byte); ok {
        return bytes.NewReader(s), nil
    }
    return nil, errors.New("ByteConverter not support Serialize ")
}

func (c *ByteConverter) CanSerialize(o interface{}, mediaType MediaType) bool {
    if !c.CanHandler(mediaType) {
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

func (c *ByteConverter) Deserialize(r io.Reader, result interface{}) (int, error) {
    buf := bytes.NewBuffer(nil)
    n, err := io.Copy(buf, r)
    if err != nil {
        return int(n), err
    }

    d := buf.Bytes()
    //在CanDeserialize中已经明确了result的类型
    v := reflect.ValueOf(result)
    v = v.Elem()
    v.SetBytes(d)
    return int(n), nil
}

func (c *ByteConverter) CanDeserialize(o interface{}, mediaType MediaType) bool {
    if !c.CanHandler(mediaType) {
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

func NewStringConverter() *StringConverter {
    return &StringConverter{
        BaseConverter{[]MediaType{
            ParseMediaType(MediaTypeTextPlain),
            ParseMediaType(MediaTypeAll),
        }},
    }
}

func (c *StringConverter) Serialize(i interface{}) (io.Reader, error) {
    if s, ok := i.(string); ok {
        return strings.NewReader(s), nil
    }
    return nil, errors.New("StringConverter not support Serialize ")
}

func (c *StringConverter) CanSerialize(o interface{}, mediaType MediaType) bool {
    if !c.CanHandler(mediaType) {
        return false
    }
    t := reflect.TypeOf(o)
    if t.Kind() == reflect.String {
        return true
    }
    return false
}

func (c *StringConverter) Deserialize(r io.Reader, result interface{}) (int, error) {
    buf := bytes.NewBuffer(nil)
    n, err := io.Copy(buf, r)
    if err != nil {
        return int(n), err
    }

    d := buf.Bytes()
    //在CanDeserialize中已经明确了result的类型
    v := reflect.ValueOf(result)
    v = v.Elem()
    v.SetString(string(d))
    return int(n), nil
}

func (c *StringConverter) CanDeserialize(o interface{}, mediaType MediaType) bool {
    if !c.CanHandler(mediaType) {
        return false
    }

    t := reflect.TypeOf(o)
    if t.Kind() != reflect.Ptr {
        return false
    }
    t = t.Elem()
    if t.Kind() == reflect.String {
        return true
    }
    return false
}

type JsonConverter struct {
    BaseConverter
}

func NewJsonConverter() *JsonConverter {
    return &JsonConverter{
        BaseConverter{[]MediaType{
            ParseMediaType(MediaTypeJson),
            BuildMediaType("application", "*json"),
        }},
    }
}

func (c *JsonConverter) Serialize(i interface{}) (io.Reader, error) {
    d, err := json.Marshal(i)
    if err != nil {
        return nil, err
    }
    return bytes.NewReader(d), nil
}

func (c *JsonConverter) CanSerialize(o interface{}, mediaType MediaType) bool {
    if !c.CanHandler(mediaType) {
        return false
    }
    t := reflect.TypeOf(o)
    if t.Kind() == reflect.Struct {
        return true
    }
    return false
}

func (c *JsonConverter) Deserialize(r io.Reader, result interface{}) (int, error) {
    buf := bytes.NewBuffer(nil)
    n, err := io.Copy(buf, r)
    if err != nil {
        return int(n), err
    }

    d := buf.Bytes()
    return int(n), json.Unmarshal(d, result)
}

func (c *JsonConverter) CanDeserialize(o interface{}, mediaType MediaType) bool {
    if !c.CanHandler(mediaType) {
        return false
    }

    t := reflect.TypeOf(o)
    if t.Kind() != reflect.Ptr {
        return false
    }
    t = t.Elem()
    if t.Kind() == reflect.Struct {
        return true
    }
    return false
}


type XmlConverter struct {
    BaseConverter
}

func NewXmlConverter() *XmlConverter {
    return &XmlConverter{
        BaseConverter{[]MediaType{
            ParseMediaType(MediaTypeXml),
            BuildMediaType("application", "*xml"),
        }},
    }
}

func (c *XmlConverter) Serialize(i interface{}) (io.Reader, error) {
    d, err := xml.Marshal(i)
    if err != nil {
        return nil, err
    }
    return bytes.NewReader(d), nil
}

func (c *XmlConverter) CanSerialize(o interface{}, mediaType MediaType) bool {
    if !c.CanHandler(mediaType) {
        return false
    }
    t := reflect.TypeOf(o)
    if t.Kind() == reflect.Struct {
        return true
    }
    return false
}

func (c *XmlConverter) Deserialize(r io.Reader, result interface{}) (int, error) {
    buf := bytes.NewBuffer(nil)
    n, err := io.Copy(buf, r)
    if err != nil {
        return int(n), err
    }

    d := buf.Bytes()
    return int(n), xml.Unmarshal(d, result)
}

func (c *XmlConverter) CanDeserialize(o interface{}, mediaType MediaType) bool {
    if !c.CanHandler(mediaType) {
        return false
    }

    t := reflect.TypeOf(o)
    if t.Kind() != reflect.Ptr {
        return false
    }
    t = t.Elem()
    if t.Kind() == reflect.Struct {
        return true
    }
    return false
}

func doSerialize(converters []Converter, o interface{}, mediaType MediaType) (io.Reader, Converter, error) {
    l := len(converters)
    for l>0 {
        l--
        if converters[l].CanSerialize(o, mediaType) {
            r, err := converters[l].Serialize(o)
            if err == nil {
                return r, converters[l], nil
            }
        }
    }
    return nil, nil, errors.New("Cannot Serialize Object ")
}

func doDeserialize(converters []Converter, r io.Reader, ret interface{}, mediaType MediaType) (int, error) {
    l := len(converters)
    for l>0 {
        l--
        if converters[l].CanDeserialize(ret, mediaType) {
            n, err := converters[l].Deserialize(r, ret)
            if err == nil {
                return n, nil
            }
        }
    }
    return 0, errors.New("Cannot Deserialize Object ")
}

func getDefaultMediaType(converter Converter) MediaType {
    return converter.SupportMediaType()[0]
}
