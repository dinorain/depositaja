// You can edit this code!
// Click here and start typing.

package depositaja

import "fmt"

//Example 1:
//Input: s = "anagram", t = "nagaram"
//Output: true
//
//Example 2:
//Input: s = "rat", t = "car"
//Output: false
//
//Example 2:
//Input: s = "silent", t = "listen"
//Output: true

func main() {
	fmt.Println("Hello, 世界")
}

func check(s string, t string) bool {
	lnS := len(s)
	ref := make(map[string]int, lnS)

	for _, val := range s {
		v := fmt.Sprint(val)
		if _, ok := ref[v]; !ok {
			ref[v] = 0
		}
		ref[v] += 1
	}

	for _, val := range t {
		v := fmt.Sprint(val)
		if _, ok := ref[v]; !ok {
			return false
		}
		ref[v] -= 1
		if ref[v] < 0 {
			return false
		}
	}

	return true
}
