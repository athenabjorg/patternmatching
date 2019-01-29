package patternmatcher

// Replace by importing actual eve event structs
type event struct {
	name  string
	value int // For testing predicates and contexts
}

type context struct {
	name          string
	validityCheck func() bool // Probably needs to accept some parameters
}

func (c *context) isValid() bool {
	return c.validityCheck()
}

// A step the pattern has to fulfill before moving onto the next one
type step struct {
	event              event
	predicates         []func(event, []event) bool // Functions on the fly, but all must have the same signiature (current event, previous events)
	contextsToActivate []context
}

type pattern struct {
	nextSteps      []step
	activeContexts []context
	prevEvents     []event
}

func (p *pattern) getNextEvent() event {
	return p.nextSteps[0].event
}

func (p *pattern) eventMatchesCurrentEvent(e event) bool {
	if p.eventTypeMatches(e) && p.predicatesMatch(e) {
		return true
	}
	return false
}

func (p *pattern) eventTypeMatches(e event) bool {
	nextEvent := p.nextSteps[0].event

	if e.name != nextEvent.name {
		return false
	}
	return true
}

func (p *pattern) predicatesMatch(e event) bool {
	predicates := p.nextSteps[0].predicates

	// Predicates are functions on the fly
	for _, predicate := range predicates {
		// Compare predicate with what has happened before
		if !predicate(e, p.prevEvents) {
			return false
		}
	}

	return true
}

func (p *pattern) contextsAreValid(e event) bool {
	contexts := p.activeContexts

	for _, context := range contexts {
		if !context.isValid() {
			return false
		}
	}
	return true
}

func (p *pattern) advance(e event) {

	// TODO remove previous contexts?
	p.activateContexts()
	p.nextSteps = p.nextSteps[1:]
	p.prevEvents = append(p.prevEvents, e)
}

func (p *pattern) activateContexts() {
	contexts := p.nextSteps[0].contextsToActivate
	p.activeContexts = contexts
}

func (p *pattern) isCompleted() bool {
	if len(p.nextSteps) == 0 {
		return true
	}
	return false
}
