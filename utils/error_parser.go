package utils

import (
	"regexp"
	"fmt"
)

func ParseDupError() {
	//var re = regexp.MustCompile(`index:\s(?P<Key>[a-z]+).*{\s:\s\\"(?P<Value>[a-z.@" "]+)`)
    var str = `E11000 duplicate key error collection: goauth_microservice.users index: email_1 dup key: { : \"rivals.remy@gmail.com\" }`
	//
    ////for i, match := range re.FindAllString(str, -1) {
    ////    fmt.Println(match, "found at index", i)
    ////}
    //fmt.Println(re.FindStringSubmatch(str))
    params := getParams(`index:\s(?P<Key>[a-z]+).*{\s:\s\\"(?P<Value>[a-z.@" "]+)`, str)
	fmt.Println(params["Key"])
	fmt.Println(params["Value"])
}

func getParams(regEx, url string) (paramsMap map[string]string) {

    var compRegEx = regexp.MustCompile(regEx)
    match := compRegEx.FindStringSubmatch(url)

    paramsMap = make(map[string]string)
    for i, name := range compRegEx.SubexpNames() {
        if i > 0 && i <= len(match) {
            paramsMap[name] = match[i]
        }
    }
    return
}