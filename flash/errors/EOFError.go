package errors

import "errors"

// ErrorEOF 如果尝试读取的内容超出可用数据的末尾, 则会引发 EOFError 异常,
var ErrorEOF = errors.New("EOFError")

// ErrorRange 如果数值不在可接受的范围内, 则会引发 RangeError 异常
var ErrorRange = errors.New("RangeError")
