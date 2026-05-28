package form

import (
    "testing"

    assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestPointer_Bitmasks(t *testing.T) {
    tests := []struct {
        name     string
        mask     pointer
        checkAny []pointer
        wantAny  bool
        wantNone bool
    }{
        {
            name:     "Prompt active has Prompt",
            mask:     pointerPrompt,
            checkAny: []pointer{pointerPrompt},
            wantAny:  true,
            wantNone: false,
        },
        {
            name:     "Prompt active does not have Gutter",
            mask:     pointerPrompt,
            checkAny: []pointer{pointerGutter},
            wantAny:  false,
            wantNone: true,
        },
        {
            name:     "Combo active has any of them",
            mask:     pointerPrompt | pointerGutter,
            checkAny: []pointer{pointerGutter},
            wantAny:  true,
            wantNone: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.wantAny, tt.mask.hasAny(tt.checkAny...))
            assert.Equal(t, tt.wantNone, tt.mask.hasNone(tt.checkAny...))
        })
    }
}

func TestPointer_Navigation(t *testing.T) {
    t.Run("findPointer bounds checking", func(t *testing.T) {
        assert.Equal(t, pointerPrompt, findPointer(0))
        assert.Equal(t, pointerGutter, findPointer(1))

        assert.Equal(t, pointerPrompt, findPointer(3)) 
        assert.Equal(t, pointerPrompt, findPointer(255))
    })

    t.Run("nextPointer cycling logic", func(t *testing.T) {
        assert.Equal(t, uint8(1), nextPointer(0))
        assert.Equal(t, uint8(2), nextPointer(1))
        assert.Equal(t, uint8(0), nextPointer(2))
    })
}
