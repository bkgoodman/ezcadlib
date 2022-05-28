package main
import "testing"
import "fmt"


func TestBkg(t *testing.T) {
	var rmap =make(map[string]string)
	rmap["name"]="Franklin"
	rmap["color"]="red"
	rmap["food"]="strawberry"
	result :=replaceMap("My pets name is ${name}, favorite color ${color}, eats ${food}, ${age} years old",&rmap)
	fmt.Println(result)
}
