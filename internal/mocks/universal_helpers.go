package mocks

import "github.com/stretchr/testify/mock"

func GetMockResult[T any](args mock.Arguments) (T, error) {
	if result, ok := args.Get(0).(T); ok {
		return result, args.Error(1)
	}
	var zeroValue T
	return zeroValue, args.Error(1)
}
