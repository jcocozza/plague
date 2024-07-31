package main

// check if a char is 0,1,2,3,4,5,6,7,8,9
func isNumber(char string) bool {
    if len(char) > 1 {
        panic("unable to check char of len > 1")
    }
    if char == "" {
        return false
    }

    rn := []rune(char)[0]
    return rn >= '0' && rn <= '9'
}

// check if a char is a-z (only lowercase)
func isLetter(char string) bool {
    if len(char) > 1 {
        panic("unable to check char of len > 1")
    }
    if char == "" {
        return false
    }

    rn := []rune(char)[0]
    return rn >= 'a' && rn <= 'z'
}
