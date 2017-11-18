package bean

import (
	"testing"
	"reflect"
	"github.com/stretchr/testify/assert"
)

type testBean1 struct{
	Bool bool
	String string
	Int int
	Int64 int64
	Float64 float64
	Float32 float32
	Uint uint
	Uint64 uint64
}

func TestSetPrimitiveFieldStringValue(t *testing.T) {
	var tb testBean1

	var v=reflect.ValueOf(&tb)
	var pv=v.Elem()

	assert.Nil(t,SetPrimitiveFieldStringValue(pv.FieldByName("Bool"),"true"),"set bool")
	assert.True(t,tb.Bool==true,"check bool")

	assert.Nil(t,SetPrimitiveFieldStringValue(pv.FieldByName("String"),"string"),"set string")
	assert.True(t,tb.String=="string","check string")

	assert.Nil(t,SetPrimitiveFieldStringValue(pv.FieldByName("Int"),"99"),"set int")
	assert.True(t,tb.Int==99,"check int")

	assert.Nil(t,SetPrimitiveFieldStringValue(pv.FieldByName("Int64"),"99"),"set int64")
	assert.True(t,tb.Int64==99,"check int64")

	assert.Nil(t,SetPrimitiveFieldStringValue(pv.FieldByName("Float64"),"9.9"),"set float64")
	assert.True(t,tb.Float64==9.9,"check float64")

	assert.Nil(t,SetPrimitiveFieldStringValue(pv.FieldByName("Float32"),"9.9"),"set float32")
	assert.True(t,tb.Float32==9.9,"check float32")

	assert.Nil(t,SetPrimitiveFieldStringValue(pv.FieldByName("Uint"),"99"),"set uint")
	assert.True(t,tb.Uint==99,"check uint")

	assert.Nil(t,SetPrimitiveFieldStringValue(pv.FieldByName("Uint64"),"99"),"set uint64")
	assert.True(t,tb.Uint64==99,"check uint64")
}


func TestIsFieldPrimitive(t *testing.T) {
	var tb testBean1
	var v=reflect.ValueOf(&tb)
	var pv=v.Elem()

	assert.True(t,IsFieldPrimitive(pv.FieldByName("Bool")),"isFieldPrimitive bool")
	assert.True(t,IsFieldPrimitive(pv.FieldByName("String")),"isFieldPrimitive String")
	assert.True(t,IsFieldPrimitive(pv.FieldByName("Int")),"isFieldPrimitive Int")
	assert.True(t,IsFieldPrimitive(pv.FieldByName("Int64")),"isFieldPrimitive Int64")
	assert.True(t,IsFieldPrimitive(pv.FieldByName("Float64")),"isFieldPrimitive Float64")
	assert.True(t,IsFieldPrimitive(pv.FieldByName("Float32")),"isFieldPrimitive Float32")
	assert.True(t,IsFieldPrimitive(pv.FieldByName("Uint")),"isFieldPrimitive Uint")
	assert.True(t,IsFieldPrimitive(pv.FieldByName("Uint64")),"isFieldPrimitive Uint64")
}



func TestSetBeanMap(t *testing.T) {
	var tb testBean1

	mapParams:=map[string][]string{
		"bool":[]string{"true"},
		"string":[]string{"hello"},
		"int":[]string{"90"},
		"int64":[]string{"91"},
		"float64":[]string{"92"},
		"float32":[]string{"93"},
		"uint":[]string{"94"},
		"uint64":[]string{"95"},
	}

	err:=SetBeanMap(&tb,mapParams)
	assert.Nil(t,err,"SetBeanMap fail")

	assert.True(t,tb.Bool==true,"check bool")
	assert.True(t,tb.String=="hello","check string")
	assert.True(t,tb.Int==90,"check int")
	assert.True(t,tb.Int64==91,"check int64")
	assert.True(t,tb.Float64==92,"check float64")
	assert.True(t,tb.Float32==93,"check float32")
	assert.True(t,tb.Uint==94,"check uint")
	assert.True(t,tb.Uint64==95,"check uint64")

	//测试require
	type requireStruct struct{
		RequireStr string `require:"true"`
		RegexpInt int `regexp:"^[1-3]{2,2}$"`
		DefaultInt int `default:"88"`
		int `default:"88"`
		AliasInt int `bean:"a"`
	}
	mapParams=map[string][]string{
		"string":[]string{"hello"},
	}
	var rB requireStruct
	err=SetBeanMap(&rB,mapParams)
	assert.NotNil(t,err,"require fail")

	mapParams=map[string][]string{
		"requireStr":[]string{"require"},
		"regexpInt":[]string{"11"},
		"a":[]string{"99"},
	}
	err=SetBeanMap(&rB,mapParams)
	assert.Nil(t,err,"regexp fail")
	assert.True(t,rB.RequireStr=="require","requireStr")
	assert.True(t,rB.RegexpInt==11,"regexpInt")
	assert.True(t,rB.DefaultInt==88,"defaultInt")
	assert.True(t,rB.AliasInt==99,"alias fail")

	mapParams=map[string][]string{
		"requireStr":[]string{"require"},
		"regexpInt":[]string{"88"},
	}
	err=SetBeanMap(&rB,mapParams)
	assert.NotNil(t,err,"regexp check fail")

	//loop
	type bag struct{
		Name string
		Price float64 `bean:"p"`
	}

	type student struct{
		Name string
		Bag *bag `require:"true"`
	}
	mapParams=map[string][]string{
		"name":[]string{"abc"},
		"bag.name":[]string{"haha"},
		"bag.p":[]string{"99.9"},
	}
	var sb student
	err=SetBeanMap(&sb,mapParams)
	assert.Nil(t,err,"student setBean fail")
	assert.True(t,sb.Name=="abc")
	assert.True(t,sb.Bag.Name=="haha","set ptr field")
	assert.True(t,sb.Bag.Price==99.9)
}