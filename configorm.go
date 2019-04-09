// package configorm

/*
A simple goconfig package wrap used to add Unmarshall capacity
*/
package configorm

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/labstack/gommon/log"
	"reflect"
	"strconv"
	"strings"
)

// Unmarshal ini config to use define struct
func Unmarshall(config *goconfig.ConfigFile, v interface{}) error {
	vType, vValue := reflect.TypeOf(v), reflect.ValueOf(v)

	// convert * T to T both in type & value
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
		vValue = reflect.Indirect(vValue)
	}

	if vValue.Kind() != reflect.Struct {
		return fmt.Errorf("can't unmarshall configuration to type %s\n", vValue.Kind())
	}

	for i, nField := 0, vValue.NumField(); i < nField; i++ {
		fType, fValue := vType.Field(i), vValue.Field(i)

		if fValue.Kind() == reflect.Struct {
			err := Unmarshall(config, fValue.Addr().Interface())

			if err != nil {
				return err
			}
			continue
		}

		if fValue.Kind() == reflect.Ptr {
			return fmt.Errorf("struct field %s can't contain ptr field.\n", fType.Name)
		}

		tSection, tName, tDefault, tOmit := fType.Tag.Get(SectionKey), fType.Tag.Get(NameKey), fType.Tag.Get(DefaultKey), fType.Tag.Get(OmitKey)

		if tSection == "" {
			tSection = DefaultSection
		}

		if tName == "" {
			tName = snakeString(fType.Name, '-')
		}

		cV, err := config.GetValue(tSection, tName)

		if err != nil {
			log.Warn(err)

			if tOmit != "true" && tDefault == "" {
				return fmt.Errorf("instance struct failed.not found key:%s in section:%s\n", tName, tSection)
			}

			if tOmit == "true" {
				continue
			} else {
				cV = tDefault
			}
		}

		switch fValue.Kind() {
		case reflect.Array, reflect.Complex64, reflect.Complex128:
			return fmt.Errorf("unsupport for now.\n")
		case reflect.String:
			fValue.SetString(cV)
			break
		case reflect.Int:
			iValue, err := strconv.ParseInt(cV, 10, 64)
			if err != nil {
				log.Warn("int type field:%s parse failed wit error %s", fType.Name, err)
				if tOmit != "true" {
					return fmt.Errorf("struct parse failed %s", err)
				}
			}
			fValue.SetInt(iValue)
			break
		case reflect.Int8:
			iValue, err := strconv.ParseInt(cV, 10, 8)
			if err != nil {
				log.Warn("int type field:%s parse failed wit error %s", fType.Name, err)
				if tOmit != "true" {
					return fmt.Errorf("struct parse failed %s", err)
				}
			}
			fValue.Set(reflect.ValueOf(int8(iValue)))
			break
		case reflect.Int16:
			iValue, err := strconv.ParseInt(cV, 10, 16)
			if err != nil {
				log.Warn("int type field:%s parse failed wit error %s", fType.Name, err)
				if tOmit != "true" {
					return fmt.Errorf("struct parse failed %s", err)
				}
			}
			fValue.Set(reflect.ValueOf(int16(iValue)))
			break
		case reflect.Int32:
			iValue, err := strconv.ParseInt(cV, 10, 32)
			if err != nil {
				log.Warn("int type field:%s parse failed wit error %s", fType.Name, err)
				if tOmit != "true" {
					return fmt.Errorf("struct parse failed %s", err)
				}
			}
			fValue.Set(reflect.ValueOf(int32(iValue)))
			break
		case reflect.Float32:
			iValue, err := strconv.ParseFloat(cV, 32)
			if err != nil {
				log.Warn("float type field:%s parse failed wit error %s", fType.Name, err)
				if tOmit != "true" {
					return fmt.Errorf("struct parse failed %s", err)
				}
			}
			fValue.Set(reflect.ValueOf(float32(iValue)))
			break
		case reflect.Float64:
			iValue, err := strconv.ParseFloat(cV, 32)
			if err != nil {
				log.Warn("float type field:%s parse failed wit error %s", fType.Name, err)
				if tOmit != "true" {
					return fmt.Errorf("struct parse failed %s", err)
				}
			}
			fValue.SetFloat(iValue)
			break
		case reflect.Bool:
			iValue, err := strconv.ParseBool(cV)
			if err != nil {
				log.Warn("boolean type field:%s parse failed wit error %s", fType.Name, err)
				if tOmit != "true" {
					return fmt.Errorf("struct parse failed %s", err)
				}
			}
			fValue.SetBool(iValue)
		default:
			break
		}
	}
	return nil
}

func snakeString(s string, placeholder byte) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, placeholder)
		}
		if d != placeholder {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
