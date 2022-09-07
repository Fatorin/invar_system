package utils

const (
	CodePrefix = "IVT"

	Prime1             = 3
	Prime2             = 5
	ReferrerCodeSalt   = 23114386
	ReferrerCodeLength = 6

	OrderSerialSalt   = 51678234
	OrderSerialLength = 8

	StackRecordSerialSalt   = 92371602
	StackRecordSerialLength = 10
)

var AlphanumericSet = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateReferrerCode(uid uint) string {
	return generateCode(uid, ReferrerCodeSalt, ReferrerCodeLength)
}

func GenerateOrderSerialCode(uid uint) string {
	return generateCode(uid, OrderSerialSalt, OrderSerialLength)
}

func GenerateStackRecordSerialCode(uid uint) string {
	return generateCode(uid, OrderSerialSalt, OrderSerialLength)
}

func generateCode(uid, salt uint, l int) string {

	uid = uid*Prime1 + salt

	var code []rune
	slIdx := make([]byte, l)

	for i := 0; i < l; i++ {
		slIdx[i] = byte(uid % uint(len(AlphanumericSet)))
		slIdx[i] = (slIdx[i] + byte(i)*slIdx[0]) % byte(len(AlphanumericSet))
		uid = uid / uint(len(AlphanumericSet))
	}

	for i := 0; i < l; i++ {
		idx := (byte(i) * byte(Prime2)) % byte(l)
		code = append(code, AlphanumericSet[slIdx[idx]])
	}

	return CodePrefix + string(code)
}
