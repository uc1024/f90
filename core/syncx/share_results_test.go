package syncx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShareCall(t *testing.T) {
	shareResults := NewShareCall()
	// Test case 1: Perform a Do operation
	go func() {
		result1, done1, err1 := shareResults.DoEx("key1", func() (interface{}, error) {
			time.Sleep(time.Second * 2)
			return "result1", nil
		})
		assert.True(t, done1)
		assert.NoError(t, err1)
		assert.Equal(t, "result1", result1)
	}()
	go func() {
		time.Sleep(time.Millisecond * 500)
		// Test case 2: Perform a Do operation with the same key, it should reuse the result from the first call
		result2, err2 := shareResults.Do("key1", func() (interface{}, error) {
			t.Fatal("This function should not be called")
			return "result2", nil
		})
		assert.NoError(t, err2)
		t.Log(result2)
		assert.Equal(t, "result1", result2)
	}()
	go func() {
		time.Sleep(time.Millisecond * 500)
		// Test case 3: Perform a DoEx operation, should return true since the result is already available
		result3, done3, err3 := shareResults.DoEx("key1", func() (interface{}, error) {
			t.Fatal("This function should not be called")
			return "result3", nil
		})
		assert.NoError(t, err3)
		assert.True(t, !done3)
		t.Log(result3)
		assert.Equal(t, "result1", result3)
	}()
	cond := NewCond()
	cond.WaitWithTimeout(time.Second * 3)
}
