package main

import "testing"
import "fmt"
import "reflect"

func TestSchema(x *testing.T) {
	fmt.Println("TestSchema")


	u := params{1,"Aluminum","Engrave","Brad",""}
	t := reflect.TypeOf(u)

	for _, fieldName := range []string{"Name", "Email"} {
		field, found := t.FieldByName(fieldName)
		if !found {
			continue
		}
		fmt.Printf("\nField: User.%s\n", fieldName)
		fmt.Printf("\tWhole tag value : %q\n", field.Tag)
		fmt.Printf("\tValue of 'mytag': %q\n", field.Tag.Get("mytag"))
	}

	m := make(map[string]string)
	fmt.Printf("DEST MAP IS %q\n",m)
	structToMap(u,m,"")
	fmt.Printf("DEST MAP IS %q\n",m)
}


func xxx() {
	type User struct {
		Name  string `mytag:"MyName"`
		Email string `mytag:"MyEmail"`
	}

	u := User{"Bob", "bob@mycompany.com"}
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(&u).Elem()

	fmt.Printf("\nV is %T %+v\n", v, v)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		fmt.Printf("\nField Name %s %s TAG:%s %T %q\n", field.Name, field.Type, field.Tag.Get("mytag"), value, value)
		v.Field(i).SetString("XX")
		/*
			fmt.Printf("\nField Type %T\n", field)
			fmt.Printf("\nValue %T %+v\n", value, value)
			fmt.Printf("\nType  %s\n", field.Type)
			fmt.Printf("\tWhole tag value : %q\n", field.Tag)
			fmt.Printf("\tValue of 'mytag': %q\n", field.Tag.Get("mytag"))
			reflect.ValueOf(&u).Elem().SetString("XXX")
		*/
	}

	fmt.Printf("U IS NOW %+v\n", u)
}
