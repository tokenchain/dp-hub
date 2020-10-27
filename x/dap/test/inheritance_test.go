package test

import (
	"fmt"
	"github.com/streadway/handy/atomic"
	"github.com/tokenchain/dp-hub/x/did/ante"
	"testing"
	"time"
)

type Ball struct {
	//signatures    []dap.IxoSignature
	signatureVals [][]byte
	Radius        int
	Material      string
	total         int
	duer          []byte
}

func NewBallInstance() Ball {
	return Ball{
		Radius:   3,
		Material: "bigass",
		duer:     make([]byte, 100),
		signatureVals: [][]byte{
		},
	}
}
func NewBallStatic() *Ball {
	return &Ball{Radius: 3, Material: "bigass"}
}

func (b *Ball) Bounce() {
	fmt.Printf("Bouncing ball %+v\n", b)
	fmt.Printf("Bouncing ball %s\n", b.Material)
	fmt.Printf("Bouncing ball %d\n", b.Radius)
}
func (b Ball) GetMaterial() string {
	return b.Material
}
func (b Ball) ChangeRed() Ball {
	b.Material = "red"
	return b
}
func (b Ball) newAddress() {
	//b.duer = sdk.AccAddress([]byte("i3o23icapKAPSkoada23213"))
	//copy(b.duer, sdk.AccAddress([]byte("h...i3o23icapKA.........939012930293PSkoada23213")))
	copy(b.duer, []byte("h...i3o23icapKA.........939012930293PSkoada23213"))
}
func (b Ball) swapTz() {
	//var bi []byte
	bi := []byte("ArTCrU37n6nyYvU3HOdd9KfAYEyI8v1gRZNG0NisPIFZd8BLLkmVgiZs77SSBAjYjBgAE+zOi4KycO8BpVP1Dw==")
	/*sample := dap.IxoTx{
		Memo: "this is the sample test",
		Signatures: []dap.IxoSignature{
			types.NewSignature(time.Now(), bi),
			types.NewSignature(time.Now(), bi),
			types.NewSignature(time.Now(), bi),
			types.NewSignature(time.Now(), bi),
		},
	}*/

	/*	b.Memo = "this is the sample test"

		b.Signatures = []dap.IxoSignature{
			types.NewSignature(time.Now(), bi),
			types.NewSignature(time.Now(), bi),
			types.NewSignature(time.Now(), bi),
			types.NewSignature(time.Now(), bi),
		}*/
	s1 := ante.NewSignature(time.Now(), bi)
	s2 := ante.NewSignature(time.Now(), bi)
	s3 := ante.NewSignature(time.Now(), bi)
	//b.signatures = append(b.signatures, s1, s2, s3, )
	b.signatureVals = append(b.signatureVals, s1.SignatureValue, s2.SignatureValue, s3.SignatureValue, )
	//copy(b.x, &sample)
	//b.sameLevelAccess()
	//b.signatureVals = esignatureVals
	//bytes.Join(b.signatureVals, s1.SignatureValue)
	//fmt.Println(len(esignatureVals))
	fmt.Println("===== BALL operation =====")
	fmt.Println(len(b.signatureVals), b.signatureVals)

	b.Material = "blueXXX"
}
func (b Ball) GetSigs() [][]byte {
	return b.signatureVals
}
func (b Ball) ballLevelRun() {
	b.sameLevelAccess()
}
func (b Ball) newAddressRaw() {
	//	b.duer = []byte("i3o23icapKAPSkoada23213")
	newplace := make([]byte, len(b.duer))
	copy(b.duer, []byte("i3o23icapKAPSkoada23213"))
	copy(newplace, b.duer[:len(b.duer)])
	fmt.Println(newplace)
}
func (b Ball) sameLevelAccess() {
	fmt.Println("===== CAN accesss ixoTx =====")
	fmt.Println(b.signatureVals)
}

type Bouncer interface {
	Bounce()
}

func BounceIt(b Bouncer) {
	b.Bounce()
}

type Football struct {
	*Ball
	menem int
}

type Can struct {
	b    Ball
	leng int
}

//var _ Bouncer = &Football{}

func (b Football) GetR() int {
	return b.Radius
}
func (b Football) runnow() {
	b.menem = 6
	b.getlength()
	fmt.Println("===== memem piece ======")
	fmt.Println(fmt.Sprintf("final total is: %d", b.total))
}
func (bv Can) runnow() {
	bv.leng = 6
	bv.getlength()
	/*	fmt.Println("===== length the last piece CAN =====")
		fmt.Println(fmt.Sprintf("final total is: %d", bv.b.total))
		fmt.Println(fmt.Sprintf("material is: %s", bv.b.Material))*/
	bv.b.newAddress()
	//fmt.Println("===== CAN access duer=====")
	//fmt.Println(bv.b.duer)
	bv.b.swapTz()
	fmt.Println("===== children invole save level access=====")
	//bv.b.ballLevelRun()
	fmt.Println(bv.b)
	fmt.Println(bv)
}
func (b Can) getlength() {
	t := b.b.Radius + b.b.Radius + b.b.Radius + b.leng
	b.b.total = t
}
func (b Football) getlength() {
	t := b.Radius + b.Radius + b.Radius + b.menem
	b.total = t
}

func NewBall() *Ball {
	return &Ball{Radius: 3, Material: "bigass"}
}
func (b Football) Bounce(distance int) {
	fmt.Printf("Bouncing football %+v for %d meters\n", b, distance)

	fmt.Printf("Read value %d later\n", b.total)

}

type XCX struct {
	Radius   int
	Material string
	total    atomic.Int
}

func NewBallXInstance() XCX {
	return XCX{
		Radius:   3,
		Material: "bigass",
		total:    atomic.Int(3),
	}
}

func (b XCX) SwapTz() XCX {
	//b.signatureVals = append(b.signatureVals, s1.SignatureValue, s2.SignatureValue, s3.SignatureValue, )
	fmt.Println("===== X operation =====")
	//fmt.Println(len(b.signatureVals), b.signatureVals)
	b.Material = "blueXXX"
	b.Radius = 212
	b.total.Add(300)
	fmt.Println(b.total.String())
	fmt.Printf("List D %+v\n", b)
	fmt.Printf("List color %s\n", b.Material)
	fmt.Printf("List Radius %d\n", b.Radius)
	fmt.Printf("List total %d\n", b.total.Get())
	fmt.Println("===== X operation end =====")
	return b
}
func (b XCX) SecondOp() {
	fmt.Println("===== Sec operation =====")
	fmt.Printf("List D %+v\n", b)
	fmt.Printf("List color %s\n", b.Material)
	fmt.Printf("List Radius %d\n", b.Radius)
	b.total.Add(30)
	fmt.Printf("List total %d\n", b.total.Get())
	fmt.Println("===== Sec operation end =====")
}
func (b XCX) SetMaterial(g string) XCX{
	b.Material = g
	return b
}
func TestMckefe(t *testing.T) {
	/*F := Football{
		Ball:  NewBallStatic(),
		menem: 30,
	}
	C := Can{
		b:    NewBallInstance(),
		leng: 32,
	}

	C.b.Bounce()
	C.b.ChangeRed()


	fmt.Println("===== material is the last piece ======")
	fmt.Printf("C material %s\n", C.b.Material)
	fmt.Printf("the football radius is now - %d\n", F.GetR())
	fmt.Printf("F material %s\n", F.GetMaterial())

	C.b.Bounce()

	C.b.Radius = 10
	C.leng = 310
	C.runnow()*/

	D := NewBallXInstance()
	state2:= D.SwapTz().SwapTz().SwapTz().SwapTz().SetMaterial("white")
	state3:= state2.SwapTz()
	state3.SecondOp()
	fmt.Println("===== ALL operations end =====")
	fmt.Printf("List D %+v\n", &state3)
	fmt.Printf("List color %s\n", state3.Material)
	fmt.Printf("List Radius %d\n", state3.Radius)
	fmt.Printf("List total %d\n", state3.total.Get())
}
