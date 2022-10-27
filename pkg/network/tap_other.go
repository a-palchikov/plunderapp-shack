//go:build !linux

package network

import "errors"

func (*Environment) CreateTap(tapName string) error {
	return errors.New("not implemented")
}

func (*Environment) DeleteTap(tapName string) error {
	return errors.New("not implemented")
}
