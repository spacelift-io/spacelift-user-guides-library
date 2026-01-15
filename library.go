package userguides

import (
	"embed"
	"fmt"
	"io/fs"
	"net/url"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed guides
var guidesFS embed.FS

type Library struct {
	Groups []Group
}

type Group struct {
	Slug        string
	Name        string
	Description string
	SkillLevel  string
	Ordering    int
	Chapters    []Chapter
}

type Chapter struct {
	Slug        string
	Name        string
	Description string
	Ordering    int
	Guides      []Guide
}

type Guide struct {
	Slug       string
	Ordering   int
	Metadata   GuideMetadata
	Steps      []GuideStep
	Completion GuideCompletion
}

type GuideMetadata struct {
	Title             string   `yaml:"title"`
	Description       string   `yaml:"description"`
	Labels            []string `yaml:"labels"`
	Difficulty        string   `yaml:"difficulty"`
	MinutesToComplete int      `yaml:"minutesToComplete"`
}

type GuideStep struct {
	Order       int        `yaml:"order"`
	Title       string     `yaml:"title"`
	Instruction string     `yaml:"instruction"`
	Hint        string     `yaml:"hint"`
	Docs        []GuideDoc `yaml:"docs"`
}

type GuideDoc struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

type GuideCompletion struct {
	SuccessMessage      string   `yaml:"successMessage"`
	RecommendedGuideIDs []string `yaml:"recommendedGuideIds"`
}

func Guides() (*Library, error) {
	lib, err := parse(guidesFS)
	if err != nil {
		panic("userguides: " + err.Error())
	}
	return lib, nil
}

func parse(f fs.FS) (*Library, error) {
	lib := &Library{
		Groups: []Group{},
	}

	groupDirs, err := fs.ReadDir(f, "guides")
	if err != nil {
		return nil, fmt.Errorf("read guides directory: %w", err)
	}

	for _, groupDir := range groupDirs {
		if !groupDir.IsDir() || strings.HasPrefix(groupDir.Name(), ".") {
			continue
		}

		group, err := parseGroup(f, groupDir.Name())
		if err != nil {
			return nil, fmt.Errorf("parse group %s: %w", groupDir.Name(), err)
		}

		lib.Groups = append(lib.Groups, group)
	}

	if err := validateLibrary(lib); err != nil {
		return nil, err
	}

	return lib, nil
}

// validateLibrary performs library-wide validation checks including:
// - Duplicate slug detection across groups, chapters, and guides
// - Referential integrity for recommendedGuideIds
func validateLibrary(lib *Library) error {
	groupSlugs := make(map[string]bool)
	allGuidePaths := make(map[string]bool)

	for _, group := range lib.Groups {
		if groupSlugs[group.Slug] {
			return fmt.Errorf("duplicate group slug: %s", group.Slug)
		}
		groupSlugs[group.Slug] = true

		chapterSlugs := make(map[string]bool)
		for _, chapter := range group.Chapters {
			if chapterSlugs[chapter.Slug] {
				return fmt.Errorf("duplicate chapter slug %s in group %s", chapter.Slug, group.Slug)
			}
			chapterSlugs[chapter.Slug] = true

			guideSlugs := make(map[string]bool)
			for _, guide := range chapter.Guides {
				if guideSlugs[guide.Slug] {
					return fmt.Errorf("duplicate guide slug %s in chapter %s/%s", guide.Slug, group.Slug, chapter.Slug)
				}
				guideSlugs[guide.Slug] = true

				guidePath := group.Slug + "/" + chapter.Slug + "/" + guide.Slug
				allGuidePaths[guidePath] = true
			}
		}
	}

	for _, group := range lib.Groups {
		for _, chapter := range group.Chapters {
			for _, guide := range chapter.Guides {
				for _, recommendedID := range guide.Completion.RecommendedGuideIDs {
					if !allGuidePaths[recommendedID] {
						guidePath := group.Slug + "/" + chapter.Slug + "/" + guide.Slug
						return fmt.Errorf("guide %s references non-existent guide in recommendedGuideIds: %s", guidePath, recommendedID)
					}
				}
			}
		}
	}

	return nil
}

func parseGroup(f fs.FS, groupSlug string) (Group, error) {
	groupPath := path.Join("guides", groupSlug)
	groupYAMLPath := path.Join(groupPath, "group.yaml")

	data, err := fs.ReadFile(f, groupYAMLPath)
	if err != nil {
		return Group{}, fmt.Errorf("read group.yaml: %w", err)
	}

	var groupMeta struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		SkillLevel  string `yaml:"skillLevel"`
		Ordering    int    `yaml:"ordering"`
	}

	if err := yaml.Unmarshal(data, &groupMeta); err != nil {
		return Group{}, fmt.Errorf("parse group.yaml: %w", err)
	}

	group := Group{
		Slug:        groupSlug,
		Name:        groupMeta.Name,
		Description: groupMeta.Description,
		SkillLevel:  groupMeta.SkillLevel,
		Ordering:    groupMeta.Ordering,
		Chapters:    []Chapter{},
	}

	if err := group.Validate(); err != nil {
		return Group{}, err
	}

	chapterDirs, err := fs.ReadDir(f, groupPath)
	if err != nil {
		return Group{}, fmt.Errorf("read group directory: %w", err)
	}

	for _, chapterDir := range chapterDirs {
		if !chapterDir.IsDir() || strings.HasPrefix(chapterDir.Name(), ".") {
			continue
		}

		chapter, err := parseChapter(f, groupSlug, chapterDir.Name())
		if err != nil {
			return Group{}, fmt.Errorf("parse chapter %s: %w", chapterDir.Name(), err)
		}

		group.Chapters = append(group.Chapters, chapter)
	}

	return group, nil
}

func parseChapter(f fs.FS, groupSlug, chapterSlug string) (Chapter, error) {
	chapterPath := path.Join("guides", groupSlug, chapterSlug)
	chapterYAMLPath := path.Join(chapterPath, "chapter.yaml")

	data, err := fs.ReadFile(f, chapterYAMLPath)
	if err != nil {
		return Chapter{}, fmt.Errorf("read chapter.yaml: %w", err)
	}

	var chapterMeta struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Ordering    int    `yaml:"ordering"`
	}

	if err := yaml.Unmarshal(data, &chapterMeta); err != nil {
		return Chapter{}, fmt.Errorf("parse chapter.yaml: %w", err)
	}

	chapter := Chapter{
		Slug:        chapterSlug,
		Name:        chapterMeta.Name,
		Description: chapterMeta.Description,
		Ordering:    chapterMeta.Ordering,
		Guides:      []Guide{},
	}

	if err := chapter.Validate(); err != nil {
		return Chapter{}, err
	}

	entries, err := fs.ReadDir(f, chapterPath)
	if err != nil {
		return Chapter{}, fmt.Errorf("read chapter directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") || entry.Name() == "chapter.yaml" {
			continue
		}

		guide, err := parseGuide(f, groupSlug, chapterSlug, entry.Name())
		if err != nil {
			return Chapter{}, fmt.Errorf("parse guide %s: %w", entry.Name(), err)
		}

		chapter.Guides = append(chapter.Guides, guide)
	}

	return chapter, nil
}

func parseGuide(f fs.FS, groupSlug, chapterSlug, guideFile string) (Guide, error) {
	guidePath := path.Join("guides", groupSlug, chapterSlug, guideFile)

	data, err := fs.ReadFile(f, guidePath)
	if err != nil {
		return Guide{}, fmt.Errorf("read guide file: %w", err)
	}

	var guideMeta struct {
		Ordering   int             `yaml:"ordering"`
		Metadata   GuideMetadata   `yaml:"metadata"`
		Steps      []GuideStep     `yaml:"steps"`
		Completion GuideCompletion `yaml:"completion"`
	}

	if err := yaml.Unmarshal(data, &guideMeta); err != nil {
		return Guide{}, fmt.Errorf("parse guide YAML: %w", err)
	}

	guideSlug := strings.TrimSuffix(guideFile, ".yaml")

	guide := Guide{
		Slug:       guideSlug,
		Ordering:   guideMeta.Ordering,
		Metadata:   guideMeta.Metadata,
		Steps:      guideMeta.Steps,
		Completion: guideMeta.Completion,
	}

	if err := guide.Validate(); err != nil {
		return Guide{}, err
	}

	return guide, nil
}

// Validate performs validation on a Group including:
// - Required fields presence (name, skill level)
// - Skill level enum validation (BEGINNER, ENABLER, COMMANDER, GUARDIAN)
func (g Group) Validate() error {
	if g.Name == "" {
		return fmt.Errorf("group %s: name cannot be empty", g.Slug)
	}
	if g.SkillLevel == "" {
		return fmt.Errorf("group %s: skill level cannot be empty", g.Slug)
	}
	validSkillLevels := map[string]bool{
		"BEGINNER":  true,
		"ENABLER":   true,
		"COMMANDER": true,
		"GUARDIAN":  true,
	}
	if !validSkillLevels[g.SkillLevel] {
		return fmt.Errorf("group %s: invalid skill level %q (must be BEGINNER, ENABLER, COMMANDER, or GUARDIAN)", g.Slug, g.SkillLevel)
	}
	return nil
}

// Validate performs validation on a Chapter including:
// - Required fields presence (name)
func (c Chapter) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("chapter %s: name cannot be empty", c.Slug)
	}
	return nil
}

// Validate performs validation on a Guide including:
// - Required fields presence
// - Step ordering (positive, unique, sequential)
// - URL format validation for documentation links
// - Difficulty enum validation
// - Label validation (non-empty strings)
// - MinutesToComplete validation (non-negative)
func (g Guide) Validate() error {
	if g.Metadata.Title == "" {
		return fmt.Errorf("guide %s: title cannot be empty", g.Slug)
	}
	if len(g.Steps) == 0 {
		return fmt.Errorf("guide %s: must have at least one step", g.Slug)
	}

	if g.Metadata.Difficulty != "" {
		validDifficulties := map[string]bool{
			"easy":   true,
			"medium": true,
			"hard":   true,
		}
		if !validDifficulties[g.Metadata.Difficulty] {
			return fmt.Errorf("guide %s: invalid difficulty %q (must be easy, medium, or hard)", g.Slug, g.Metadata.Difficulty)
		}
	}

	// Validate labels are non-empty
	for i, label := range g.Metadata.Labels {
		if strings.TrimSpace(label) == "" {
			return fmt.Errorf("guide %s: label at index %d is empty", g.Slug, i)
		}
	}

	// Collect and validate step orders
	var orders []int
	stepOrders := make(map[int]bool)
	for _, step := range g.Steps {
		if step.Order <= 0 {
			return fmt.Errorf("guide %s: step order must be positive", g.Slug)
		}
		if step.Title == "" {
			return fmt.Errorf("guide %s: step %d title cannot be empty", g.Slug, step.Order)
		}
		if step.Instruction == "" {
			return fmt.Errorf("guide %s: step %d instruction cannot be empty", g.Slug, step.Order)
		}
		if stepOrders[step.Order] {
			return fmt.Errorf("guide %s: duplicate step order %d", g.Slug, step.Order)
		}
		stepOrders[step.Order] = true
		orders = append(orders, step.Order)

		// Validate documentation URLs
		for _, doc := range step.Docs {
			if doc.Title == "" {
				return fmt.Errorf("guide %s: step %d doc title cannot be empty", g.Slug, step.Order)
			}
			if doc.URL == "" {
				return fmt.Errorf("guide %s: step %d doc URL cannot be empty", g.Slug, step.Order)
			}
			// Validate URL format
			if _, err := url.Parse(doc.URL); err != nil {
				return fmt.Errorf("guide %s: step %d doc URL %q is malformed: %w", g.Slug, step.Order, doc.URL, err)
			}
			// Ensure URL has scheme (http or https)
			parsedURL, _ := url.Parse(doc.URL)
			if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
				return fmt.Errorf("guide %s: step %d doc URL %q must use http or https scheme", g.Slug, step.Order, doc.URL)
			}
		}
	}

	// Validate step ordering is sequential (1, 2, 3, ...)
	sort.Ints(orders)
	for i, order := range orders {
		expectedOrder := i + 1
		if order != expectedOrder {
			return fmt.Errorf("guide %s: steps must be sequentially ordered starting at 1, found order %d at position %d", g.Slug, order, expectedOrder)
		}
	}

	if g.Metadata.MinutesToComplete < 0 {
		return fmt.Errorf("guide %s: minutes to complete cannot be negative", g.Slug)
	}

	return nil
}
