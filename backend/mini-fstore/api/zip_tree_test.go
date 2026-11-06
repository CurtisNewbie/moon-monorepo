package api

import (
	"encoding/json"
	"testing"
)

func TestZipTreeNodeJSON(t *testing.T) {
	t.Run("folder carries children and no size", func(t *testing.T) {
		b, err := json.Marshal(ZipTreeNode{
			Name:     "docs",
			Dir:      true,
			Children: []*ZipTreeNode{{Name: "a.txt", Size: 5}},
		})
		if err != nil {
			t.Fatalf("Marshal returned error, %v", err)
		}
		want := `{"name":"docs","dir":true,"children":[{"name":"a.txt","dir":false,"size":5}]}`
		if string(b) != want {
			t.Errorf("json = %s, want %s", b, want)
		}
	})

	t.Run("empty folder omits children", func(t *testing.T) {
		b, err := json.Marshal(ZipTreeNode{Name: "docs", Dir: true})
		if err != nil {
			t.Fatalf("Marshal returned error, %v", err)
		}
		want := `{"name":"docs","dir":true}`
		if string(b) != want {
			t.Errorf("json = %s, want %s", b, want)
		}
	})

	t.Run("file carries its size and no children", func(t *testing.T) {
		b, err := json.Marshal(ZipTreeNode{Name: "a.txt", Size: 5})
		if err != nil {
			t.Fatalf("Marshal returned error, %v", err)
		}
		want := `{"name":"a.txt","dir":false,"size":5}`
		if string(b) != want {
			t.Errorf("json = %s, want %s", b, want)
		}
	})

	t.Run("zero-byte file omits size", func(t *testing.T) {
		b, err := json.Marshal(ZipTreeNode{Name: "empty.txt"})
		if err != nil {
			t.Fatalf("Marshal returned error, %v", err)
		}
		want := `{"name":"empty.txt","dir":false}`
		if string(b) != want {
			t.Errorf("json = %s, want %s", b, want)
		}
	})

	t.Run("node round-trips", func(t *testing.T) {
		in := `{"name":"docs","dir":true,"children":[{"name":"a.txt","dir":false,"size":5}]}`
		var n ZipTreeNode
		if err := json.Unmarshal([]byte(in), &n); err != nil {
			t.Fatalf("Unmarshal returned error, %v", err)
		}
		if len(n.Children) != 1 || n.Children[0].Name != "a.txt" || n.Children[0].Size != 5 {
			t.Errorf("decoded node = %+v", n)
		}
	})
}
