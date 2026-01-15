package userguides_test

import (
	"testing"

	userguides "github.com/spacelift-io/spacelift-user-guides-library"
)

func TestGuidesLoad(t *testing.T) {
	lib, err := userguides.Guides()
	if err != nil {
		t.Fatalf("Guides() returned error: %v", err)
	}

	if lib == nil {
		t.Fatal("Guides() returned nil library")
	}

	if len(lib.Groups) == 0 {
		t.Fatal("Expected at least one group, got none")
	}

	t.Logf("Successfully loaded %d group(s)", len(lib.Groups))

	for _, group := range lib.Groups {
		t.Logf("Group: %s (%s) - %d chapter(s)", group.Name, group.Slug, len(group.Chapters))

		if group.Name == "" {
			t.Errorf("Group %s has empty name", group.Slug)
		}

		if group.SkillLevel == "" {
			t.Errorf("Group %s has empty skill level", group.Slug)
		}

		for _, chapter := range group.Chapters {
			t.Logf("  Chapter: %s (%s) - %d guide(s)", chapter.Name, chapter.Slug, len(chapter.Guides))

			if chapter.Name == "" {
				t.Errorf("Chapter %s has empty name", chapter.Slug)
			}

			for _, guide := range chapter.Guides {
				t.Logf("    Guide: %s (%s) - %d step(s)", guide.Metadata.Title, guide.Slug, len(guide.Steps))

				if guide.Metadata.Title == "" {
					t.Errorf("Guide %s has empty title", guide.Slug)
				}

				if len(guide.Steps) == 0 {
					t.Errorf("Guide %s has no steps", guide.Slug)
				}

				for _, step := range guide.Steps {
					if step.Order <= 0 {
						t.Errorf("Guide %s step has invalid order: %d", guide.Slug, step.Order)
					}
					if step.Title == "" {
						t.Errorf("Guide %s step %d has empty title", guide.Slug, step.Order)
					}
					if step.Instruction == "" {
						t.Errorf("Guide %s step %d has empty instruction", guide.Slug, step.Order)
					}
				}
			}
		}
	}
}

func TestGettingStartedGroup(t *testing.T) {
	lib, err := userguides.Guides()
	if err != nil {
		t.Fatalf("Guides() returned error: %v", err)
	}

	var gettingStarted *userguides.Group
	for i := range lib.Groups {
		if lib.Groups[i].Slug == "getting-started" {
			gettingStarted = &lib.Groups[i]
			break
		}
	}

	if gettingStarted == nil {
		t.Fatal("Expected 'getting-started' group to exist")
	}

	if gettingStarted.Name != "Getting Started" {
		t.Errorf("Expected name 'Getting Started', got %q", gettingStarted.Name)
	}

	if gettingStarted.SkillLevel != "BEGINNER" {
		t.Errorf("Expected skill level 'BEGINNER', got %q", gettingStarted.SkillLevel)
	}

	if len(gettingStarted.Chapters) == 0 {
		t.Fatal("Expected at least one chapter in 'getting-started' group")
	}
}

func TestSlugGeneration(t *testing.T) {
	lib, err := userguides.Guides()
	if err != nil {
		t.Fatalf("Guides() returned error: %v", err)
	}

	for _, group := range lib.Groups {
		if group.Slug == "" {
			t.Errorf("Group %s has empty slug", group.Name)
		}

		for _, chapter := range group.Chapters {
			if chapter.Slug == "" {
				t.Errorf("Chapter %s has empty slug", chapter.Name)
			}

			for _, guide := range chapter.Guides {
				if guide.Slug == "" {
					t.Errorf("Guide %s has empty slug", guide.Metadata.Title)
				}
			}
		}
	}
}

func TestOrdering(t *testing.T) {
	lib, err := userguides.Guides()
	if err != nil {
		t.Fatalf("Guides() returned error: %v", err)
	}

	for _, group := range lib.Groups {
		t.Logf("Group %s has ordering: %d", group.Name, group.Ordering)

		for _, chapter := range group.Chapters {
			t.Logf("  Chapter %s has ordering: %d", chapter.Name, chapter.Ordering)

			for _, guide := range chapter.Guides {
				t.Logf("    Guide %s has ordering: %d", guide.Metadata.Title, guide.Ordering)
			}
		}
	}
}
