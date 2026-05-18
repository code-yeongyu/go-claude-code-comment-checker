package filters

import (
	"testing"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
	"github.com/stretchr/testify/assert"
)

func Test_AgentMemoFilter_IsAgentMemo_ChangedFrom(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Changed from old_value to new_value"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_ModifiedTo(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# Modified to use new implementation"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_UpdatedFrom(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Updated from v1 to v2"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Refactored(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Refactored for better performance"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Added(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Added new validation logic"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Removed(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Removed deprecated function"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Implemented(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Implemented new feature"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_ThisImplements(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// This implements the new API"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_HereWe(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Here we handle the error case"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_NowThis(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Now this uses the new format"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Previously(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Previously this was handled differently"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Note(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Note: this is important"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_ArrowNotation(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// oldValue -> newValue"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_ConvertedFrom(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# Converted from callbacks to state machine"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_MigratedTo(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# Migrated to the new parser"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_SwitchedFrom(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# Switched from old code path"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Replaced(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# Replaced temporary formatter"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Deleted(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# Deleted stale branch"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_BeforeThis(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# Before this the parser was recreated"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_AfterThis(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# After this the cache is warm"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_NotAgentMemo_BDD(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# given"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.False(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_NotAgentMemo_Directive(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# noqa: E501"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.False(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_NotAgentMemo_Regular(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Calculate the sum of values"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.False(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_NotAgentMemo_RegularEnglish(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# Calculate the checksum before write"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.False(t, result)
}
