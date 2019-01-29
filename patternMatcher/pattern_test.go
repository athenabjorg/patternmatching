package patternmatcher

import (
	"testing"
)

// Events
var eventOne = event{name: "EventOne"}
var eventTwo = event{name: "EventTwo"}
var eventThree = event{name: "EventThree"}

// Predicates
// Returns true no matter what
var predicateTrue = func(incomingEvent event, prevEvents []event) bool {
	return true
}

// Returns true if the value is the same as that of the first event
var predicateFirstValueMatch = func(incomingEvent event, prevEvents []event) bool {
	if incomingEvent.value == prevEvents[0].value {
		return true
	}
	return false
}

// Steps
var stepOne = step{event: eventOne}
var stepTwo = step{event: eventTwo}
var stepThree = step{event: eventThree}

// Patterns
var patternWithSteps = pattern{
	nextSteps: []step{
		stepOne,
		stepTwo,
		stepThree,
	},
}

var patternWithPredicates = pattern{
	nextSteps: []step{
		step{
			event: eventTwo,
			predicates: []func(event, []event) bool{
				predicateTrue,
				predicateFirstValueMatch,
			},
		},
	},
	prevEvents: []event{
		event{
			name:  "EventOne",
			value: 1,
		},
	},
}

var patternWithSatisfiedContexts = pattern{
	nextSteps: []step{
		step{
			event: eventTwo,
		},
	},
	activeContexts: []context{
		context{
			name: "ContextOne",
			validityCheck: func() bool {
				return true
			},
		},
		context{
			name: "ContextTwo",
			validityCheck: func() bool {
				return true
			},
		},
	},
}

var patternWithUnsatisfiedContexts = pattern{
	nextSteps: []step{
		step{
			event: eventTwo,
		},
	},
	activeContexts: []context{
		context{
			name: "ContextOne",
			validityCheck: func() bool {
				return true
			},
		},
		context{
			name: "ContextTwo",
			validityCheck: func() bool {
				return false
			},
		},
	},
}

var patternWithContextsToActivate = pattern{
	nextSteps: []step{
		step{
			event: eventOne,
			contextsToActivate: []context{
				context{
					name: "ContextOne",
				},
			},
		},
	},
	activeContexts: []context{},
}

func Test_getNextEvent_getsTheCorrectEvent(t *testing.T) {
	p := patternWithSteps

	expected := "EventOne"
	actual := p.getNextEvent().name

	if expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}

func Test_eventMatchesCurrentEvent_eventTypesMatch(t *testing.T) {
	p := patternWithSteps

	expected := true
	actual := p.eventMatchesCurrentEvent(event{name: "EventOne"})

	if expected != actual {
		t.Errorf("Expected '%t', got '%t'", expected, actual)
	}
}

func Test_eventMatchesCurrentEvent_eventsTypesDontMatch(t *testing.T) {

	p := patternWithSteps

	expected := false
	actual := p.eventMatchesCurrentEvent(event{name: "EventTwo"})

	if expected != actual {
		t.Errorf("Expected '%t', got '%t'", expected, actual)
	}
}

func Test_eventMatchesCurrentEvent_eventPredicatesMatch(t *testing.T) {
	p := patternWithPredicates

	expected := true
	actual := p.eventMatchesCurrentEvent(event{name: "EventTwo", value: 1})

	if expected != actual {
		t.Errorf("Expected '%t', got '%t'", expected, actual)
	}
}

func Test_eventMatchesCurrentEvent_eventPredicatesDontMatch(t *testing.T) {
	p := patternWithPredicates

	expected := false
	actual := p.eventMatchesCurrentEvent(event{name: "EventTwo", value: 2})

	if expected != actual {
		t.Errorf("Expected '%t', got '%t'", expected, actual)
	}
}

// func Test_eventMatchesCurrentEvent_eventContextsSatisfied(t *testing.T) {
// 	p := patternWithSatisfiedContexts

// 	expected := true
// 	actual := p.eventMatchesCurrentEvent(event{name: "EventTwo"})

// 	if expected != actual {
// 		t.Errorf("Expected '%t', got '%t'", expected, actual)
// 	}
// }

// func Test_eventMatchesCurrentEvent_eventContextsUnsatisfied(t *testing.T) {
// 	p := patternWithUnsatisfiedContexts

// 	expected := false
// 	actual := p.eventMatchesCurrentEvent(event{name: "EventTwo"})

// 	if expected != actual {
// 		t.Errorf("Expected '%t', got '%t'", expected, actual)
// 	}
// }

func Test_advance_nextEventIsCorrect(t *testing.T) {

	p := patternWithSteps
	p.advance(event{name: "EventOne"})

	expected := "EventTwo"
	actual := p.getNextEvent().name

	if expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}

func Test_advance_prevEventIsCorrect(t *testing.T) {

	p := patternWithSteps
	p.advance(event{name: "EventOne"})
	p.advance(event{name: "EventTwo"})

	expected := "EventTwo"
	actual := p.prevEvents[len(p.prevEvents)-1].name

	if expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}

func Test_advance_contextsAreAdvanced(t *testing.T) {

	p := patternWithContextsToActivate
	p.advance(event{name: "EventOne"})

	expected := "ContextOne"
	actual := p.activeContexts[len(p.prevEvents)-1].name

	if expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}

func Test_isCompleted_yes(t *testing.T) {

	p := patternWithSteps
	p.advance(event{name: "EventOne"})
	p.advance(event{name: "EventTwo"})
	p.advance(event{name: "EventThree"})

	expected := true
	actual := p.isCompleted()

	if expected != actual {
		t.Errorf("Expected '%t', got '%t'", expected, actual)
	}
}

func Test_isCompleted_no(t *testing.T) {

	p := patternWithSteps
	p.advance(event{name: "EventOne"})
	p.advance(event{name: "EventTwo"})

	expected := false
	actual := p.isCompleted()

	if expected != actual {
		t.Errorf("Expected '%t', got '%t'", expected, actual)
	}
}
