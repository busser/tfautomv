package pretty

import (
	"fmt"
	"sort"
	"strings"

	"github.com/busser/tfautomv/pkg/engine"
)

// Summarizer is used to pretty print a summary of the tfautomv engine's
// findings.
type Summarizer struct {
	moves       []engine.Move
	comparisons []engine.ResourceComparison

	verbosity int

	// used to explain moves
	modulesWithMoves []string

	// used to explain matches and non-matches
	resourcesToCreateByID  map[string]engine.Resource
	resourcesToDeleteByID  map[string]engine.Resource
	matchCountToCreateByID map[string]int
	matchCountToDeleteByID map[string]int

	// used to build a dynamic legend
	symbolCreateUsed  bool
	symbolDeleteUsed  bool
	symbolIgnoredUsed bool
}

// NewSummarizer returns a ready-to-use Summarizer.
func NewSummarizer(moves []engine.Move, comparisons []engine.ResourceComparison, verbosity int) Summarizer {
	// We precompute some data to simplify the explanation logic.

	var modulesWithMoves []string
	for _, move := range moves {
		modulesWithMoves = append(modulesWithMoves, move.SourceModule, move.DestinationModule)
	}
	modulesWithMoves = unique(modulesWithMoves)
	sort.Strings(modulesWithMoves)

	resourcesToCreateByID := make(map[string]engine.Resource)
	resourcesToDeleteByID := make(map[string]engine.Resource)
	matchCountToCreateByID := make(map[string]int)
	matchCountToDeleteByID := make(map[string]int)
	for _, c := range comparisons {
		resourcesToCreateByID[c.ToCreate.ID()] = c.ToCreate
		resourcesToDeleteByID[c.ToDelete.ID()] = c.ToDelete
		if c.IsMatch() {
			matchCountToCreateByID[c.ToCreate.ID()]++
			matchCountToDeleteByID[c.ToDelete.ID()]++
		}
	}

	return Summarizer{
		moves:       moves,
		comparisons: comparisons,

		verbosity: verbosity,

		modulesWithMoves: modulesWithMoves,

		resourcesToCreateByID: resourcesToCreateByID,
		resourcesToDeleteByID: resourcesToDeleteByID,

		matchCountToCreateByID: matchCountToCreateByID,
		matchCountToDeleteByID: matchCountToDeleteByID,
	}
}

func (s *Summarizer) Summary() string {
	parts := []string{
		Colorf("tfautomv made %s and found %s", StyledNumComparisons(len(s.comparisons)), StyledNumMoves(len(s.moves))),
	}

	var (
		moves          = s.movesFound()
		tooManyMatches = s.tooManyMatches()
		noMatches      = s.noMatches()
		legend         = s.legend()
	)

	if legend != "" {
		parts = append(parts, legend)
	}
	if moves != "" {
		parts = append(parts, moves)
	}
	if tooManyMatches != "" {
		parts = append(parts, tooManyMatches)
	}
	if noMatches != "" {
		parts = append(parts, noMatches)
	}

	return BoxSection("Summary", strings.Join(parts, "\n\n"), "cyan")
}

// legend returns a string explaining the symbols used in the Moves and Matches
// methods.
func (s *Summarizer) legend() string {
	if !s.symbolCreateUsed && !s.symbolDeleteUsed && !s.symbolIgnoredUsed {
		return ""
	}

	lines := []string{"the following symbols are used below:"}

	if s.symbolCreateUsed {
		lines = append(lines, Colorf("  %s the resource Terraform plans to [green][bold]create[reset] has this attribute", s.symbolCreate()))
	}
	if s.symbolDeleteUsed {
		lines = append(lines, Colorf("  %s the resource Terraform plans to [red][bold]delete[reset] has this attribute", s.symbolDelete()))
	}
	if s.symbolIgnoredUsed {
		lines = append(lines, Colorf("  %s differences in this attribute are [yellow][bold]ignored[reset] because of a rule", s.symbolIgnored()))
	}

	return strings.Join(lines, "\n")
}

func (s *Summarizer) movesFound() string {
	return strings.Join(s.moveExplanations(), "\n\n")
}

func (s *Summarizer) tooManyMatches() string {
	return strings.Join(s.multipleMatchExplanations(), "\n\n")
}

func (s *Summarizer) noMatches() string {
	return strings.Join(s.noMatchExplanations(), "\n\n")
}

func (s *Summarizer) symbolCreate() string {
	s.symbolCreateUsed = true
	return Color("[green][bold]+")
}

func (s *Summarizer) symbolDelete() string {
	s.symbolDeleteUsed = true
	return Color("[red][bold]-")
}

func (s *Summarizer) symbolIgnored() string {
	s.symbolIgnoredUsed = true
	return Color("[yellow][bold]~")
}

func (s *Summarizer) styledAddress(addr string) string {
	return Colorf("[bold]%s", addr)
}

func (s *Summarizer) StyledModule(module string) string {
	if module == "." {
		module = "current directory"
	}
	return Colorf("[bold]%s", module)
}

func (s *Summarizer) annotationCreate() string {
	return Color("([green][bold]create[reset])")
}

func (s *Summarizer) annotationDelete() string {
	return Color("([red][bold]delete[reset])")
}

func (s *Summarizer) annotatedResource(r engine.Resource, annotation string) string {
	return Colorf("%s %s in %s", s.styledAddress(r.Address), annotation, s.StyledModule(r.ModuleID))
}

const (
	verbosityListMoves       = 1
	verbosityListComparisons = 1
	verbosityCountAttributes = 2
	verbosityListAttributes  = 3
)

func (s *Summarizer) styledIgnored(comp engine.ResourceComparison) string {
	if len(comp.IgnoredAttributes) == 0 {
		return ""
	}

	if s.verbosity < verbosityCountAttributes {
		return ""
	}

	if s.verbosity < verbosityListAttributes {
		return Colorf("%s  %s", s.symbolIgnored(), s.styledNumAttributes(len(comp.IgnoredAttributes)))
	}

	var lines []string
	for _, attr := range comp.IgnoredAttributes {
		lines = append(lines, Colorf("%s %s", s.symbolIgnored(), attr))
	}

	return strings.Join(lines, "\n")
}

func (s *Summarizer) styledNumAttributes(n int) string {
	if n == 1 {
		return "1 attribute"
	}

	return fmt.Sprintf("%d attributes", n)
}

func (s *Summarizer) styledMismatches(comp engine.ResourceComparison) string {
	if len(comp.MismatchingAttributes) == 0 {
		return ""
	}

	if s.verbosity < verbosityCountAttributes {
		return ""
	}

	if s.verbosity < verbosityListAttributes {
		return Colorf("%s%s %s", s.symbolCreate(), s.symbolDelete(), s.styledNumAttributes(len(comp.MismatchingAttributes)))
	}

	var lines []string
	for _, attr := range comp.MismatchingAttributes {
		lines = append(lines, Colorf("%s %s = %#v", s.symbolCreate(), attr, comp.ToCreate.Attributes[attr]))
		lines = append(lines, Colorf("%s %s = %#v", s.symbolDelete(), attr, comp.ToDelete.Attributes[attr]))
	}

	return strings.Join(lines, "\n")
}

func (s *Summarizer) styledMove(m engine.Move) string {
	comp := s.findComparison(m)

	var lines []string

	lines = append(lines, Colorf("from %s", s.styledAddress(comp.ToDelete.Address)))
	lines = append(lines, Colorf("to   %s", s.styledAddress(comp.ToCreate.Address)))

	ignored := s.styledIgnored(comp)
	if ignored != "" {
		lines = append(lines, "")
		lines = append(lines, ignored)
	}

	return strings.Join(lines, "\n")
}

func StyledNumMoves(n int) string {
	if n == 1 {
		return Color("[bold][green]1 move")
	}

	return Colorf("[bold][green]%d moves", n)
}

func StyledNumComparisons(n int) string {
	if n == 1 {
		return Color("[bold][magenta]1 comparison")
	}

	return Colorf("[bold][magenta]%d comparisons", n)
}

func (s *Summarizer) styledMovesWithinModule(module string) string {
	var styledMoves []string
	for _, move := range s.moves {
		if move.SourceModule == module && move.DestinationModule == module {
			styledMoves = append(styledMoves, s.styledMove(move))
		}
	}

	if len(styledMoves) == 0 {
		return ""
	}

	header := Colorf("%s within %s", StyledNumMoves(len(styledMoves)), s.StyledModule(module))

	if s.verbosity < verbosityListMoves {
		return header
	}

	list := BoxItems(styledMoves, "green")

	return header + "\n" + list
}

func (s *Summarizer) styledMovesBetweenModules(fromModule, toModule string) string {
	var styledMoves []string
	for _, move := range s.moves {
		if move.SourceModule == fromModule && move.DestinationModule == toModule {
			styledMoves = append(styledMoves, s.styledMove(move))
		}
	}

	if len(styledMoves) == 0 {
		return ""
	}

	header := Colorf("%s from %s and %s", StyledNumMoves(len(styledMoves)), s.StyledModule(fromModule), s.StyledModule(toModule))

	if s.verbosity < verbosityListMoves {
		return header
	}

	list := BoxItems(styledMoves, "green")

	return header + "\n" + list
}

func (s *Summarizer) moveExplanations() []string {
	var explanations []string

	for _, module := range s.modulesWithMoves {
		exp := s.styledMovesWithinModule(module)
		if exp != "" {
			explanations = append(explanations, exp)
		}
	}

	for _, fromModule := range s.modulesWithMoves {
		for _, toModule := range s.modulesWithMoves {
			if fromModule == toModule {
				continue
			}

			exp := s.styledMovesBetweenModules(fromModule, toModule)
			if exp != "" {
				explanations = append(explanations, exp)
			}
		}
	}

	return explanations
}

func (s *Summarizer) styledAttributes(c engine.ResourceComparison) string {
	var lines []string

	ignored := s.styledIgnored(c)
	if ignored != "" {
		lines = append(lines, ignored)
	}

	mismatches := s.styledMismatches(c)
	if mismatches != "" {
		lines = append(lines, mismatches)
	}

	return strings.Join(lines, "\n")
}

func StyledNumMatches(n int) string {
	if n == 0 {
		return Color("[bold][red]0 matches")
	}

	if n == 1 {
		return Color("[bold][magenta]1 match")
	}

	return Colorf("[bold][magenta]%d matches", n)
}

func (s *Summarizer) styledMatchesForResourceToCreate(r engine.Resource) string {
	var styledMatches []string
	for _, c := range s.comparisons {
		if c.ToCreate.ID() == r.ID() && c.IsMatch() {
			parts := []string{
				s.annotatedResource(c.ToDelete, s.annotationDelete()),
			}
			styledAttributes := s.styledAttributes(c)
			if styledAttributes != "" {
				parts = append(parts, "", styledAttributes)
			}

			styledMatches = append(styledMatches, strings.Join(parts, "\n"))
		}
	}

	header := Colorf("%s for %s", StyledNumMatches(len(styledMatches)), s.annotatedResource(r, s.annotationCreate()))

	if s.verbosity < verbosityListComparisons {
		return header
	}

	list := BoxItems(styledMatches, "magenta")

	return header + "\n" + list
}

func (s *Summarizer) styledMatchesForResourceToDelete(r engine.Resource) string {
	var styledMatches []string
	for _, c := range s.comparisons {
		if c.ToDelete.ID() == r.ID() && c.IsMatch() {
			parts := []string{
				s.annotatedResource(c.ToCreate, s.annotationCreate()),
			}
			styledAttributes := s.styledAttributes(c)
			if styledAttributes != "" {
				parts = append(parts, "", styledAttributes)
			}

			styledMatches = append(styledMatches, strings.Join(parts, "\n"))
		}
	}

	header := Colorf("%s for %s", StyledNumMatches(len(styledMatches)), s.annotatedResource(r, s.annotationDelete()))

	if s.verbosity < verbosityListComparisons {
		return header
	}

	list := BoxItems(styledMatches, "magenta")

	return header + "\n" + list
}

func (s *Summarizer) multipleMatchExplanations() []string {
	var explanations []string

	for id, toCreate := range s.resourcesToCreateByID {
		if s.matchCountToCreateByID[id] > 1 {
			explanations = append(explanations, s.styledMatchesForResourceToCreate(toCreate))
		}
	}

	for id, toDelete := range s.resourcesToDeleteByID {
		if s.matchCountToDeleteByID[id] > 1 {
			explanations = append(explanations, s.styledMatchesForResourceToDelete(toDelete))
		}
	}

	return explanations
}

func (s *Summarizer) styledNoMatchForResourceToCreate(r engine.Resource) string {
	var styledMatches []string
	for _, c := range s.comparisons {
		if c.ToCreate.ID() == r.ID() && !c.IsMatch() && s.matchCountToDeleteByID[c.ToDelete.ID()] == 0 {
			parts := []string{
				s.annotatedResource(c.ToDelete, s.annotationDelete()),
			}
			styledAttributes := s.styledAttributes(c)
			if styledAttributes != "" {
				parts = append(parts, "", styledAttributes)
			}

			styledMatches = append(styledMatches, strings.Join(parts, "\n"))
		}
	}

	header := Colorf("%s for %s", StyledNumMatches(0), s.annotatedResource(r, s.annotationCreate()))

	if s.verbosity < verbosityListComparisons {
		return header
	}

	list := BoxItems(styledMatches, "red")

	return header + "\n" + list
}

func (s *Summarizer) styledNoMatchForResourceToDelete(r engine.Resource) string {
	var styledMatches []string
	for _, c := range s.comparisons {
		if c.ToDelete.ID() == r.ID() && !c.IsMatch() && s.matchCountToCreateByID[c.ToCreate.ID()] == 0 {
			parts := []string{
				s.annotatedResource(c.ToCreate, s.annotationCreate()),
			}
			styledAttributes := s.styledAttributes(c)
			if styledAttributes != "" {
				parts = append(parts, "", styledAttributes)
			}

			styledMatches = append(styledMatches, strings.Join(parts, "\n"))
		}
	}

	header := Colorf("%s for %s", StyledNumMatches(0), s.annotatedResource(r, s.annotationDelete()))

	if s.verbosity < verbosityListComparisons {
		return header
	}

	list := BoxItems(styledMatches, "red")

	return header + "\n" + list
}

func (s *Summarizer) noMatchExplanations() []string {
	var explanations []string

	for id, toCreate := range s.resourcesToCreateByID {
		if s.matchCountToCreateByID[id] == 0 {
			explanations = append(explanations, s.styledNoMatchForResourceToCreate(toCreate))
		}
	}

	for id, toDelete := range s.resourcesToDeleteByID {
		if s.matchCountToDeleteByID[id] == 0 {
			explanations = append(explanations, s.styledNoMatchForResourceToDelete(toDelete))
		}
	}

	return explanations
}

func (s *Summarizer) findComparison(m engine.Move) engine.ResourceComparison {
	for _, c := range s.comparisons {
		if c.ToCreate.ModuleID == m.DestinationModule &&
			c.ToCreate.Address == m.DestinationAddress &&
			c.ToDelete.ModuleID == m.SourceModule &&
			c.ToDelete.Address == m.SourceAddress {
			return c
		}
	}

	return engine.ResourceComparison{}
}

func unique(s []string) []string {
	seen := make(map[string]struct{})

	var unique []string
	for _, e := range s {
		if _, ok := seen[e]; !ok {
			unique = append(unique, e)
			seen[e] = struct{}{}
		}
	}

	return unique
}
