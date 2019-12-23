package image_test

import (
	"testing"

	"github.com/albertwidi/go-project-example/internal/entity/image"
)

func TestValidateMode(t *testing.T) {
	cases := []struct {
		mode        image.Mode
		expectError bool
	}{
		{
			mode:        image.Mode("invalid"),
			expectError: true,
		},
		{
			mode:        image.ModePrivate,
			expectError: false,
		},
	}

	for _, c := range cases {
		err := c.mode.Validate()
		if err != nil && !c.expectError {
			t.Errorf("not expecting error but got %v", err)
			return
		}

		if err == nil && c.expectError {
			t.Error("expecting error but got nil")
			return
		}
	}
}

func TestValidateGroup(t *testing.T) {
	cases := []struct {
		group       image.Group
		expectError bool
	}{
		{
			group:       image.Group("invalid"),
			expectError: true,
		},
		{
			group:       image.GroupAmenities,
			expectError: false,
		},
	}

	for _, c := range cases {
		err := c.group.Validate()
		if err != nil && !c.expectError {
			t.Errorf("not expecting error but got %v", err)
			return
		}

		if err == nil && c.expectError {
			t.Error("expecting error but got nil")
			return
		}
	}
}
