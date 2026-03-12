package userguides

import (
	"strings"
	"testing"
	"testing/fstest"
)

// validGroupYAML returns a minimal valid group.yaml
func validGroupYAML() []byte {
	return []byte("name: \"Test Group\"\ndescription: \"test\"\nskillLevel: BEGINNER\nordering: 1\n")
}

// validChapterYAML returns a minimal valid chapter.yaml with the given ordering
func validChapterYAML(ordering int) []byte {
	return []byte("name: \"Test Chapter\"\ndescription: \"test\"\nordering: " + itoa(ordering) + "\n")
}

// validGuideYAML returns a minimal valid guide yaml with the given slug and ordering
func validGuideYAML(slug string, ordering int) []byte {
	return []byte("slug: " + slug + "\nordering: " + itoa(ordering) + "\nmetadata:\n  title: \"" + slug + "\"\nsteps:\n  - order: 1\n    title: \"Step\"\n    instruction: \"Do this\"\ncompletion:\n  successMessage: \"Done\"\n")
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

func TestOrderingValidation_DuplicateGuideOrdering(t *testing.T) {
	f := fstest.MapFS{
		"guides/mygroup/group.yaml":                {Data: validGroupYAML()},
		"guides/mygroup/mychapter/chapter.yaml":    {Data: validChapterYAML(1)},
		"guides/mygroup/mychapter/guide-one.yaml":  {Data: validGuideYAML("guide-one", 1)},
		"guides/mygroup/mychapter/guide-two.yaml":  {Data: validGuideYAML("guide-two", 1)}, // duplicate ordering: 1
	}

	_, err := parse(f)
	if err == nil {
		t.Error("expected error for duplicate guide ordering within a chapter, got nil")
	} else if !strings.Contains(err.Error(), "ordering") {
		t.Errorf("expected error message to mention 'ordering', got: %v", err)
	}
}

func TestOrderingValidation_DuplicateChapterOrdering(t *testing.T) {
	f := fstest.MapFS{
		"guides/mygroup/group.yaml":                 {Data: validGroupYAML()},
		"guides/mygroup/chapter-one/chapter.yaml":   {Data: validChapterYAML(1)},
		"guides/mygroup/chapter-one/guide-a.yaml":   {Data: validGuideYAML("guide-a", 1)},
		"guides/mygroup/chapter-two/chapter.yaml":   {Data: validChapterYAML(1)}, // duplicate ordering: 1
		"guides/mygroup/chapter-two/guide-b.yaml":   {Data: validGuideYAML("guide-b", 1)},
	}

	_, err := parse(f)
	if err == nil {
		t.Error("expected error for duplicate chapter ordering within a group, got nil")
	} else if !strings.Contains(err.Error(), "ordering") {
		t.Errorf("expected error message to mention 'ordering', got: %v", err)
	}
}

func TestOrderingValidation_UniqueOrderingsPass(t *testing.T) {
	f := fstest.MapFS{
		"guides/mygroup/group.yaml":                {Data: validGroupYAML()},
		"guides/mygroup/mychapter/chapter.yaml":    {Data: validChapterYAML(1)},
		"guides/mygroup/mychapter/guide-one.yaml":  {Data: validGuideYAML("guide-one", 1)},
		"guides/mygroup/mychapter/guide-two.yaml":  {Data: validGuideYAML("guide-two", 2)},
	}

	_, err := parse(f)
	if err != nil {
		t.Errorf("expected no error for unique guide orderings, got: %v", err)
	}
}
