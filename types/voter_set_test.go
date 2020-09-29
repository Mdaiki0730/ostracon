package types

import (
	"bytes"
	"math"
	s "sort"
	"strconv"
	"testing"
	"time"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/rand"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmmath "github.com/tendermint/tendermint/libs/math"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

//-------------------------------------------------------------------

// Check VerifyCommit, VerifyCommitLight and VerifyCommitLightTrusting basic
// verification.
func TestVoterSet_VerifyCommit_All(t *testing.T) {
	var (
		privKey = ed25519.GenPrivKey()
		pubKey  = privKey.PubKey()
		v1      = NewValidator(pubKey, 1000)
		vset    = ToVoterAll([]*Validator{v1})

		chainID = "Lalande21185"
	)

	vote := examplePrecommit()
	vote.ValidatorAddress = pubKey.Address()
	v := vote.ToProto()
	sig, err := privKey.Sign(VoteSignBytes(chainID, v))
	require.NoError(t, err)
	vote.Signature = sig

	commit := NewCommit(vote.Height, vote.Round, vote.BlockID, []CommitSig{vote.CommitSig()})

	vote2 := *vote
	sig2, err := privKey.Sign(VoteSignBytes("EpsilonEridani", v))
	require.NoError(t, err)
	vote2.Signature = sig2

	testCases := []struct {
		description string
		chainID     string
		blockID     BlockID
		height      int64
		commit      *Commit
		expErr      bool
	}{
		{"good", chainID, vote.BlockID, vote.Height, commit, false},

		{"wrong signature (#0)", "EpsilonEridani", vote.BlockID, vote.Height, commit, true},
		{"wrong block ID", chainID, makeBlockIDRandom(), vote.Height, commit, true},
		{"wrong height", chainID, vote.BlockID, vote.Height - 1, commit, true},

		{"wrong set size: 1 vs 0", chainID, vote.BlockID, vote.Height,
			NewCommit(vote.Height, vote.Round, vote.BlockID, []CommitSig{}), true},

		{"wrong set size: 1 vs 2", chainID, vote.BlockID, vote.Height,
			NewCommit(vote.Height, vote.Round, vote.BlockID,
				[]CommitSig{vote.CommitSig(), {BlockIDFlag: BlockIDFlagAbsent}}), true},

		{"insufficient voting power: got 0, needed more than 666", chainID, vote.BlockID, vote.Height,
			NewCommit(vote.Height, vote.Round, vote.BlockID, []CommitSig{{BlockIDFlag: BlockIDFlagAbsent}}), true},

		{"wrong signature (#0)", chainID, vote.BlockID, vote.Height,
			NewCommit(vote.Height, vote.Round, vote.BlockID, []CommitSig{vote2.CommitSig()}), true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			err := vset.VerifyCommit(tc.chainID, tc.blockID, tc.height, tc.commit)
			if tc.expErr {
				if assert.Error(t, err, "VerifyCommit") {
					assert.Contains(t, err.Error(), tc.description, "VerifyCommit")
				}
			} else {
				assert.NoError(t, err, "VerifyCommit")
			}

			err = vset.VerifyCommitLight(tc.chainID, tc.blockID, tc.height, tc.commit)
			if tc.expErr {
				if assert.Error(t, err, "VerifyCommitLight") {
					assert.Contains(t, err.Error(), tc.description, "VerifyCommitLight")
				}
			} else {
				assert.NoError(t, err, "VerifyCommitLight")
			}
		})
	}
}

func TestVoterSet_VerifyCommit_CheckAllSignatures(t *testing.T) {
	var (
		chainID = "test_chain_id"
		h       = int64(3)
		blockID = makeBlockIDRandom()
	)

	voteSet, _, voterSet, vals := randVoteSet(h, 0, tmproto.PrecommitType, 4, 10)
	commit, err := MakeCommit(blockID, h, 0, voteSet, vals, time.Now())
	require.NoError(t, err)

	// malleate 4th signature
	vote := voteSet.GetByIndex(3)
	v := vote.ToProto()
	err = vals[3].SignVote("CentaurusA", v)
	require.NoError(t, err)
	vote.Signature = v.Signature
	commit.Signatures[3] = vote.CommitSig()

	err = voterSet.VerifyCommit(chainID, blockID, h, commit)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "wrong signature (#3)")
	}
}

func TestVoterSet_VerifyCommitLight_ReturnsAsSoonAsMajorityOfVotingPowerSigned(t *testing.T) {
	var (
		chainID = "test_chain_id"
		h       = int64(3)
		blockID = makeBlockIDRandom()
	)

	voteSet, _, voterSet, vals := randVoteSet(h, 0, tmproto.PrecommitType, 4, 10)
	commit, err := MakeCommit(blockID, h, 0, voteSet, vals, time.Now())
	require.NoError(t, err)

	// malleate 4th signature (3 signatures are enough for 2/3+)
	vote := voteSet.GetByIndex(3)
	v := vote.ToProto()
	err = vals[3].SignVote("CentaurusA", v)
	require.NoError(t, err)
	vote.Signature = v.Signature
	commit.Signatures[3] = vote.CommitSig()

	err = voterSet.VerifyCommitLight(chainID, blockID, h, commit)
	assert.NoError(t, err)
}

func TestVoterSet_VerifyCommitLightTrusting_ReturnsAsSoonAsTrustLevelOfVotingPowerSigned(t *testing.T) {
	var (
		chainID = "test_chain_id"
		h       = int64(3)
		blockID = makeBlockIDRandom()
	)

	voteSet, _, voterSet, vals := randVoteSet(h, 0, tmproto.PrecommitType, 4, 10)
	commit, err := MakeCommit(blockID, h, 0, voteSet, vals, time.Now())
	require.NoError(t, err)

	// malleate 3rd signature (2 signatures are enough for 1/3+ trust level)
	vote := voteSet.GetByIndex(2)
	v := vote.ToProto()
	err = vals[2].SignVote("CentaurusA", v)
	require.NoError(t, err)
	vote.Signature = v.Signature
	commit.Signatures[2] = vote.CommitSig()

	err = voterSet.VerifyCommitLightTrusting(chainID, commit, tmmath.Fraction{Numerator: 1, Denominator: 3})
	assert.NoError(t, err)
}

func TestValidatorSet_VerifyCommitLightTrusting(t *testing.T) {
	var (
		blockID                    = makeBlockIDRandom()
		voteSet, _, voterSet, vals = randVoteSet(1, 1, tmproto.PrecommitType, 6, 1)
		commit, err                = MakeCommit(blockID, 1, 1, voteSet, vals, time.Now())
		_, newVoterSet, _          = RandVoterSet(2, 1)
	)
	require.NoError(t, err)

	testCases := []struct {
		voterSet *VoterSet
		err      bool
	}{
		// good
		0: {
			voterSet: voterSet,
			err:      false,
		},
		// bad - no overlap between voter sets
		1: {
			voterSet: newVoterSet,
			err:      true,
		},
		// good - first two are different but the rest of the same -> >1/3
		2: {
			voterSet: WrapValidatorsToVoterSet(append(newVoterSet.Voters, voterSet.Voters...)),
			err:      false,
		},
	}

	for _, tc := range testCases {
		err = tc.voterSet.VerifyCommitLightTrusting("test_chain_id", commit,
			tmmath.Fraction{Numerator: 1, Denominator: 3})
		if tc.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidatorSet_VerifyCommitLightTrustingErrorsOnOverflow(t *testing.T) {
	var (
		blockID                    = makeBlockIDRandom()
		voteSet, _, voterSet, vals = randVoteSet(1, 1, tmproto.PrecommitType, 1, MaxTotalStakingPower)
		commit, err                = MakeCommit(blockID, 1, 1, voteSet, vals, time.Now())
	)
	require.NoError(t, err)

	err = voterSet.VerifyCommitLightTrusting("test_chain_id", commit,
		tmmath.Fraction{Numerator: 25, Denominator: 55})
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "int64 overflow")
	}
}

func countZeroStakingPower(vals []*Validator) int {
	count := 0
	for _, v := range vals {
		if v.StakingPower == 0 {
			count++
		}
	}
	return count
}

func verifyVoterSetSame(t *testing.T, vset1, vset2 *VoterSet) {
	assert.True(t, vset1.Size() == vset2.Size())
	for i, v1 := range vset1.Voters {
		v2 := vset2.Voters[i]
		assert.True(t, v1.Address.String() == v2.Address.String())
		assert.True(t, v1.VotingPower == v2.VotingPower)
		assert.True(t, v1.StakingPower == v2.StakingPower)
	}
}

func verifyVoterSetDifferent(t *testing.T, vset1, vset2 *VoterSet) {
	result := vset1.Size() != vset2.Size()
	if !result {
		for i, v1 := range vset1.Voters {
			v2 := vset2.Voters[i]
			if v1.Address.String() != v2.Address.String() ||
				v1.StakingPower != v2.StakingPower ||
				v1.VotingPower != v2.VotingPower {
				result = true
				break
			}
		}
	}
	assert.True(t, result)
}

func TestSelectVoter(t *testing.T) {
	valSet := randValidatorSet(30)
	valSet.Validators[0].StakingPower = 0

	zeroVals := countZeroStakingPower(valSet.Validators)
	genDoc := &GenesisDoc{
		GenesisTime: tmtime.Now(),
		ChainID:     "tendermint-test",
		VoterParams: &VoterParams{10, 20},
		Validators:  toGenesisValidators(valSet.Validators),
	}
	hash := genDoc.Hash()

	// verifying determinism
	voterSet1 := SelectVoter(valSet, hash, genDoc.VoterParams)
	voterSet2 := SelectVoter(valSet, hash, genDoc.VoterParams)
	verifyVoterSetSame(t, voterSet1, voterSet2)

	// verifying randomness
	hash[0] = (hash[0] & 0xFE) | (^(hash[0] & 0x01) & 0x01) // reverse 1 bit of hash
	voterSet3 := SelectVoter(valSet, hash, genDoc.VoterParams)
	verifyVoterSetDifferent(t, voterSet1, voterSet3)

	// verifying zero-staking removed
	assert.True(t, countZeroStakingPower(voterSet1.Voters) == 0)

	// case that all validators are voters
	voterSet := SelectVoter(valSet, hash, &VoterParams{30, 1})
	assert.True(t, voterSet.Size() == 30-zeroVals)
	voterSet = SelectVoter(valSet, nil, genDoc.VoterParams)
	assert.True(t, voterSet.Size() == 30-zeroVals)
}

func zeroIncluded(valSet *ValidatorSet) bool {
	for _, v := range valSet.Validators {
		if v.StakingPower == 0 {
			return true
		}
	}
	return false
}

func areSame(a *ValidatorSet, b *VoterSet) bool {
	if a.Size() != b.Size() {
		return false
	}
	for i, v := range a.Validators {
		if !v.PubKey.Equals(b.Voters[i].PubKey) {
			return false
		}
		if v.Address.String() != b.Voters[i].Address.String() {
			return false
		}
		if v.StakingPower != b.Voters[i].StakingPower {
			return false
		}
	}
	return true
}

func TestToVoterAll(t *testing.T) {
	valSet := randValidatorSet(30)
	vals := valSet.Validators
	vals[0].StakingPower = 0
	vals[5].StakingPower = 0
	vals[28].StakingPower = 0
	zeroRemovedVoters := ToVoterAll(vals)
	assert.True(t, zeroRemovedVoters.Size() == 27)

	valSet = randValidatorSet(3)
	vals = valSet.Validators
	vals[0].StakingPower = 0
	vals[1].StakingPower = 0
	vals[2].StakingPower = 0
	zeroRemovedVoters = ToVoterAll(vals)
	assert.True(t, zeroRemovedVoters.Size() == 0)

	for i := 0; i < 100; i++ {
		valSet = randValidatorSet(10)
		if zeroIncluded(valSet) {
			continue
		}
		voters := ToVoterAll(valSet.Validators)
		assert.True(t, areSame(valSet, voters), "[%d] %+v != %+v", i, valSet, voters)
	}
}

func toGenesisValidators(vals []*Validator) []GenesisValidator {
	genVals := make([]GenesisValidator, len(vals))
	for i, val := range vals {
		genVals[i] = GenesisValidator{Address: val.Address, PubKey: val.PubKey, Power: val.StakingPower, Name: "name"}
	}
	return genVals
}

/**
The result when we set LoopCount to 10000
  << min power=100, max power=100, actual average voters=10, max voters=10 >> largest gap: 0.040000
  << min power=100, max power=100, actual average voters=20, max voters=20 >> largest gap: 0.030000
  << min power=100, max power=100, actual average voters=29, max voters=29 >> largest gap: 0.010000
  << min power=100, max power=10000, actual average voters=10, max voters=10 >> largest gap: 0.183673
  << min power=100, max power=10000, actual average voters=20, max voters=20 >> largest gap: 0.128788
  << min power=100, max power=10000, actual average voters=28, max voters=29 >> largest gap: 0.304348
  << min power=100, max power=1000000, actual average voters=10, max voters=10 >> largest gap: 0.093158
  << min power=100, max power=1000000, actual average voters=20, max voters=20 >> largest gap: 0.094404
  << min power=100, max power=1000000, actual average voters=28, max voters=29 >> largest gap: 0.194133
  << min power=100, max power=100000000, actual average voters=10, max voters=10 >> largest gap: 0.076536
  << min power=100, max power=100000000, actual average voters=20, max voters=20 >> largest gap: 0.076547
  << min power=100, max power=100000000, actual average voters=29, max voters=29 >> largest gap: 0.147867
*/
func TestSelectVoterReasonableStakingPower(t *testing.T) {
	// Raise LoopCount to get smaller gap over 10000. But large LoopCount takes a lot of time
	const LoopCount = 100
	for minMaxRate := 1; minMaxRate <= 1000000; minMaxRate *= 100 {
		findLargestStakingPowerGap(t, LoopCount, minMaxRate, 10)
		findLargestStakingPowerGap(t, LoopCount, minMaxRate, 20)
		findLargestStakingPowerGap(t, LoopCount, minMaxRate, 29)
	}
}

func findLargestStakingPowerGap(t *testing.T, loopCount int, minMaxRate int, maxVoters int) {
	valSet, privMap := randValidatorSetWithMinMax(30, 100, 100*int64(minMaxRate))
	genDoc := &GenesisDoc{
		GenesisTime: tmtime.Now(),
		ChainID:     "tendermint-test",
		VoterParams: DefaultVoterParams(),
		Validators:  toGenesisValidators(valSet.Validators),
	}
	hash := genDoc.Hash()
	accumulation := make(map[string]int64)
	totalVoters := 0
	for i := 0; i < loopCount; i++ {
		voterSet := SelectVoter(valSet, hash, genDoc.VoterParams)
		for _, voter := range voterSet.Voters {
			accumulation[voter.Address.String()] += voter.StakingPower
		}
		proposer := valSet.SelectProposer(hash, int64(i), 0)
		message := MakeRoundHash(hash, int64(i), 0)
		proof, _ := privMap[proposer.Address.String()].GenerateVRFProof(message)
		pubKey, _ := privMap[proposer.Address.String()].GetPubKey()
		hash, _ = pubKey.VRFVerify(proof, message)
		totalVoters += voterSet.Size()
	}
	largestGap := float64(0)
	for _, val := range valSet.Validators {
		acc := accumulation[val.Address.String()] / int64(loopCount)
		if math.Abs(float64(val.StakingPower-acc))/float64(val.StakingPower) > largestGap {
			largestGap = math.Abs(float64(val.StakingPower-acc)) / float64(val.StakingPower)
		}
	}
	t.Logf("<< min power=100, max power=%d, actual average voters=%d, max voters=%d >> largest gap: %f",
		100*minMaxRate, totalVoters/loopCount, maxVoters, largestGap)
}

/**
  This test is a test to see the difference between MaxVoters and the actual number of elected voters.
  This test is to identify the minimum MaxVoters that cannot be selected as much as MaxVoters by fixing
	MaxSamplingLoopTry.
  If MaxSamplingLoopTry is very large then actual elected voters is up to MaxVoters,
  but large MaxSamplingLoopTry takes too much time.
*/
func TestSelectVoterMaxVarious(t *testing.T) {
	hash := 0
	for minMaxRate := 1; minMaxRate <= 100000000; minMaxRate *= 10000 {
		t.Logf("<<< min: 100, max: %d >>>", 100*minMaxRate)
		for validators := 16; validators <= 256; validators *= 4 {
			for voters := 1; voters <= validators; voters += 10 {
				valSet, _ := randValidatorSetWithMinMax(validators, 100, 100*int64(minMaxRate))
				voterSet := SelectVoter(valSet, []byte{byte(hash)}, &VoterParams{int32(voters), 20})
				if voterSet.Size() < voters {
					t.Logf("Cannot elect voters up to MaxVoters: validators=%d, MaxVoters=%d, actual voters=%d",
						validators, voters, voterSet.Size())
					break
				}
				hash++
			}
		}
	}
}

func TestCalVotersNum(t *testing.T) {
	total := int64(200)
	byzantine := 0.2
	accuracy := 0.99999
	selection := CalNumOfVoterToElect(total, byzantine, accuracy)
	assert.Equal(t, selection, int64(88))

	total = int64(100)
	selection = CalNumOfVoterToElect(total, byzantine, accuracy)
	assert.Equal(t, selection, int64(58))

	assert.Panics(t, func() { CalNumOfVoterToElect(total, 0.3, 10) })
	assert.Panics(t, func() { CalNumOfVoterToElect(total, 1.1, 0.9999) })
}

func TestCalNumOfVoterToElect(t *testing.T) {
	// result of CalNumOfVoterToElect(1, 0.2, 0.99999) ~ CalNumOfVoterToElect(260, 0.2, 0.99999)
	result := []int64{1, 1, 1, 1, 4, 4, 4, 4, 4, 7, 7, 7, 7, 7, 10, 10, 10, 10, 10, 13,
		13, 13, 13, 13, 16, 16, 16, 16, 16, 19, 19, 19, 19, 19, 22, 22, 22, 22, 22, 25,
		25, 25, 25, 25, 28, 28, 28, 28, 28, 31, 31, 31, 31, 31, 34, 34, 34, 34, 34, 37,
		37, 37, 37, 37, 40, 40, 40, 40, 40, 43, 43, 43, 43, 43, 46, 46, 46, 46, 46, 49,
		49, 49, 49, 49, 52, 52, 52, 52, 49, 55, 52, 52, 52, 52, 55, 55, 55, 55, 55, 58,
		58, 58, 58, 58, 61, 61, 58, 58, 58, 61, 61, 61, 61, 61, 64, 64, 64, 64, 61, 67,
		67, 64, 64, 64, 67, 67, 67, 67, 67, 70, 70, 70, 67, 67, 70, 70, 70, 70, 70, 73,
		73, 73, 70, 70, 73, 73, 73, 73, 73, 76, 76, 76, 76, 73, 79, 76, 76, 76, 76, 79,
		79, 79, 76, 76, 79, 79, 79, 79, 79, 82, 82, 82, 79, 79, 82, 82, 82, 82, 82, 85,
		85, 82, 82, 82, 85, 85, 85, 85, 85, 88, 88, 85, 85, 85, 88, 88, 88, 88, 85, 88,
		88, 88, 88, 88, 91, 91, 88, 88, 88, 91, 91, 91, 91, 88, 94, 91, 91, 91, 91, 94,
		94, 94, 91, 91, 94, 94, 94, 94, 94, 97, 94, 94, 94, 94, 97, 97, 97, 94, 94, 97,
		97, 97, 97, 97, 101, 97, 97, 97, 97, 104, 101, 101, 97, 97, 104, 104, 101, 101, 97, 104}

	for i := 1; i <= len(result); i++ {
		assert.True(t, CalNumOfVoterToElect(int64(i), 0.2, 0.99999) == result[i-1])
	}
}

func makeByzantine(valSet *ValidatorSet, rate float64) map[string]bool {
	result := make(map[string]bool)
	byzantinePower := int64(0)
	threshold := int64(float64(valSet.TotalStakingPower()) * rate)
	for _, v := range valSet.Validators {
		if byzantinePower+v.StakingPower > threshold {
			break
		}
		result[v.Address.String()] = true
		byzantinePower += v.StakingPower
	}
	return result
}

func byzantinesPower(voters []*Validator, byzantines map[string]bool) int64 {
	power := int64(0)
	for _, v := range voters {
		if byzantines[v.Address.String()] {
			power += v.VotingPower
		}
	}
	return power
}

func countByzantines(voters []*Validator, byzantines map[string]bool) int {
	count := 0
	for _, v := range voters {
		if byzantines[v.Address.String()] {
			count++
		}
	}
	return count
}

func electVotersForLoop(t *testing.T, hash []byte, valSet *ValidatorSet, privMap map[string]PrivValidator,
	byzantines map[string]bool, loopCount int, byzantinePercent, accuracy int32) {
	byzantineFault := 0
	totalVoters := 0
	totalByzantines := 0
	for i := 0; i < loopCount; i++ {
		voterSet := SelectVoter(valSet, hash, &VoterParams{1, byzantinePercent})
		byzantineThreshold := int64(float64(voterSet.TotalVotingPower())*0.33) + 1
		if byzantinesPower(voterSet.Voters, byzantines) >= byzantineThreshold {
			byzantineFault++
		}
		totalVoters += voterSet.Size()
		totalByzantines += countByzantines(voterSet.Voters, byzantines)
		proposer := valSet.SelectProposer(hash, int64(i), 0)
		message := MakeRoundHash(hash, int64(i), 0)
		proof, _ := privMap[proposer.Address.String()].GenerateVRFProof(message)
		pubKey, _ := privMap[proposer.Address.String()].GetPubKey()
		hash, _ = pubKey.VRFVerify(proof, message)
	}
	t.Logf("voters=%d, fault=%d, avg byzantines=%f",
		totalVoters/loopCount, byzantineFault, float64(totalByzantines)/float64(loopCount))
	assert.True(t, float64(byzantineFault) < float64(loopCount))
}

func TestCalVotersNum2(t *testing.T) {
	valSet, privMap := randValidatorSetWithMinMax(100, 100, 10000)
	byzantinePercent := int32(20)
	byzantines := makeByzantine(valSet, float64(byzantinePercent)/100)
	genDoc := &GenesisDoc{
		GenesisTime: tmtime.Now(),
		ChainID:     "tendermint-test",
		Validators:  toGenesisValidators(valSet.Validators),
	}
	hash := genDoc.Hash()

	loopCount := 1000
	electVotersForLoop(t, hash, valSet, privMap, byzantines, loopCount, byzantinePercent, 1)
	electVotersForLoop(t, hash, valSet, privMap, byzantines, loopCount, byzantinePercent, 2)
	electVotersForLoop(t, hash, valSet, privMap, byzantines, loopCount, byzantinePercent, 3)
	electVotersForLoop(t, hash, valSet, privMap, byzantines, loopCount, byzantinePercent, 4)
	electVotersForLoop(t, hash, valSet, privMap, byzantines, loopCount, byzantinePercent, 5)
}

func TestVoterSetProtoBuf(t *testing.T) {
	_, voterSet, _ := RandVoterSet(10, 100)
	_, voterSet2, _ := RandVoterSet(10, 100)
	voterSet2.Voters[0] = &Validator{}

	testCase := []struct {
		msg      string
		v1       *VoterSet
		expPass1 bool
		expPass2 bool
	}{
		{"success", voterSet, true, true},
		{"fail voterSet2, pubkey empty", voterSet2, false, false},
		{"fail empty voterSet", &VoterSet{}, true, false},
		{"false nil", nil, true, false},
	}
	for _, tc := range testCase {
		protoVoterSet, err := tc.v1.ToProto()
		if tc.expPass1 {
			require.NoError(t, err, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}

		vSet, err := VoterSetFromProto(protoVoterSet)
		if tc.expPass2 {
			require.NoError(t, err, tc.msg)
			require.EqualValues(t, tc.v1, vSet, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}
	}
}

func testVotingPower(t *testing.T, valSet *ValidatorSet) {
	voterParams := &VoterParams{
		VoterElectionThreshold:          100,
		MaxTolerableByzantinePercentage: 20,
	}

	voterSetNoSampling := SelectVoter(valSet, []byte{0}, voterParams)
	for _, v := range voterSetNoSampling.Voters {
		assert.True(t, v.StakingPower == v.VotingPower)
	}

	for i := 90; i > 50; i-- {
		voterParams.VoterElectionThreshold = int32(i)
		voterSetSampling := SelectVoter(valSet, []byte{0}, voterParams)
		allSame := true
		for _, v := range voterSetSampling.Voters {
			if v.StakingPower != v.VotingPower {
				allSame = false
				break
			}
		}
		assert.False(t, allSame)
		assert.True(t, valSet.TotalStakingPower() > voterSetSampling.TotalVotingPower())
		// total voting power can not be less than total staking power - precisionForSelection(1000)

		//TODO: make test code for new voting power
		//assert.True(t, valSet.TotalStakingPower()-voterSetSampling.TotalVotingPower() <= 1000)
	}
}

func TestVotingPower(t *testing.T) {
	testVotingPower(t, randValidatorSet(100))
	vals := make([]*Validator, 100)
	for i := 0; i < len(vals); i++ {
		vals[i] = newValidator(tmrand.Bytes(32), 100)
	}
	testVotingPower(t, NewValidatorSet(vals))
	vals2 := make([]*Validator, 100)
	for i := 0; i < len(vals2); i++ {
		vals2[i] = newValidator(rand.Bytes(32), MaxTotalStakingPower/100)
	}
	testVotingPower(t, NewValidatorSet(vals2))
}

func resetPoints(validators *ValidatorSet) {
	for _, v := range validators.Validators {
		v.VotingPower = 0
	}
}

func isByzantine(voters []voter, totalPriority, tolerableByzantinePercent int64) bool {
	tolerableByzantinePower := totalPriority * tolerableByzantinePercent / 100
	topFVotersVotingPower := countVoters(voters, tolerableByzantinePower)
	return topFVotersVotingPower >= totalPriority/3
}

func TestElectVotersNonDupCandidate(t *testing.T) {
	candidates := newValidatorSet(100, func(i int) int64 { return int64(1000 * (i + 1)) })

	winners := electVotersNonDup(candidates, 0, 20)
	assert.True(t, !isByzantine(winners, candidates.totalStakingPower, 20))
}

// test samplingThreshold
func TestElectVotersNonDupSamplingThreshold(t *testing.T) {
	candidates := newValidatorSet(100, func(i int) int64 { return int64(1000 * (i + 1)) })

	for i := int64(1); i <= 20; i++ {
		winners := electVotersNonDup(candidates, 0, i)
		assert.True(t, !isByzantine(winners, candidates.totalStakingPower, i))
		resetPoints(candidates)
	}
}

// test downscale of win point cases
func TestElectVotersNonDupDownscale(t *testing.T) {
	candidates := newValidatorSet(10, func(i int) int64 {
		if i == 0 {
			return MaxTotalStakingPower >> 1
		}
		if i == 1 {
			return 1 << 55
		}
		if i == 3 {
			return 1 << 54
		}
		if i == 4 {
			return 1 << 55
		}
		return int64(i)
	})
	electVotersNonDup(candidates, 0, 20)
}

// test random election should be deterministic
func TestElectVotersNonDupDeterministic(t *testing.T) {
	candidates1 := newValidatorSet(100, func(i int) int64 { return int64(i + 1) })
	candidates2 := newValidatorSet(100, func(i int) int64 { return int64(i + 1) })
	for i := 1; i <= 100; i++ {
		winners1 := electVotersNonDup(candidates1, uint64(i), 50)
		winners2 := electVotersNonDup(candidates2, uint64(i), 50)
		sameVoters(winners1, winners2)
		resetPoints(candidates1)
		resetPoints(candidates2)
	}
}

func TestElectVotersNonDupIncludingZeroStakingPower(t *testing.T) {
	// first candidate's priority is 0
	candidates1 := newValidatorSet(100, func(i int) int64 { return int64(i) })
	winners1 := electVotersNonDup(candidates1, 0, 20)
	assert.True(t, !isByzantine(winners1, candidates1.totalStakingPower, 20))

	//half of candidates has 0 priority
	candidates2 := newValidatorSet(100, func(i int) int64 {
		if i < 50 {
			return 0
		}
		return int64(i)
	})
	winners2 := electVotersNonDup(candidates2, 0, 20)
	assert.True(t, !isByzantine(winners2, candidates2.totalStakingPower, 20))
}

func TestElectVotersNonDupOverflow(t *testing.T) {
	number := 98
	candidates := newValidatorSet(number, func(i int) int64 { return MaxTotalStakingPower / int64(number+2) })
	totalPriority := candidates.totalStakingPower
	assert.True(t, totalPriority < math.MaxInt64)
	winners := electVotersNonDup(candidates, rand.Uint64(), 20)
	assert.True(t, !isByzantine(winners, totalPriority, 20))
	for _, w := range winners {
		assert.True(t, w.winPoint > 0)
	}
}

//func accumulateAndResetReward(voters []voter, acc []uint64) uint64 {
//	totalWinPoint := uint64(0)
//	for _, v := range voters {
//
//		winPoint := uint64(v.winPoint * float64(precisionForSelection))
//		idx, err := strconv.Atoi(string(v.val.Address.Bytes()))
//		if err != nil {
//			panic(err)
//		}
//		acc[idx] += winPoint
//		totalWinPoint += winPoint
//	}
//	return totalWinPoint
//}

// test reward fairness
//FAILTEST
//func TestElectVotersNonDupReward(t *testing.T) {
//	candidates := newValidatorSet(100, func(i int) uint64 { return uint64(i + 1) })
//
//	accumulatedRewards := make([]uint64, 100)
//	for i := 0; i < 100000; i++ {
//		// 25 samplingThreshold is minimum to pass this test
//		// If samplingThreshold is less than 25, the result says the reward is not fair
//		winners := electVotersNonDup(candidates, uint64(i), 20)
//		accumulateAndResetReward(winners, accumulatedRewards)
//	}
//	for i := 0; i < 99; i++ {
//		assert.True(t, accumulatedRewards[i] < accumulatedRewards[i+1])
//	}
//
//	accumulatedRewards = make([]uint64, 100)
//	for i := 0; i < 50000; i++ {
//		winners := electVotersNonDup(candidates, uint64(i), 20)
//		accumulateAndResetReward(winners, accumulatedRewards)
//	}
//	for i := 0; i < 99; i++ {
//		assert.True(t, accumulatedRewards[i] < accumulatedRewards[i+1])
//	}
//
//	//fail
//	//accumulatedRewards = make([]uint64, 100)
//	//for i := 0; i < 10000; i++ {
//	//	winners := electVotersNonDup(candidates, uint64(i), 20)
//	//	accumulateAndResetReward(winners, accumulatedRewards)
//	//}
//	//for i := 0; i < 99; i++ {
//	//	assert.True(t, accumulatedRewards[i] < accumulatedRewards[i+1])
//	//}
//}

/**
conditions for fair reward
1. even staking power(less difference between min staking and max staking)
2. large total staking(a small total staking power makes a large error when converting float into int)
3. many sampling count
4. loop count
*/
//failed: not fit to new voting
//func TestElectVotersNonDupEquity(t *testing.T) {
//	loopCount := 10000
//
//	// good condition
//	candidates := newValidatorSet(100, func(i int) uint64 { return 1000000 + rand.Uint64()&0xFFFFF })
//	totalStaking := uint64(0)
//	for _, c := range candidates {
//		totalStaking += c.Priority()
//	}
//
//	accumulatedRewards := make([]uint64, 100)
//	totalAccumulateRewards := uint64(0)
//	for i := 0; i < loopCount; i++ {
//		electVotersNonDup(candidates, uint64(i), 20)
//		totalAccumulateRewards += accumulateAndResetReward(candidates, accumulatedRewards)
//	}
//	for i := 0; i < 99; i++ {
//		rewardRate := float64(accumulatedRewards[i]) / float64(totalAccumulateRewards)
//		stakingRate := float64(candidates[i].Priority()) / float64(totalStaking)
//		rate := rewardRate / stakingRate
//		rewardPerStakingDiff := math.Abs(1 - rate)
//		assert.True(t, rewardPerStakingDiff < 0.01)
//	}
//
//	// =======================================================================================================
//	// The codes below are not test codes to verify logic,
//	// but codes to find out what parameters are that weaken the equity of rewards.
//
//	// violation of condition 1
//	candidates = newValidatorSet(100, func(i int) uint64 { return rand.Uint64() & 0xFFFFFFFFF })
//	accumulatedRewards = make([]uint64, 100)
//	for i := 0; i < loopCount; i++ {
//		electVotersNonDup(candidates, uint64(i), 20)
//		accumulateAndResetReward(candidates, accumulatedRewards)
//	}
//	maxRewardPerStakingDiff := float64(0)
//	for i := 0; i < 99; i++ {
//		rewardPerStakingDiff :=
//			math.Abs(float64(accumulatedRewards[i])/float64(candidates[i].Priority())/float64(loopCount) - 1)
//		if maxRewardPerStakingDiff < rewardPerStakingDiff {
//			maxRewardPerStakingDiff = rewardPerStakingDiff
//		}
//	}
//	t.Logf("[! condition 1] max reward per staking difference: %f", maxRewardPerStakingDiff)
//
//	// violation of condition 2
//	candidates = newValidatorSet(100, func(i int) uint64 { return rand.Uint64() & 0xFFFFF })
//	accumulatedRewards = make([]uint64, 100)
//	for i := 0; i < loopCount; i++ {
//		electVotersNonDup(candidates, uint64(i), 20)
//		accumulateAndResetReward(candidates, accumulatedRewards)
//	}
//	maxRewardPerStakingDiff = float64(0)
//	for i := 0; i < 99; i++ {
//		rewardPerStakingDiff :=
//			math.Abs(float64(accumulatedRewards[i])/float64(candidates[i].Priority())/float64(loopCount) - 1)
//		if maxRewardPerStakingDiff < rewardPerStakingDiff {
//			maxRewardPerStakingDiff = rewardPerStakingDiff
//		}
//	}
//	t.Logf("[! condition 2] max reward per staking difference: %f", maxRewardPerStakingDiff)
//
//	// violation of condition 3
//	candidates = newValidatorSet(100, func(i int) uint64 { return 1000000 + rand.Uint64()&0xFFFFF })
//	accumulatedRewards = make([]uint64, 100)
//	for i := 0; i < loopCount; i++ {
//		electVotersNonDup(candidates, uint64(i), 20)
//		accumulateAndResetReward(candidates, accumulatedRewards)
//	}
//	maxRewardPerStakingDiff = float64(0)
//	for i := 0; i < 99; i++ {
//		rewardPerStakingDiff :=
//			math.Abs(float64(accumulatedRewards[i])/float64(candidates[i].Priority())/float64(loopCount) - 1)
//		if maxRewardPerStakingDiff < rewardPerStakingDiff {
//			maxRewardPerStakingDiff = rewardPerStakingDiff
//		}
//	}
//	t.Logf("[! condition 3] max reward per staking difference: %f", maxRewardPerStakingDiff)
//
//	// violation of condition 4
//	loopCount = 100
//	candidates = newValidatorSet(100, func(i int) uint64 { return 1000000 + rand.Uint64()&0xFFFFF })
//	accumulatedRewards = make([]uint64, 100)
//	for i := 0; i < loopCount; i++ {
//		electVotersNonDup(candidates, uint64(i), 99)
//		accumulateAndResetReward(candidates, accumulatedRewards)
//	}
//	maxRewardPerStakingDiff = float64(0)
//	for i := 0; i < 99; i++ {
//		rewardPerStakingDiff :=
//			math.Abs(float64(accumulatedRewards[i])/float64(candidates[i].Priority())/float64(loopCount) - 1)
//		if maxRewardPerStakingDiff < rewardPerStakingDiff {
//			maxRewardPerStakingDiff = rewardPerStakingDiff
//		}
//	}
//	t.Logf("[! condition 4] max reward per staking difference: %f", maxRewardPerStakingDiff)
//}

func newValidatorSet(length int, prio func(int) int64) *ValidatorSet {
	validators := make([]*Validator, length)
	totalStakingPower := int64(0)
	for i := 0; i < length; i++ {
		validators[i] = &Validator{
			Address:      crypto.AddressHash([]byte(strconv.Itoa(i))),
			StakingPower: prio(i),
			VotingPower:  0,
		}
		totalStakingPower += prio(i)
	}

	return &ValidatorSet{
		Validators:        validators,
		totalStakingPower: totalStakingPower,
	}
}

func sameVoters(c1 []voter, c2 []voter) bool {
	if len(c1) != len(c2) {
		return false
	}
	s.Slice(c1, func(i, j int) bool {
		return bytes.Compare(c1[i].val.Address.Bytes(), c1[j].val.Address.Bytes()) == -1
	})
	s.Slice(c2, func(i, j int) bool {
		return bytes.Compare(c2[i].val.Address.Bytes(), c2[j].val.Address.Bytes()) == -1
	})
	for i := 0; i < len(c1); i++ {
		if bytes.Compare(c1[i].val.Address.Bytes(), c2[i].val.Address.Bytes()) == 1 {
			return false
		}
		if c1[i].val.StakingPower != c2[i].val.StakingPower {
			return false
		}
		if c1[i].winPoint != c2[i].winPoint {
			return false
		}
		if c1[i].val.VotingPower != c2[i].val.VotingPower {
			return false
		}
	}
	return true
}

func TestElectVotersNonDup(t *testing.T) {
	candidates := newValidatorSet(5, func(i int) int64 { return 10 })
	expectedPercentage := []float64{
		0.262,
		0.239,
		0.210,
		0.172,
		0.114,
	}
	expectedVotingPower := []int64{
		13,
		11,
		10,
		8,
		5,
	}

	totalWinPoint := float64(0)
	byzantinePercent := int64(10)
	voters := electVotersNonDup(candidates, 0, byzantinePercent)
	assert.True(t, !isByzantine(voters, candidates.totalStakingPower, 10))

	for _, voter := range voters {
		totalWinPoint += voter.winPoint
	}

	for i, voter := range voters {
		assert.True(t, expectedPercentage[i] == float64(int64(voter.winPoint*1000/totalWinPoint))/1000)
		assert.True(t, expectedVotingPower[i] == voter.val.VotingPower)
	}

}

func TestElectVoter(t *testing.T) {
	validators := newValidatorSet(10, func(i int) int64 { return int64(i + 1) })
	total := int64(0)
	for _, val := range validators.Validators {
		total += val.StakingPower
	}
	seed := uint64(0)

	candidates := validators.Validators

	//if fail to voting, panic
	for i := range validators.Validators {
		idx, winner := electVoter(&seed, candidates, i, total)
		total -= winner.val.StakingPower
		moveWinnerToLast(candidates, idx)
	}
}
