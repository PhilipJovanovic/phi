package phi

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type TestStruct struct {
	ClassicString string `json:"classicString"`
	ClassicInt    int    `json:"classicInt"`
}

type TestRequiredStruct struct {
	ClassicString string `json:"classicString,required"`
	ClassicInt    int    `json:"classicInt"`
}

type TestComplexStruct struct {
	ClassicString      string                            `json:"classicString,required"`
	ClassicStruct      TestRequiredStruct                `json:"classicStruct"`
	ClassicPtr         *TestRequiredStruct               `json:"classicPtr,required"`
	ClassicSlice       []TestRequiredStruct              `json:"classicSlice,required"`
	ClassicPtrSlice    []*TestRequiredStruct             `json:"classicPtrSlice,required"`
	ClassicPtrPtrSlice []**TestRequiredStruct            `json:"classicPtrPtrSlice,required"`
	ClassicPtrSlicePtr *[]*TestRequiredStruct            `json:"classicPtrSlicePtr,required"`
	ClassicMap         map[string]TestRequiredStruct     `json:"classicMap,required"`
	ClassicPtrMap      map[string]*TestRequiredStruct    `json:"classicPtrMap,required"`
	ClassicPtrMapPtr   *map[string]*TestRequiredStruct   `json:"classicPtrMapPtr,required"`
	ClassicMapSlicePtr *map[string][]*TestRequiredStruct `json:"classicMapSlicePtr,required"`
}

type TestWeirdStruct struct {
	DumbString *string                                  `json:"dumbString,required"`
	WeirdMap   *map[int]TestRequiredStruct              `json:"weirdMap,required"`
	CrackSlice *[]*[]*TestRequiredStruct                `json:"crackSlice,required"`
	WtfMap     map[string]map[string]TestRequiredStruct `json:"wtfMap,required"`
}

func TestRequest(t *testing.T) {
	t.Run("test basic", TestBasic)
	t.Run("test missing param", TestMissingParam)
	t.Run("test slice", TestSlice)
	t.Run("test complex", TestComplex)
	t.Run("test weird", TestWeird)
	t.Run("test redundant ptr", TestRedPtr)
	t.Run("test playground", TestPlayground)
}

func TestBasic(t *testing.T) {
	body, err := handleValidate(&TestStruct{
		ClassicString: "test",
		ClassicInt:    1337,
	})
	if err != nil {
		t.Error("error should be nil")
	}

	if body == nil {
		t.Error("body is nil")
	}
}

func TestMissingParam(t *testing.T) {
	_, err := handleValidate(&TestRequiredStruct{
		ClassicInt: 1337,
	})
	if err == nil || err.Error != "missingBodyParameters" {
		t.Error("error should be missingBodyParameters")
	}
}

func TestSlice(t *testing.T) {
	_, err := handleValidate(&[]TestRequiredStruct{
		{ClassicString: "test1", ClassicInt: 1337},
		{ClassicString: "test1", ClassicInt: 1338},
	})
	if err != nil {
		t.Error(err.Message)
	}

	_, err = handleValidate(&[2]TestRequiredStruct{
		{ClassicString: "test1", ClassicInt: 1337},
		{ClassicInt: 1338},
	})
	if err == nil {
		t.Error("error should not be nil")
		return
	}

	if err.Message != "missing '[1].classicString'" {
		t.Errorf("validation went wrong have (%s) want (%s)", err.Message, "missing '[1].classicString'")
	}
}

func TestComplex(t *testing.T) {
	p1 := &TestRequiredStruct{ /* ClassicString: "test", */ ClassicInt: 1337}
	p2 := &TestRequiredStruct{ClassicString: "test", ClassicInt: 1338}

	_, err := handleValidate(&TestComplexStruct{
		/* ClassicString: "test", */
		ClassicStruct: TestRequiredStruct{
			/* ClassicString: "test", */
			ClassicInt: 1337,
		},
		ClassicPtr: &TestRequiredStruct{
			/* ClassicString: "test", */
			ClassicInt: 1337,
		},
		ClassicSlice: []TestRequiredStruct{
			{ /* ClassicString: "test", */ ClassicInt: 1337},
			{ClassicString: "test", ClassicInt: 1338},
		},
		ClassicPtrSlice: []*TestRequiredStruct{
			{ /* ClassicString: "test", */ ClassicInt: 1337},
			{ClassicString: "test", ClassicInt: 1338},
		},
		ClassicPtrSlicePtr: &[]*TestRequiredStruct{
			{ /* ClassicString: "test", */ ClassicInt: 1337},
			{ClassicString: "test", ClassicInt: 1338},
		},
		ClassicPtrPtrSlice: []**TestRequiredStruct{
			&p1,
			&p2,
		},
		ClassicMap: map[string]TestRequiredStruct{
			"test1": { /* ClassicString: "test", */ ClassicInt: 1337},
			"test2": {ClassicString: "test", ClassicInt: 1338},
		},
		ClassicPtrMap: map[string]*TestRequiredStruct{
			"test1": { /* ClassicString: "test", */ ClassicInt: 1337},
			"test2": {ClassicString: "test", ClassicInt: 1338},
		},
		ClassicPtrMapPtr: &map[string]*TestRequiredStruct{
			"test1": { /* ClassicString: "test", */ ClassicInt: 1337},
			"test2": {ClassicString: "test", ClassicInt: 1338},
		},
		ClassicMapSlicePtr: &map[string][]*TestRequiredStruct{
			"test1": {
				{ /* ClassicString: "test", */ ClassicInt: 1337},
				{ClassicString: "test", ClassicInt: 1338},
			},
			"test2": {
				{ClassicString: "test", ClassicInt: 1337},
				{ /* ClassicString: "test", */ ClassicInt: 1338},
			},
		},
	})
	errorS := "missing 'classicString, classicStruct.classicString, classicPtr.classicString, classicSlice[0].classicString, classicPtrSlice[0].classicString, classicPtrPtrSlice[0].classicString, classicPtrSlicePtr[0].classicString, classicMap[test1].classicString, classicPtrMap[test1].classicString, classicPtrMapPtr[test1].classicString, classicMapSlicePtr[test1][0].classicString, classicMapSlicePtr[test2][1].classicString'"
	if err == nil || err.Message != errorS {
		t.Errorf("validation went wrong have (%s) want (%s)", err.Message, errorS)

	}
}

func TestWeird(t *testing.T) {
	_, err := handleValidate(&TestWeirdStruct{
		DumbString: nil,
		WeirdMap:   nil,
		CrackSlice: nil,
		WtfMap:     nil,
	})
	if err == nil || err.Message != "missing 'dumbString, weirdMap, crackSlice, wtfMap'" {
		t.Errorf("validation went wrong have (%s) want (%s)", err.Message, "missing 'dumbString, weirdMap, crackSlice, wtfMap'")
	}

	_, err = handleValidate(&TestWeirdStruct{
		DumbString: nil,
		WeirdMap: &map[int]TestRequiredStruct{
			1337: { /* ClassicString: "test2", */ ClassicInt: 1338},
		},
		CrackSlice: &[]*[]*TestRequiredStruct{
			{{ClassicString: "test2", ClassicInt: 1337}, {ClassicString: "test2", ClassicInt: 1338}},
			{{ClassicString: "test2", ClassicInt: 1338}, { /* ClassicString: "test2", */ ClassicInt: 1338}},
		},
		WtfMap: map[string]map[string]TestRequiredStruct{
			"test1": {
				"test2": { /* ClassicString: "test2", */ ClassicInt: 1338},
			},
		},
	})
	errS := "missing 'dumbString, weirdMap[1337].classicString, crackSlice[1][1].classicString, wtfMap[test1][test2].classicString'"
	if err == nil || err.Message != errS {
		t.Errorf("validation went wrong have (%s) want (%s)", err.Message, errS)
	}
}

func TestRedPtr(t *testing.T) {
	e := &TestRequiredStruct{
		ClassicInt: 1337,
	}

	_, err := handleValidate(&e)
	if err == nil {
		t.Error("error should not be nil")
		return
	}

	if err.Message != "missing 'classicString'" {
		t.Errorf("validation went wrong have (%s) want (%s)", err.Message, "missing 'classicString'")
	}
}

type TestPlaygroundStruct struct {
	Arr *[]***TestRequiredStruct `json:"arr,required"`
}

func TestPlayground(t *testing.T) {
	p1 := &TestRequiredStruct{ /* ClassicString: "test", */ ClassicInt: 1337}
	p2 := &TestRequiredStruct{ClassicString: "test", ClassicInt: 1338}

	p11 := &p1
	p12 := &p2

	data := []***TestRequiredStruct{
		&p11, &p12,
	}
	_, err := handleValidate(&TestPlaygroundStruct{
		Arr: &data,
	})
	spew.Dump(err)
}
