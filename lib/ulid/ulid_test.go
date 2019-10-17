package ulid

import "testing"

func TestGenUlid(t *testing.T) {
	u := New(3)
	res := u.Ulid()
	if res == "" {
		t.Error("result from ulid is empty")
	}
}

func TestMock(t *testing.T) {
	ulidTest := "111122223333"
	um := NewMock(ulidTest)

	if res := um.Ulid(); res != ulidTest {
		t.Error("ulid doesn't match")
	}
}

func TestEmptyMock(t *testing.T) {
	um := NewMock()

	if res := um.Ulid(); res != um.DefaultValue() {
		t.Error("default ulid doesn't exist")
	}
}

func TestMultipleMock(t *testing.T) {
	ulidTest := "111122223333"
	ulidTest2 := "111122224444"

	um := NewMock(ulidTest, ulidTest2)

	if res1 := um.Ulid(); res1 != ulidTest {
		t.Errorf("expect %v, got %v\n", ulidTest, res1)
	}

	if res2 := um.Ulid(); res2 != ulidTest2 {
		t.Errorf("expect %v, got %v\n", ulidTest2, res2)
	}
}
