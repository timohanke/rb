package blscgo

/*
#cgo CFLAGS:
#cgo LDFLAGS:-lbls -lbls_if -lmcl -lgmp -lgmpxx -lstdc++ -lcrypto
#include "bls_if.h"
*/
import "C"
import "fmt"
import "unsafe"

// Init --
func Init() {
	C.blsInit()
}

// ID -- 
type ID struct {
	v [4]C.uint64_t
}

// getPointer --
func (id *ID) getPointer() (p *C.blsId) {
	return (*C.blsId)(unsafe.Pointer(&id.v[0]))
}

// String --
func (id *ID) String() string {
	buf := make([]byte, 1024)
	n := C.blsIdGetStr(id.getPointer(), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	if n == 0 {
		panic("implementation err. size of buf is small")
	}
	return string(buf[:n])
}

// SetStr --
func (id *ID) SetStr(s string) error {
	buf := []byte(s)
	err := C.blsIdSetStr(id.getPointer(), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	if err > 0 {
		return fmt.Errorf("bad string:%s", s)
	}
	return nil
}

// Set --
func (id *ID) Set(v []uint64) error {
	if len(v) != 4 {
		return fmt.Errorf("bad size (%d), expected size 4", len(v))
	}
	C.blsIdSet(id.getPointer(), (*C.uint64_t)(unsafe.Pointer(&v[0])))
	return nil
}

// SecretKey --
type SecretKey struct {
	v [4]C.uint64_t
}

// getPointer --
func (sec *SecretKey) getPointer() (p *C.blsSecretKey) {
	return (*C.blsSecretKey)(unsafe.Pointer(&sec.v[0]))
}

// String --
func (sec *SecretKey) String() string {
	buf := make([]byte, 1024)
	n := C.blsSecretKeyGetStr(sec.getPointer(), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	if n == 0 {
		panic("implementation err. size of buf is small")
	}
	return string(buf[:n])
}

// SetStr --
func (sec *SecretKey) SetStr(s string) error {
	buf := []byte(s)
	err := C.blsSecretKeySetStr(sec.getPointer(), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	if err > 0 {
		return fmt.Errorf("bad string:%s", s)
	}
	return nil
}

// SetArray --
func (sec *SecretKey) SetArray(v []uint64) error {
	if len(v) != 4 {
		return fmt.Errorf("bad size (%d), expected size 4", len(v))
	}
	C.blsSecretKeySetArray(sec.getPointer(), (*C.uint64_t)(unsafe.Pointer(&v[0])))
	return nil
}

// Init --
func (sec *SecretKey) Init() {
	C.blsSecretKeyInit(sec.getPointer())
}

// Add --
func (sec *SecretKey) Add(rhs *SecretKey) {
	C.blsSecretKeyAdd(sec.getPointer(), rhs.getPointer())
}

// GetMasterSecretKey --
func (sec *SecretKey) GetMasterSecretKey(k int) (msk []SecretKey) {
	msk = make([]SecretKey, k)
	msk[0] = *sec
	for i := 1; i < k; i++ {
		msk[i].Init()
	}
	return msk
}

// GetMasterPublicKey --
func GetMasterPublicKey(msk []SecretKey) (mpk []PublicKey) {
	n := len(msk)
	mpk = make([]PublicKey, n)
	for i := 0; i < n; i++ {
		mpk[i] = *msk[i].GetPublicKey()
	}
	return mpk
}

// Set --
func (sec *SecretKey) Set(msk []SecretKey, id *ID) {
	C.blsSecretKeySet(sec.getPointer(), msk[0].getPointer(), C.size_t(len(msk)), id.getPointer())
}

// Recover --
func (sec *SecretKey) Recover(secVec []SecretKey, idVec []ID) {
	C.blsSecretKeyRecover(sec.getPointer(), secVec[0].getPointer(), idVec[0].getPointer(), C.size_t(len(secVec)))
}

// GetPop --
func (sec *SecretKey) GetPop() (sign *Sign) {
	sign = new(Sign)
	C.blsSecretKeyGetPop(sec.getPointer(), sign.getPointer())
	return sign
}

// PublicKey --
type PublicKey struct {
	v [4 * 2 * 3]C.uint64_t
}

// getPointer --
func (pub *PublicKey) getPointer() (p *C.blsPublicKey) {
	return (*C.blsPublicKey)(unsafe.Pointer(&pub.v[0]))
}

// String --
func (pub *PublicKey) String() string {
	buf := make([]byte, 1024)
	n := C.blsPublicKeyGetStr(pub.getPointer(), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	if n == 0 {
		panic("implementation err. size of buf is small")
	}
	return string(buf[:n])
}

// SetStr --
func (pub *PublicKey) SetStr(s string) error {
	buf := []byte(s)
	err := C.blsPublicKeySetStr(pub.getPointer(), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	if err > 0 {
		return fmt.Errorf("bad string:%s", s)
	}
	return nil
}

// Add --
func (pub *PublicKey) Add(rhs *PublicKey) {
	C.blsPublicKeyAdd(pub.getPointer(), rhs.getPointer())
}
// Set --
func (pub *PublicKey) Set(msk []PublicKey, id *ID) {
	C.blsPublicKeySet(pub.getPointer(), msk[0].getPointer(), C.size_t(len(msk)), id.getPointer())
}

// Recover --
func (pub *PublicKey) Recover(pubVec []PublicKey, idVec []ID) {
	C.blsPublicKeyRecover(pub.getPointer(), pubVec[0].getPointer(), idVec[0].getPointer(), C.size_t(len(pubVec)))
}

// Sign  --
type Sign struct {
	v [4 * 3]C.uint64_t
}

// getPointer --
func (sign *Sign) getPointer() (p *C.blsSign) {
	return (*C.blsSign)(unsafe.Pointer(&sign.v[0]))
}

// String --
func (sign *Sign) String() string {
	buf := make([]byte, 1024)
	n := C.blsSignGetStr(sign.getPointer(), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	if n == 0 {
		panic("implementation err. size of buf is small")
	}
	return string(buf[:n])
}

// SetStr --
func (sign *Sign) SetStr(s string) error {
	buf := []byte(s)
	err := C.blsSignSetStr(sign.getPointer(), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	if err > 0 {
		return fmt.Errorf("bad string:%s", s)
	}
	return nil
}

// GetPublicKey --
func (sec *SecretKey) GetPublicKey() (pub *PublicKey) {
	pub = new(PublicKey)
	C.blsSecretKeyGetPublicKey(sec.getPointer(), pub.getPointer())
	return pub
}

// Sign --
func (sec *SecretKey) Sign(m string) (sign *Sign) {
	sign = new(Sign)
	buf := []byte(m)
	C.blsSecretKeySign(sec.getPointer(), sign.getPointer(), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	return sign
}

// Add --
func (sign *Sign) Add(rhs *Sign) {
	C.blsSignAdd(sign.getPointer(), rhs.getPointer())
}
// Recover --
func (sign *Sign) Recover(signVec []Sign, idVec []ID) {
	C.blsSignRecover(sign.getPointer(), signVec[0].getPointer(), idVec[0].getPointer(), C.size_t(len(signVec)))
}

// Verify --
func (sign *Sign) Verify(pub *PublicKey, m string) bool {
	buf := []byte(m)
	return C.blsSignVerify(sign.getPointer(), pub.getPointer(), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf))) == 1
}

// VerifyPop --
func (sign *Sign) VerifyPop(pub *PublicKey) bool {
	return C.blsSignVerifyPop(sign.getPointer(), pub.getPointer()) == 1
}
