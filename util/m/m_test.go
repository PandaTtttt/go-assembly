package m

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestM_SortedKeys(t *testing.T) {
	mm := map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"baz": 3,
		"qux": 4,
	}
	assert.Equal(t, M(mm).SortedKeys(), []string{"bar", "baz", "foo", "qux"})
}

func TestM_Contains(t *testing.T) {
	mm := map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"baz": 3,
		"qux": 4,
	}
	assert.Equal(t, M(mm).Contains([]string{"foo"}), true)
	assert.Equal(t, M(mm).Contains([]string{"foo", "bar"}), true)
	assert.Equal(t, M(mm).Contains([]string{"foo", "bar", "baz"}), true)
	assert.Equal(t, M(mm).Contains([]string{"foo", "bar", "baz", "qux"}), true)
	assert.Equal(t, M(mm).Contains([]string{"foo", "bar", "baz", "qux", "quux"}), false)

}

func TestM_ExactlyContains(t *testing.T) {
	mm := map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"baz": 3,
		"qux": 4,
	}
	assert.Equal(t, M(mm).ExactlyContains([]string{"foo"}), false)
	assert.Equal(t, M(mm).ExactlyContains([]string{"foo", "bar"}), false)
	assert.Equal(t, M(mm).ExactlyContains([]string{"foo", "bar", "baz"}), false)
	assert.Equal(t, M(mm).ExactlyContains([]string{"foo", "bar", "baz", "qux"}), true)
	assert.Equal(t, M(mm).ExactlyContains([]string{"foo", "bar", "baz", "qux", "quux"}), false)
}

