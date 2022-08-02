// Code generated by "stringer -type=eventType"; DO NOT EDIT.

package agreement

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[none-0]
	_ = x[votePresent-1]
	_ = x[payloadPresent-2]
	_ = x[bundlePresent-3]
	_ = x[voteVerified-4]
	_ = x[payloadVerified-5]
	_ = x[bundleVerified-6]
	_ = x[roundInterruption-7]
	_ = x[timeout-8]
	_ = x[fastTimeout-9]
	_ = x[speculationTimeout-10]
	_ = x[softThreshold-11]
	_ = x[certThreshold-12]
	_ = x[nextThreshold-13]
	_ = x[proposalCommittable-14]
	_ = x[proposalAccepted-15]
	_ = x[voteFiltered-16]
	_ = x[voteMalformed-17]
	_ = x[bundleFiltered-18]
	_ = x[bundleMalformed-19]
	_ = x[payloadRejected-20]
	_ = x[payloadMalformed-21]
	_ = x[payloadPipelined-22]
	_ = x[payloadAccepted-23]
	_ = x[proposalFrozen-24]
	_ = x[voteAccepted-25]
	_ = x[newRound-26]
	_ = x[newPeriod-27]
	_ = x[readStaging-28]
	_ = x[readPinned-29]
	_ = x[readLowestValue-30]
	_ = x[readLowestPayload-31]
	_ = x[voteFilterRequest-32]
	_ = x[voteFilteredStep-33]
	_ = x[nextThresholdStatusRequest-34]
	_ = x[nextThresholdStatus-35]
	_ = x[freshestBundleRequest-36]
	_ = x[freshestBundle-37]
	_ = x[dumpVotesRequest-38]
	_ = x[dumpVotes-39]
	_ = x[wrappedAction-40]
	_ = x[checkpointReached-41]
}

const _eventType_name = "nonevotePresentpayloadPresentbundlePresentvoteVerifiedpayloadVerifiedbundleVerifiedroundInterruptiontimeoutfastTimeoutspeculationTimeoutsoftThresholdcertThresholdnextThresholdproposalCommittableproposalAcceptedvoteFilteredvoteMalformedbundleFilteredbundleMalformedpayloadRejectedpayloadMalformedpayloadPipelinedpayloadAcceptedproposalFrozenvoteAcceptednewRoundnewPeriodreadStagingreadPinnedreadLowestValuereadLowestPayloadvoteFilterRequestvoteFilteredStepnextThresholdStatusRequestnextThresholdStatusfreshestBundleRequestfreshestBundledumpVotesRequestdumpVoteswrappedActioncheckpointReached"

var _eventType_index = [...]uint16{0, 4, 15, 29, 42, 54, 69, 83, 100, 107, 118, 136, 149, 162, 175, 194, 210, 222, 235, 249, 264, 279, 295, 311, 326, 340, 352, 360, 369, 380, 390, 405, 422, 439, 455, 481, 500, 521, 535, 551, 560, 573, 590}

func (i eventType) String() string {
	if i < 0 || i >= eventType(len(_eventType_index)-1) {
		return "eventType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _eventType_name[_eventType_index[i]:_eventType_index[i+1]]
}
