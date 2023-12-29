package core

import (
	"context"
	"log"
	"testing"
)

func TestResolverFunc(t *testing.T) {
	resolver := newResolverImpl(defaultLogger)
	domain, err := resolver.Resolve(context.TODO(),
		[]string{
			"sd-rtn.com",
			"agora.io",
		}, "api")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("domain:%s", domain)
}
