// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"bytes"
	"github.com/xfali/restclient"
	"reflect"
	"testing"
)

type testStruct struct {
	I int
	S string
	F float64
}

func TestJsonConverter(t *testing.T) {
	t.Run("encoder interface", func(t *testing.T) {
		conv := restclient.NewJsonConverter()
		var in interface{} = &testStruct{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeJson))
		if !b {
			t.Fatal("cannot encode interface to json")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &testStruct{I: 1, S: "a", F: 2.2}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeJson))
		if !b {
			t.Fatal("cannot decode interface to json")
		}
		de := conv.CreateDecoder(buf)
		_, err = de.Decode(in)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(origin, in) {
			t.Log("success! ", in)
		} else {
			t.Fatal("diff ", in, "origin: ", origin)
		}
	})

	t.Run("encoder struct", func(t *testing.T) {
		conv := restclient.NewJsonConverter()
		in := &testStruct{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeJson))
		if !b {
			t.Fatal("cannot encode interface to json")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &testStruct{I: 1, S: "a", F: 2.2}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeJson))
		if !b {
			t.Fatal("cannot decode interface to json")
		}
		de := conv.CreateDecoder(buf)
		_, err = de.Decode(in)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(origin, in) {
			t.Log("success! ", in)
		} else {
			t.Fatal("diff ", in, "origin: ", origin)
		}
	})

	t.Run("encoder map", func(t *testing.T) {
		conv := restclient.NewJsonConverter()
		in := &map[string]interface{}{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeJson))
		if !b {
			t.Fatal("cannot encode interface to json")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		_, err := en.Encode(testStruct{I: 1, S: "a", F: 2.2})
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeJson))
		if !b {
			t.Fatal("cannot decode interface to json")
		}
		de := conv.CreateDecoder(buf)
		_, err = de.Decode(in)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("success! ", in)
	})

	t.Run("encoder slice", func(t *testing.T) {
		conv := restclient.NewJsonConverter()
		in := &[]testStruct{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeJson))
		if !b {
			t.Fatal("cannot encode interface to json")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &[]testStruct{{I: 1, S: "a", F: 2.2}, {I: 2, S: "b", F: 3.3}}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeJson))
		if !b {
			t.Fatal("cannot decode interface to json")
		}
		de := conv.CreateDecoder(buf)
		_, err = de.Decode(in)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(origin, in) {
			t.Log("success! ", in)
		} else {
			t.Fatal("diff ", in, "origin: ", origin)
		}
	})

	t.Run("encoder other", func(t *testing.T) {
		conv := restclient.NewJsonConverter()
		var ss = "123"
		in := &ss
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeJson))
		if !b {
			t.Fatal("cannot encode interface to json")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		_, err := en.Encode([]testStruct{{I: 1, S: "a", F: 2.2}, {I: 2, S: "b", F: 3.3}})
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeJson))
		if b {
			t.Fatal("cannot decode interface to json")
		}
	})
}


func TestYamlConverter(t *testing.T) {
	t.Run("encoder interface", func(t *testing.T) {
		conv := restclient.NewYamlConverter()
		var in interface{} = &testStruct{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeYaml))
		if !b {
			t.Fatal("cannot encode interface to yaml")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &testStruct{I: 1, S: "a", F: 2.2}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeYaml))
		if !b {
			t.Fatal("cannot decode interface to yaml")
		}
		de := conv.CreateDecoder(buf)
		_, err = de.Decode(in)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(origin, in) {
			t.Log("success! ", in)
		} else {
			t.Fatal("diff ", in, "origin: ", origin)
		}
	})

	t.Run("encoder struct", func(t *testing.T) {
		conv := restclient.NewYamlConverter()
		in := &testStruct{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeYaml))
		if !b {
			t.Fatal("cannot encode interface to yaml")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &testStruct{I: 1, S: "a", F: 2.2}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeYaml))
		if !b {
			t.Fatal("cannot decode interface to yaml")
		}
		de := conv.CreateDecoder(buf)
		_, err = de.Decode(in)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(origin, in) {
			t.Log("success! ", in)
		} else {
			t.Fatal("diff ", in, "origin: ", origin)
		}
	})

	t.Run("encoder map", func(t *testing.T) {
		conv := restclient.NewYamlConverter()
		in := &map[string]interface{}{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeYaml))
		if !b {
			t.Fatal("cannot encode interface to yaml")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &testStruct{I: 1, S: "a", F: 2.2}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeYaml))
		if !b {
			t.Fatal("cannot decode interface to yaml")
		}
		de := conv.CreateDecoder(buf)
		_, err = de.Decode(in)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("success! ", in)
	})

	t.Run("encoder slice", func(t *testing.T) {
		conv := restclient.NewYamlConverter()
		in := &[]testStruct{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeYaml))
		if !b {
			t.Fatal("cannot encode interface to yaml")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &[]testStruct{{I: 1, S: "a", F: 2.2}, {I: 2, S: "b", F: 3.3}}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeYaml))
		if !b {
			t.Fatal("cannot decode interface to yaml")
		}
		de := conv.CreateDecoder(buf)
		_, err = de.Decode(in)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(origin, in) {
			t.Log("success! ", in)
		} else {
			t.Fatal("diff ", in, "origin: ", origin)
		}
	})

	t.Run("encoder other", func(t *testing.T) {
		conv := restclient.NewYamlConverter()
		var ss = "123"
		in := &ss
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeYaml))
		if !b {
			t.Fatal("cannot encode interface to yaml")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		_, err := en.Encode([]testStruct{{I: 1, S: "a", F: 2.2}, {I: 2, S: "b", F: 3.3}})
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeYaml))
		if b {
			t.Fatal("cannot decode interface to yaml")
		}
	})
}

func TestXmlConverter(t *testing.T) {
	t.Run("encoder interface", func(t *testing.T) {
		conv := restclient.NewXmlConverter()
		var in interface{} = &testStruct{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeXml))
		if !b {
			t.Fatal("cannot encode interface to xml")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &testStruct{I: 1, S: "a", F: 2.2}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeXml))
		if !b {
			t.Fatal("cannot decode interface to xml")
		}
		de := conv.CreateDecoder(buf)
		_, err = de.Decode(in)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(origin, in) {
			t.Log("success! ", in)
		} else {
			t.Fatal("diff ", in, "origin: ", origin)
		}
	})

	t.Run("encoder struct", func(t *testing.T) {
		conv := restclient.NewXmlConverter()
		in := &testStruct{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeXml))
		if !b {
			t.Fatal("cannot encode interface to xml")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &testStruct{I: 1, S: "a", F: 2.2}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeXml))
		if !b {
			t.Fatal("cannot decode interface to xml")
		}
		de := conv.CreateDecoder(buf)
		_, err = de.Decode(in)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(origin, in) {
			t.Log("success! ", in)
		} else {
			t.Fatal("diff ", in, "origin: ", origin)
		}
	})

	t.Run("encoder map", func(t *testing.T) {
		conv := restclient.NewXmlConverter()
		in := &map[string]interface{}{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeXml))
		if !b {
			t.Fatal("cannot encode interface to xml")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &testStruct{I: 1, S: "a", F: 2.2}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeXml))
		if b {
			t.Fatal("cannot decode interface to xml")
		} else {
			t.Log("xml cannot decode map")
		}
	})

	t.Run("encoder slice", func(t *testing.T) {
		conv := restclient.NewXmlConverter()
		in := &[]testStruct{}
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeXml))
		if !b {
			t.Fatal("cannot encode interface to xml")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		origin := &[]testStruct{{I: 1, S: "a", F: 2.2}, {I: 2, S: "b", F: 3.3}}
		_, err := en.Encode(origin)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeXml))
		if b {
			t.Fatal("cannot decode interface to xml")
		} else {
			t.Log("xml cannot decode to slice")
		}
	})

	t.Run("encoder other", func(t *testing.T) {
		conv := restclient.NewXmlConverter()
		var ss = "123"
		in := &ss
		b := conv.CanEncode(in, restclient.ParseMediaType(restclient.MediaTypeXml))
		if !b {
			t.Fatal("cannot encode interface to xml")
		}
		buf := bytes.NewBuffer(nil)
		en := conv.CreateEncoder(buf)
		_, err := en.Encode([]testStruct{{I: 1, S: "a", F: 2.2}, {I: 2, S: "b", F: 3.3}})
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(buf.String())
		}

		b = conv.CanDecode(in, restclient.ParseMediaType(restclient.MediaTypeXml))
		if b {
			t.Fatal("cannot decode interface to xml")
		}
	})
}
