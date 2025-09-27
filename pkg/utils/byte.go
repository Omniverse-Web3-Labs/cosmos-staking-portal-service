package utils

var ByteUtil = newByteUtil()

type byteUtil struct {
}

func newByteUtil() *byteUtil {
	return &byteUtil{}
}

func (*byteUtil) preNum(data byte) int {
	var mask byte = 0x80
	var num int = 0
	//How many 1bits are there before the first 0bit in the 8bit
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}

func (util *byteUtil) ByteIsUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := util.preNum(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//Judge whether the following num - 1 byte starts with 10
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			//Other conditions indicate that it is not utf-8
			return false
		}
	}
	return true
}
