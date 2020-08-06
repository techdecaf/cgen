package app

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestBump(t *testing.T){
  test := Goblin(t)
  test.Describe("given: an engineer wants to bump the current version of their code", func(){
    test.Describe("when: the increment is invalid", func(){
      test.It("then: it returns an InvalidIncrement error", func(){
        increment := VersionIncrement("test")
        version, err := increment.Bump("0.0.1")

        test.Assert(version).Equal("")
        test.Assert(err.Error()).Equal(`Invalid VersionIncrement wanted major, minor, patch or pre`)
      })
    })

    test.Describe("when: the version is invalid", func(){
      test.It("then: it returns an invalid version error", func(){
        increment := VersionIncrement("patch")
        version, err := increment.Bump("v0.0.1")

        test.Assert(version).Equal("")
        test.Assert(err.Error()).Equal(`Invalid character(s) found in major number "v0"`)
      })
    })

    test.Describe("when: a valid previous version is found", func(){
      test.It("then: it increments the major version", func(){
        increment := VersionIncrement("major")
        version, err := increment.Bump("0.0.1")

        if err != nil {
          test.Fail(err)
        }

        test.Assert(version).Equal("1.0.0")

        version2, err := increment.Bump("1.0.1-2")

        if err != nil {
          test.Fail(err)
        }

        test.Assert(version2).Equal("2.0.0")
      })

      test.It("then: it increments the minor version", func(){
        increment := VersionIncrement("minor")
        version, err := increment.Bump("0.0.1")

        if err != nil {
          test.Fail(err)
        }

        test.Assert(version).Equal("0.1.0")

        version2, err := increment.Bump("0.0.1-3")

        if err != nil {
          test.Fail(err)
        }

        test.Assert(version2).Equal("0.1.0")
      })

      test.It("then: it increments the patch version", func(){
        increment := VersionIncrement("patch")
        version, err := increment.Bump("1.0.1")

        if err != nil {
          test.Fail(err)
        }

        test.Assert(version).Equal("1.0.2")

        version2, err := increment.Bump("0.1.1-rc.1")

        if err != nil {
          test.Fail(err)
        }

        test.Assert(version2).Equal("0.1.2")
      })

      test.It("then: it increments the pre-release version", func(){
        increment := VersionIncrement("pre-release")
        version, err := increment.Bump("0.0.1")

        if err != nil {
          test.Fail(err)
        }

        test.Assert(version).Equal("0.0.1-1")

        version2, err := increment.Bump("0.0.1-1")
        if err != nil {
          test.Fail(err)
        }
        test.Assert(version2).Equal("0.0.1-2")

        version3, err := increment.Bump("0.0.1-rc.1")
        if err != nil {
          test.Fail(err)
        }
        test.Assert(version3).Equal("0.0.1-1")
      })

    })
  })
}