package patternmatcher

type patternMatcher struct {
	// They key is the next event type name the pattern is expecting
	templatePatterns map[string][]pattern // Only to be used to create new active patterns
	activePatterns   map[string][]pattern
}

// Receive new events to work with
func (pm *patternMatcher) ProcessEvent(e event) {

	// TODO cancel patterns?

	// Active patterns processed first to that new patterns from templates don't get double processed
	pm.processPatterns(e, pm.activePatterns[e.name], true)
	pm.processPatterns(e, pm.templatePatterns[e.name], false)
}

// Process and update patterns if needed, depending on the event
func (pm *patternMatcher) processPatterns(e event, patterns []pattern, isActive bool) {

	for i := 0; i < len(patterns); i++ {
		if patterns[i].eventMatchesCurrentEvent(e) {

			// Duplicate in case it is a template
			pattern := patterns[i]

			// Cancel pattern if context is not valid
			if !pattern.contextsAreValid(e) {
				if isActive {
					// Remove the pattern from this events mapping
					patterns = append(patterns[:i], patterns[i+1:]...)
					// Back the index up by one to account for the pattern removed
					i--
				}
				// Proceed to next pattern
				continue
			}

			// Advance the pattern, will advance following contexts as well
			pattern.advance(e)

			if pattern.isCompleted() {
				// TODO call handler
			} else {
				// Add the pattern to the list corresponding to it's next event
				nextEvent := pattern.getNextEvent()
				pm.activePatterns[nextEvent.name] = append(pm.activePatterns[nextEvent.name], pattern)
			}

			if isActive {
				// Remove the pattern from this events mapping
				patterns = append(patterns[:i], patterns[i+1:]...)
				// Back the index up by one to account for the pattern removed
				i--
			}
		}
	}
}
