package custom

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Date is a custom time type with a custom unmarshaler for ISO 8601 date format
type Date time.Time

func (d *Date) String() string {
	t := time.Time(*d).Format("2006-01-02")
	return t
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d *Date) ToTime() time.Time {
	return time.Time(*d)
}

var timeConverter = func(value string) reflect.Value {
	fmt.Println("timeConverter", value)
	if v, err := time.Parse("2006-01-02", value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.Value{}
}

var CustomDate = fiber.ParserType{
	Customtype: Date{},
	Converter:  timeConverter,
}
