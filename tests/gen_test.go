package tests

import "testing"

func TestGenCmd(t *testing.T) {
	// test each flags
	t.Run("all", genCmd_All)
	t.Run("dir", genCmd_Dir)
	t.Run("tags", genCmd_Tags)
	t.Run("recursive", genCmd_Recursive)
	t.Run("taggedFieldOnly", genCmd_TaggedOnly)
	t.Run("MissingTagPolicy", genCmd_MiggingTagValPolicy)
}

func genCmd_All(t *testing.T) {

}

func genCmd_Dir(t *testing.T) {

}

func genCmd_Tags(t *testing.T) {

}

func genCmd_Recursive(t *testing.T) {

}

func genCmd_TaggedOnly(t *testing.T) {

}

func genCmd_MiggingTagValPolicy(t *testing.T) {

}
