package mypkg

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//
// Most common scalars
//

// type YesNo bool

// // UnmarshalGQL implements the graphql.Unmarshaler interface
// func (y *YesNo) UnmarshalGQL(v interface{}) error {
// 	yes, ok := v.(string)
// 	if !ok {
// 		return fmt.Errorf("YesNo must be a string")
// 	}

// 	if yes == "yes" {
// 		*y = true
// 	} else {
// 		*y = false
// 	}
// 	return nil
// }

// // MarshalGQL implements the graphql.Marshaler interface
// func (y YesNo) MarshalGQL(w io.Writer) {
// 	if y {
// 		w.Write([]byte(`"yes"`))
// 	} else {
// 		w.Write([]byte(`"no"`))
// 	}
// }

//
// Scalars that need access to the request context
//

// type Length float64

// // UnmarshalGQLContext implements the graphql.ContextUnmarshaler interface
// func (l *Length) UnmarshalGQLContext(ctx context.Context, v interface{}) error {
// 	s, ok := v.(string)
// 	if !ok {
// 		return fmt.Errorf("Length must be a string")
// 	}
// 	length, err := ParseLength(s)
// 	if err != nil {
// 		return err
// 	}
// 	*l = length
// 	return nil
// }

// // MarshalGQLContext implements the graphql.ContextMarshaler interface
// func (l Length) MarshalGQLContext(ctx context.Context, w io.Writer) error {
// 	s, err := l.FormatContext(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	w.Write([]byte(strconv.Quote(s)))
// 	return nil
// }

// // ParseLength parses a length measurement string with unit on the end (eg: "12.45in")
// func ParseLength(string) (Length, error)

// // ParseLength formats the string using a value in the context to specify format
// func (l Length) FormatContext(ctx context.Context) (string, error)

// Myarray is a custom type for an array of strings
type Myarray []string

func (ma *Myarray) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("src value cannot be cast to []byte")
	}
	*ma = strings.Split(string(bytes), ",")
	return nil
}

func (ma Myarray) Value() (driver.Value, error) {
	if len(ma) == 0 {
		return nil, nil
	}
	return strings.Join(ma, ","), nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for Myarray
func (ma *Myarray) UnmarshalGQL(v interface{}) error {
	strArr, ok := v.([]interface{})
	if !ok {
		return fmt.Errorf("Myarray must be an array of strings")
	}

	*ma = make([]string, len(strArr))
	for i, elem := range strArr {
		str, ok := elem.(string)
		if !ok {
			return fmt.Errorf("Myarray element must be a string")
		}
		(*ma)[i] = str
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface for Myarray
func (ma Myarray) MarshalGQL(w io.Writer) {
	strArr := make([]string, len(ma))
	for i, str := range ma {
		strArr[i] = strconv.Quote(str)
	}
	io.WriteString(w, "[")
	io.WriteString(w, strings.Join(strArr, ","))
	io.WriteString(w, "]")
}
