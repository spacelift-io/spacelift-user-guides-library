package userguides_test

import (
	"strings"
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

func TestFoundationsGroup(t *testing.T) {
	lib, err := userguides.Guides()
	if err != nil {
		t.Fatalf("Guides() returned error: %v", err)
	}

	var foundations *userguides.Group
	for i := range lib.Groups {
		if lib.Groups[i].Slug == "foundations" {
			foundations = &lib.Groups[i]
			break
		}
	}

	if foundations == nil {
		t.Fatal("Expected 'foundations' group to exist")
	}

	if foundations.Name != "Foundations" {
		t.Errorf("Expected name 'Foundations', got %q", foundations.Name)
	}

	if foundations.SkillLevel != "BEGINNER" {
		t.Errorf("Expected skill level 'BEGINNER', got %q", foundations.SkillLevel)
	}

	if len(foundations.Chapters) == 0 {
		t.Fatal("Expected at least one chapter in 'foundations' group")
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

func TestValidationRules(t *testing.T) {
	tests := []struct {
		name      string
		guide     userguides.Guide
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid guide",
			guide: userguides.Guide{
				Slug:     "test-guide",
				Ordering: 1,
				Metadata: userguides.GuideMetadata{
					Title:             "Test Guide",
					Description:       "A test guide",
					Labels:            []string{"test"},
					Difficulty:        "easy",
					MinutesToComplete: 5,
				},
				Steps: []userguides.GuideStep{
					{Order: 1, Title: "Step 1", Instruction: "Do this"},
					{Order: 2, Title: "Step 2", Instruction: "Do that"},
				},
				Completion: userguides.GuideCompletion{
					SuccessMessage:      "Done",
					RecommendedGuideIDs: []string{},
				},
			},
			expectErr: false,
		},
		{
			name: "invalid difficulty",
			guide: userguides.Guide{
				Slug:     "test-guide",
				Ordering: 1,
				Metadata: userguides.GuideMetadata{
					Title:             "Test Guide",
					Difficulty:        "super-hard",
					MinutesToComplete: 5,
				},
				Steps: []userguides.GuideStep{
					{Order: 1, Title: "Step 1", Instruction: "Do this"},
				},
				Completion: userguides.GuideCompletion{},
			},
			expectErr: true,
			errMsg:    "invalid difficulty",
		},
		{
			name: "empty label",
			guide: userguides.Guide{
				Slug:     "test-guide",
				Ordering: 1,
				Metadata: userguides.GuideMetadata{
					Title:             "Test Guide",
					Labels:            []string{"valid", ""},
					MinutesToComplete: 5,
				},
				Steps: []userguides.GuideStep{
					{Order: 1, Title: "Step 1", Instruction: "Do this"},
				},
				Completion: userguides.GuideCompletion{},
			},
			expectErr: true,
			errMsg:    "label at index 1 is empty",
		},
		{
			name: "non-sequential steps",
			guide: userguides.Guide{
				Slug:     "test-guide",
				Ordering: 1,
				Metadata: userguides.GuideMetadata{
					Title:             "Test Guide",
					MinutesToComplete: 5,
				},
				Steps: []userguides.GuideStep{
					{Order: 1, Title: "Step 1", Instruction: "Do this"},
					{Order: 3, Title: "Step 3", Instruction: "Skip step 2"},
				},
				Completion: userguides.GuideCompletion{},
			},
			expectErr: true,
			errMsg:    "steps must be sequentially ordered",
		},
		{
			name: "invalid URL scheme",
			guide: userguides.Guide{
				Slug:     "test-guide",
				Ordering: 1,
				Metadata: userguides.GuideMetadata{
					Title:             "Test Guide",
					MinutesToComplete: 5,
				},
				Steps: []userguides.GuideStep{
					{
						Order:       1,
						Title:       "Step 1",
						Instruction: "Do this",
						Docs: []userguides.GuideDoc{
							{Title: "Docs", URL: "ftp://example.com"},
						},
					},
				},
				Completion: userguides.GuideCompletion{},
			},
			expectErr: true,
			errMsg:    "must use http or https scheme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.guide.Validate()
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error containing %q but got no error", tt.errMsg)
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing %q but got: %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
