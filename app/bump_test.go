package app

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestBump(t *testing.T){
  test := Goblin(t)
  test.Describe("given: an engineer wants to bump the current version of their code", func(){
    test.Describe("when: a valid previous version is found", func(){
      test.It("then: it increments the minor version", func(){
        increment := VersionIncrement(minor)
        version, err := increment.Bump("0.0.1")

        if err != nil {
          test.Fail(err)
        }

        test.Assert(version).Equal("0.1.0")
      })

      test.It("then: it increments the pre-release version", func(){
        increment := VersionIncrement(pre)
        version, err := increment.Bump("0.0.1-1.pre")

        if err != nil {
          test.Fail(err)
        }

        test.Assert(version).Equal("0.0.1-2")
      })
    })
  })
}