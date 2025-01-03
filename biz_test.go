package main

import (
	"log"
	"testing"

	"github.com/ccding/go-stun/stun"
)

func TestNattype(t *testing.T) {
	nat, host, err := stun.NewClient().Discover()
	log.Println("local:", nat, host, err)
}
