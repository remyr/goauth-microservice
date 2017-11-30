package utils

import (
	"log"
	"strings"
	"regexp"
)

func ParseDuplicateKey() {
	var e = "E11000 duplicate key error collection: goauth_microservice.users index: email_1 dup key: { : \"rivals.remy@gmail.com\" }"
	var stp1 = strings.Split(e, "index: ")[1]
	var stp2 = strings.Split(stp1, " dup key")[0]
	var i = strings.Index(stp2, "_")
	var result = stp2[:i]

	r, _ := regexp.Compile(`index:\s([a-z]+).*{\s:\s\\"([a-z.@" "]+)`)
	res := r.FindAllString(e, -1)

	log.Printf("%s", result)
	log.Printf("%s", res)
}