package bls

// types
type Pop Signature

// Proof-of-Possesion
func GeneratePop(sec Seckey, pub Pubkey) Pop {
    return Pop(Sign(sec, []byte(pub.String())))
}

// Verification
func VerifyPop(pub Pubkey, pop Pop) bool {
  return VerifySig(pub, []byte(pub.String()), Signature(pop))
}

