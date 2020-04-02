package structure

import "testing"

func TestVertices_String(t *testing.T) {
	const val = "UnitTest"
	v := &Vertex{ val}

	vs := v.String()
	if val != vs {
		t.Error("Vertex value stored incorrectly")
	}
}
