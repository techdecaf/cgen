package tests

import (
	"testing"

	. "github.com/franela/goblin"
	"github.com/techdecaf/cgen/git"
)

func TestGitRepository(t *testing.T){
  test := Goblin(t)
  test.Describe("given: package git", func(){
    test.Describe("when: a new repo is created", func(){
      repo := git.Repository{Directory: "."}

      test.It("then: lists the git tags associated with that directory", func(){
        tags, err := repo.ListTags()
        test.Assert(len(tags)).IsNotZero()
        test.Assert(err).Equal(nil)

        stable, err := tags.LatestStable(); {
          test.Assert(err).Equal(nil)
          test.Assert(stable == "").IsFalse()
        }

        latest, err := tags.Latest(); {
          test.Assert(err).Equal(nil)
          test.Assert(latest == "").IsFalse()
          // test.Assert(latest).Equal("")
        }
      })
    })
  })
}