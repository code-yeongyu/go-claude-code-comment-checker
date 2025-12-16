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

func Test_AgentMemoFilter_IsAgentMemo_Korean_YeogiseoBarwim(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# 여기서 값이 변경됨"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Korean_Guhyeonham(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# 새로운 기능 구현함"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Korean_Chugaham(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# 에러 처리 추가함"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Korean_Sujeongham(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# 버그 수정함"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Korean_ByeongyeongDwaem(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# 여기에서 A로 변경됨"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Korean_Refactoring(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# 리팩토링 진행"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.True(t, result)
}

func Test_AgentMemoFilter_IsAgentMemo_Korean_GijonEneun(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# 기존에는 다르게 동작했음"}

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

func Test_AgentMemoFilter_IsAgentMemo_NotAgentMemo_RegularKorean(t *testing.T) {
	// given
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "# 값의 합계를 계산"}

	// when
	result := filter.IsAgentMemo(comment)

	// then
	assert.False(t, result)
}
