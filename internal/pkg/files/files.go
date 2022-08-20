package files

import (
	"fmt"
	"strconv"
)

func SizeString(value int64) string {
	size := float64(value)
	unit := float64(1000.0)

	switch {
	case size < unit:
		{
			return fmt.Sprintf("%v Bytes", strconv.FormatFloat(size, 'f', 1, 64))
		}
	case size/unit < unit:
		{
			return fmt.Sprintf("%v KB", strconv.FormatFloat(size/unit, 'f', 1, 64))
		}
	case size/unit/unit < unit:
		{
			return fmt.Sprintf("%v MB", strconv.FormatFloat(size/unit/unit, 'f', 1, 64))
		}
	case size/unit/unit/unit < unit:
		{
			return fmt.Sprintf("%v GB", strconv.FormatFloat(size/unit/unit/unit, 'f', 1, 64))
		}
	default:
		{
			return fmt.Sprintf("%v TB", strconv.FormatFloat(size/unit/unit/unit/unit, 'f', 1, 64))
		}
	}
}
