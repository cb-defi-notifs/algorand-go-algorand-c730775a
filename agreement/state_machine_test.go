// Copyright (C) 2019-2024 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package agreement

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/algorand/go-algorand/logging"
)

const truncateIOTrace = false

/*
 * testable.go
 * -----------
 *
 * This file defines a number of interfaces that state machines can implement, so that they
 * are easier to unit test. In particular, implementing these interfaces:
 * - allows us to validate state machines against expected traces
 *
 * We generally model state machines as I/O automata, even though our machines are single-threaded
 * and synchronous. We can generate "traces", which
 * are sequences of input and output actions visible to the outside world. "Traces" do not
 * include state, or internal actions. We validate our state machines against sequences
 * of expected traces (which are actually glorified safety properties).
 * Traces are associated with safety properties, and liveness properties.
 * A trace is not a fragment; it is a history since the instatiation of the automata at its
 * start state.
 *
 * These traces are heavily related to state machine contracts, which are a stateful version
 * of safety property (and liveness property) verification.
 *
 * -----------Rationale for using the model-----------
 * 1. the model defines "parallel composition" of automata A and B.
 *    If A satisfies trace property P_A, and B satisfies trace property P_B, then A*B satisfies P_A*P_B
 *    (which are just the rules applied over the combined trace). This gives an immediate way
 *    to compose trace safety checking when we compose machines together (in particular, making
 *    sure we satisfy the parallel composition criteria: internal actions of A are disjoint with B,
 *    outputs of A and B are disjoint)
 * 2. The model also defines simulation relations: if there is a simulation relation from A to B, then
 *    traces(A) \subseteq \traces(B) (A is lower level, B is a higher level (more abstract) automaton.
 *    (A implements, or simualates, B)
 * 3. Instead of reasoning about inputs and outputs, and introducing an additional router
 *    assumption to pipe outputs back into inputs, we can abstract away the router and reason about the
 *    underlying state machine. This allows us to put post-conditions on event dispatches.
 *
 * -----------
 * Key Ideas:
 *   - test cases (specific traces) are safety properties
 *   - even though I/O automata are asynchronous, we still gain benefits from modeling our (synchronous) state
 *     machines this way. In particular, in the implementation, which is stricter than the model:
 *       - (a) only one action can be enabled at a time (I'm conflating events and actions)
 *       - (b) input events are always spaced out far enough to allow output actions to complete (demux'd)
 *       - (c) state machines can generate only a single action in response to an event
 *       - note that traditional Synchronous State machine models transition as a function of ALL incoming
 *         events. Algorand agreement only delivers one event at a time (due to router mux). This necessitates
 *         part (c) above.
 *       - In particular, this means that traces will look like; Out; In; Out; In; Out (if the mapping from
 *         listeners to IO automata are implemented correctly).
 */

// ioTrace is a collapsed execution trace generated by a state machine.
// It cannot contain nil events. It can handle dispatched event stacks that
// aren't immediately succeeded by an output event (wrapping internal actions)
// if we decide not to hide internal actions (fed through the router) and expose them
// when composing traces of compositions of state machines.
type ioTrace struct {
	events []event // input and output actions
}

func (t *ioTrace) length() int {
	return len(t.events)
}

func (t *ioTrace) extend(eventsToAppend ...event) error {
	if eventsToAppend == nil {
		return fmt.Errorf("Cannot extend trace with nil event")
	}
	t.events = append(t.events, eventsToAppend...)
	return nil
}

func (t *ioTrace) checkWellFormed() error {
	for i := 0; i < len(t.events); i++ {
		if t.events[i] == nil {
			return fmt.Errorf("Trace contains nil event")
		}
	}
	return nil
}

// Primarily for debug purposes
func (t *ioTrace) String() string {
	const TruncLen = 500
	var buf bytes.Buffer
	buf.WriteString("{\n")
	for i := 0; i < len(t.events); i++ {
		buf.WriteString(fmt.Sprintf("\t%v |", t.events[i].ComparableStr()))
		if i%2 == 0 {
			buf.WriteString("\n")
		}
	}
	buf.WriteString("}\n")

	if !truncateIOTrace {
		return buf.String()
	}

	l := buf.Len()
	g := l - TruncLen // try to truncate
	prepend := ""
	if g < 0 {
		g = 0
	} else {
		prepend = "(truncated...)\t"
	}
	return prepend + string(buf.Bytes()[g:])
}

// test helpers
func (t ioTrace) Contains(e event) bool {
	return t.ContainsFn(func(b event) bool {
		return e.ComparableStr() == b.ComparableStr()
	})
}

func (t ioTrace) ContainsString(s string) bool {
	return t.ContainsFn(func(b event) bool {
		return strings.Contains(b.ComparableStr(), s)
	})
}

func (t ioTrace) CountEvent(b event) (count int) {
	for _, e := range t.events {
		if e.ComparableStr() == b.ComparableStr() {
			count++
		}
	}
	return
}

// for each event, passes it into the given fn; if returns true, returns true.
func (t ioTrace) ContainsFn(compareFn func(b event) bool) bool {
	for _, ev := range t.events {
		if compareFn(ev) {
			return true
		}
	}
	return false
}

func (t ioTrace) countAction() (count int) {
	for _, ev := range t.events {
		if ev.t() == wrappedAction {
			count++
		}
	}
	return
}

// ioSafetyProp denotes whether some trace is "safe" according to itself
type ioSafetyProp interface {
	// returns bool whether trace is in the safety property. If false,
	// optionally accompanied by an informational message pertaining to why
	// the trace is not in the safety property. Err is set if we saw an
	// unforeseen error.
	containsTrace(trace ioTrace) (contains bool, info string, err error)

	// every safety prop also exposes the option to check it dynamically
	newPropChecker() ioPropChecker
}

// the safety prop that contains all traces (can be used as a stub)
type ioPropAll struct {
}

func (s ioPropAll) containsTrace(trace ioTrace) (bool, string, error) {
	return true, "", nil
}

func (s ioPropAll) newPropChecker() ioPropChecker {
	return new(ioPropAllChecker)
}

type ioPropAllChecker struct {
}

func (c ioPropAllChecker) addEvent(e event) error {
	return nil
}

// ioPropChecker is a stateful safety prop validator
type ioPropChecker interface {
	// add another event in order from the trace, returns error if the
	// addition of this event excludes the source trace from the safety prop
	addEvent(e event) error
}

type ioPropCheckerFactory interface {
	newPropChecker() ioPropChecker
}

// ioPropWrapper implements ioSafetyProp, wrapping an iopropChecker
type ioPropWrapper struct {
	ioPropCheckerFactory
}

func (w ioPropWrapper) containsTrace(trace ioTrace) (contains bool, info string, err error) {
	checker := w.newPropChecker()
	err = trace.checkWellFormed()
	if err != nil {
		return false, "", err
	}
	for _, e := range trace.events {
		valid := checker.addEvent(e)
		if valid != nil {
			return false, valid.Error(), nil
		}
	}
	return true, "", nil
}

// directMatchIoSafetyProp is a safety prop that returns "safe" if a trace
// contains the specified test trace as a prefix, or matches a prefix of the direct match
type directMatchIoSafetyProp struct {
	directMatchTrace ioTrace
}

// containsTrace validates traces if and only if they match our expected actions
func (e *directMatchIoSafetyProp) containsTrace(trace ioTrace) (bool, string, error) {
	for i := 0; i < trace.length(); i++ {
		if i >= e.directMatchTrace.length() {
			return true, "", nil
		}
		// compare using String(), in case event is uncomparable
		// This comparison is very loose, but we only need to match event types for now
		// Only exception: an error field type matches anything of same type
		if trace.events[i].ComparableStr() != e.directMatchTrace.events[i].ComparableStr() {
			return false, "", nil
		}
	}
	return true, "", nil
}

func (e *directMatchIoSafetyProp) newPropChecker() ioPropChecker {
	panic("Unsupported; direct match safety prop cannot dynamically check traces (yet)")
}

// ioAutomata is a traceable state machine. The trace hides internal actions.
// Why is this useful when listener is already checked? We can impose
// test-only safety properties, for instance. This is, in fact, how input/output
// unit test event matching works - it's expressed as a safety property.
// Eventually, checkedListener should implement ioAutomata! For now, we wrap
// checked listeners in ioAutomataConcrete to implement ioAutomata.
type ioAutomata interface {
	ioTraceable

	// #todo eventually, post refactor, these methods should be equivalent
	// to the listener interface. However, listener currently requires router
	// and player, which ioAutomata should fully encapsulate. Once we remove the
	// ioAutomata dependency on player (and wrap the router) the interfaces
	// can probably be combined.
	transition(input event) (err error, panicError error)
	transitionAll(inputs []event) (err error, panicError error)
}

type ioTraceable interface {
	// getTrace returns the trace of the execution so far.
	// note that it always starts at initialization, at the start state.
	getTrace() ioTrace
	// getTraceVisible returns a trace without hiding internal events.
	getTraceVisible() ioTrace
	// resetTrace resets the stored trace, erasing history, in effect
	// restarting it at the point when resetTrace is called
	resetTrace()
}

// ioAutomataConcrete is a concrete wrapper around listener, implementing ioAutomata.
// ioAutomataConcrete also implements router to help with assembling traces.
// Wraps listeners with a router and player object. Also catches panics and wraps them in an error.
type ioAutomataConcrete struct {
	listener

	// listeners need additional context. For now, we keep it static.
	routerCtx router // optional, set to router{} on default
	playerCtx player // optional, set to player{} by default

	// private
	savedHiddenTrace ioTrace // hides internal events, output of getTrace
	savedTrace       ioTrace
	rHandle          *routerHandle
}

func (w *ioAutomataConcrete) getTrace() ioTrace {
	return w.savedHiddenTrace
}

// resets the trace, for instance to make testing for recent events easier
func (w *ioAutomataConcrete) resetTrace() {
	w.savedHiddenTrace = ioTrace{}
	w.savedTrace = ioTrace{}
}

func (w *ioAutomataConcrete) getTraceVisible() ioTrace {
	return w.savedTrace
}

// Hijack router so that we can track internal events dispatched between state machines.
// Alternatively, we create a tracer interface and pass ourselves in
// as the tracer - but hijacking router seems to be less impactful since an interface
// already exists.
func (w *ioAutomataConcrete) dispatch(t *tracer, state player, e event, src stateMachineTag, dest stateMachineTag, r round, p period, s step) event {
	_ = w.savedTrace.extend(e)
	out := w.routerCtx.dispatch(t, state, e, src, dest, r, p, s)
	_ = w.savedTrace.extend(out)
	return out
}

func (w *ioAutomataConcrete) callHandler(inputTraceEvent event) (outEvent event, panicErr error) {
	logging.Base().SetOutput(nullWriter{})
	defer func() {
		logging.Base().SetOutput(os.Stderr)
		r := recover()
		if r != nil {
			panicErr = fmt.Errorf("%v", r)
		}
	}()
	if w.rHandle == nil {
		w.rHandle = &routerHandle{t: &tracer{log: serviceLogger{logging.Base()}}, r: w}
	}
	outEvent = w.listener.handle(*w.rHandle, w.playerCtx, inputTraceEvent)
	return
}

func (w *ioAutomataConcrete) transition(inputTraceEvent event) (err error, panicErr error) {
	out, callPanicErr := w.callHandler(inputTraceEvent)
	if callPanicErr != nil {
		// the first err will be more useful once state machines propagate errors upwards
		return err, callPanicErr
	}

	// extend saved traces
	err = w.savedHiddenTrace.extend(inputTraceEvent)
	if err != nil {
		return err, nil
	}
	err = w.savedHiddenTrace.extend(out)
	if err != nil {
		return err, nil
	}
	err = w.savedTrace.extend(inputTraceEvent)
	if err != nil {
		return err, nil
	}
	err = w.savedTrace.extend(out)
	if err != nil {
		return err, nil
	}

	return nil, nil
}

func (w *ioAutomataConcrete) transitionAll(inputs []event) (error, error) {
	for i := 0; i < len(inputs); i++ {
		err, panicErr := w.transition(inputs[i]) // a nil event is interpreted as no input
		if err != nil || panicErr != nil {
			return err, panicErr
		}
	}
	return nil, nil
}

/* Testing Utils */

type blackhole struct{}

func (blackhole) Write(data []byte) (int, error) {
	return len(data), nil
}

// deterministicTraceTestCase encapsulates a traditional unit test test case.
type determisticTraceTestCase struct {
	inputs          []event
	expectedOutputs []event
	safetyProps     []ioSafetyProp
}

// Validate takes a given automata at zero state, drives it with the test case input,
// and validates the output.
func (testCase *determisticTraceTestCase) Validate(automaton ioAutomata) (invalidErr error, runtimeErr error) {
	return testCase.ValidateAsExtension(automaton)
}

// ValidateAsExtension takes a given automata that is already in some state, drives it
// with some addition input (an "extension"), and validates the output.
func (testCase *determisticTraceTestCase) ValidateAsExtension(automaton ioAutomata) (invalidErr error, runtimeErr error) {
	// suppress error logging from contract-checkers
	logging.Base().SetOutput(blackhole{})
	defer func() {
		logging.Base().SetOutput(os.Stderr)
	}()

	if len(testCase.inputs) != len(testCase.expectedOutputs) && len(testCase.inputs) != len(testCase.expectedOutputs)+1 {
		return nil, fmt.Errorf("Malformed test case: either inputs and outputs must be same length, or inputs should be one longer than outputs")
	}

	// the automata may have already run some and generated a trace
	existingTraceLength := len(automaton.getTrace().events)

	// construct partial input and final expected extension traces
	allEvents := make([]event, len(testCase.expectedOutputs)*2)
	for i := 0; i < len(testCase.inputs); i++ {
		if i < len(testCase.expectedOutputs) {
			allEvents[2*i] = testCase.inputs[i]
			allEvents[2*i+1] = testCase.expectedOutputs[i]
		}
	}
	expectedFinalTrace := ioTrace{allEvents}
	err := expectedFinalTrace.checkWellFormed()
	if err != nil {
		return nil, fmt.Errorf("Outputs cannot contain nil events; %v", err)
	}
	err, panicErr := automaton.transitionAll(testCase.inputs)
	if err != nil {
		return nil, err
	}
	outputTrace := automaton.getTrace()
	outputTraceLen := outputTrace.length()
	outputTraceExtension := ioTrace{outputTrace.events[existingTraceLength:]}
	validator := directMatchIoSafetyProp{expectedFinalTrace}
	traceValid, _, runtimeErr := validator.containsTrace(outputTraceExtension)
	if runtimeErr != nil {
		return nil, runtimeErr
	}

	// any trace should be valid up to the point of panicking
	if !traceValid {
		invalidErr = errIOTraceDiverge{expected: expectedFinalTrace.String(), actual: outputTraceExtension.String()}
		return invalidErr, nil
	}

	if len(testCase.inputs) == len(testCase.expectedOutputs) {
		// we have one output for each input if and only if we did not ever panic
		if panicErr != nil {
			invalidErr = fmt.Errorf("Panicked when we were not expecting it: %v", panicErr)
			return invalidErr, nil
		}
	} else if len(testCase.inputs) == len(testCase.expectedOutputs)+1 {
		// we have a dangling final output if and only if we panicked
		if panicErr == nil {
			invalidErr = fmt.Errorf("Did not panic when we were expecting to")
			return invalidErr, nil
		}
	} else {
		invalidErr = fmt.Errorf("Input size (%v) is inconsistent with output size (%v)", testCase.inputs, testCase.expectedOutputs)
		return invalidErr, nil
	}

	if outputTraceLen < expectedFinalTrace.length() {
		if panicErr != nil {
			invalidErr = fmt.Errorf("Panicked early: %d:%d:\t%v",
				outputTraceLen, expectedFinalTrace.length(), panicErr)
			return invalidErr, nil
		}
		// since we are validating a synchronous state machine, the output trace should be as long as expected
		invalidErr = fmt.Errorf("Trace too short (%d:%d). %v %v",
			outputTraceLen, expectedFinalTrace.length(), expectedFinalTrace, outputTraceExtension)
		return invalidErr, nil
	}

	// finally, validate (the entire) output trace (not just the extension) against specified safety properties, if any
	for _, sp := range testCase.safetyProps {
		good, msg, err := sp.containsTrace(outputTrace)
		if err != nil {
			return nil, fmt.Errorf("Error evaluating safety property %v", sp)
		}
		if !good {
			return fmt.Errorf("Trace not in safety property: %v, %s", sp, msg), nil
		}
	}

	return nil, nil
}

// a convenience helper
type testCaseBuilder struct {
	inputs          []event
	expectedOutputs []event
	safetyProps     []ioSafetyProp
}

func (b *testCaseBuilder) Build() *determisticTraceTestCase {
	return &determisticTraceTestCase{b.inputs, b.expectedOutputs, b.safetyProps}
}

func (b *testCaseBuilder) AddInOutPair(input event, output event) {
	if b.inputs == nil {
		b.inputs = make([]event, 0, 1)
	}
	if b.expectedOutputs == nil {
		b.expectedOutputs = make([]event, 0, 1)
	}
	b.inputs = append(b.inputs, input)
	b.expectedOutputs = append(b.expectedOutputs, output)
}

func (b *testCaseBuilder) AddSafetyProp(prop ioSafetyProp) {
	if b.safetyProps == nil {
		b.safetyProps = make([]ioSafetyProp, 0, 1)
	}
	b.safetyProps = append(b.safetyProps, prop)
}

type errIOTraceDiverge struct {
	expected string
	actual   string
}

func (err errIOTraceDiverge) Error() string {
	return fmt.Sprintf("Expected: %s, Actual %s", err.expected, err.actual)
}

/* Utils for player testing */

// wrap actions as events so we can test player as a listener
func ev(a action) event {
	return wrappedActionEvent{a}
}

type wrappedActionEvent struct {
	action
}

func (e wrappedActionEvent) t() eventType {
	return wrappedAction
}

func (e wrappedActionEvent) String() string {
	return e.action.String()
}

func (e wrappedActionEvent) ComparableStr() string {
	return e.action.ComparableStr()
}

// ioAutomataConcretePlayer is a concrete wrapper around root router, implementing ioAutomata.
type ioAutomataConcretePlayer struct {
	*rootRouter

	savedTrace *ioTrace

	// need to stub out these objects
	t *tracer
}

func (w *ioAutomataConcretePlayer) getTrace() ioTrace {
	return *w.savedTrace
}

// resets the trace, for instance to make testing for recent events easier
func (w *ioAutomataConcretePlayer) resetTrace() {
	w.savedTrace = nil
}

func (w *ioAutomataConcretePlayer) getTraceVisible() ioTrace {
	panic("unsupported")
}

func (w *ioAutomataConcretePlayer) underlying() *player {
	return w.rootRouter.root.underlying().(*player)
}

func (w *ioAutomataConcretePlayer) callSubmitTop(inputTraceEvent event) (outEvents []event, panicErr error) {
	logging.Base().SetOutput(nullWriter{})
	defer func() {
		logging.Base().SetOutput(os.Stderr)
		r := recover()
		if r != nil {
			panicErr = fmt.Errorf("Panic: %v", r)
		}
	}()
	_, actions := w.rootRouter.submitTop(w.t, *w.underlying(), inputTraceEvent)
	// wrap all actions as events
	outEvents = make([]event, len(actions))
	for i, a := range actions {
		outEvents[i] = ev(a)
	}
	return
}

func (w *ioAutomataConcretePlayer) transition(inputTraceEvent event) (err error, panicErr error) {
	if w.savedTrace == nil {
		w.savedTrace = &ioTrace{}
	}
	if w.t == nil {
		w.t = &tracer{log: serviceLogger{logging.Base()}}
	}
	outEvents, callPanicErr := w.callSubmitTop(inputTraceEvent)
	if callPanicErr != nil {
		// the first err will be more useful once state machines propagate errors upwards
		return err, callPanicErr
	}

	// extend saved trace
	err = w.savedTrace.extend(inputTraceEvent)
	if err != nil {
		return err, nil
	}
	err = w.savedTrace.extend(outEvents...)
	if err != nil {
		return err, nil
	}

	return nil, nil
}

func (w *ioAutomataConcretePlayer) transitionAll(inputs []event) (error, error) {
	for i := 0; i < len(inputs); i++ {
		err, panicErr := w.transition(inputs[i]) // a nil event is interpreted as no input
		if err != nil || panicErr != nil {
			return err, panicErr
		}
	}
	return nil, nil
}
