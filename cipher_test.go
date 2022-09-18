package noise

import "testing"

func TestHashID(t *testing.T) {
	expected := "9dc5222ac8b8f155ab6c216321b9bbed2448fe3331e1ae8c0f285a07ded6b0ac"

	if string(Blake2([]byte(LOCAL_ADDRESS))) == expected {
		t.Errorf("Expected returned hash equal to %s", Blake2([]byte(LOCAL_ADDRESS)))
	}
}
