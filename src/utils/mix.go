package utils

import (
	"log"
	"strconv"
	"time"
)

// LogErr logs error if there is one
func LogErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Mute wrap multi-values from a function call to access specific returned value by specifying an index
func Mute(a ...interface{}) []interface{} {
	return a
}

// MuteStr return first value as string
func MuteStr(a ...interface{}) string {
	return a[0].(string)
}

// DoEvery call function every x seconds
func DoEvery(d int, fnc func()) chan struct{} {
	ticker := time.NewTicker(time.Duration(d) * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fnc()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

// Sec2Time converts seconds into time format hh:mm:ss
func Sec2Time(ms int64) string {
	min := ms / 60
	sek := ms % 60
	hour := min / 60
	if hour != 0 {
		min = min % 60
	}
	return strconv.FormatInt(hour, 10) + ":" + strconv.FormatInt(min, 10) + ":" + strconv.FormatInt(sek, 10)
}

/*
// DoEvery2 call function every x seconds and send return value to channel
func DoEvery2(d int, fnc func() interface{}) (*time.Ticker, chan interface{}) {
	ch := make(chan interface{})
	ticker := time.NewTicker(time.Duration(d) * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				ch <- fnc()
			}
		}
	}()

	return ticker, ch
}

// DoEvery3 call function via reflection every x seconds and send return value to channel
func DoEvery3(d int, myClass interface{}, funcName string, params ...interface{}) (*time.Ticker, chan interface{}) {
	ch := make(chan interface{})
	ticker := time.NewTicker(time.Duration(d) * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				ch <- CallFuncByName(myClass, funcName, params...)
			}
		}
	}()

	return ticker, ch
}

func testDoEvery2() {
	checkFailure := func(str1 string) func() interface{} {
		str2 := "str2"
		return func() interface{} {
			return str1 + str2
		}
	}

	_, ch := DoEvery2(1, checkFailure("str1"))
	DoEvery(1, func() { fmt.Println("res = ", <-ch) })
}

type rect struct {
	width, height int
}

func (r rect) Area(x int) int {
	return r.width * r.height * x
}

func testDoEvery3() {
	r := rect{width: 10, height: 5}
	refArgs := make([]interface{}, 1)
	refArgs[0] = 2
	_, ch := DoEvery3(1, &r, "Area", refArgs...)
	DoEvery(1, func() { fmt.Println("res = ", (<-ch).([]reflect.Value)[0].Interface().(int)) })
}

// GetStructField get field of struct via reflection
func GetStructField(s interface{}, field string) interface{} {
	r := reflect.ValueOf(s)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Interface()
}

// CallFuncByName call function by name via reflection
func CallFuncByName(myClass interface{}, funcName string, params ...interface{}) []reflect.Value {
	myClassValue := reflect.ValueOf(myClass)
	m := myClassValue.MethodByName(funcName)
	if !m.IsValid() {
		fmt.Println("Method not found " + funcName)
		return make([]reflect.Value, 0)
	}
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}
	return m.Call(in)
}
*/
