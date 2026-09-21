package api

import "fmt"

func (p GetDevicesParams) Validate() error {
	if p.Interval != nil && *p.Interval < 0 {
		return fmt.Errorf("cannot specify negative integer: GetDevicesParams.interval")
	}
	return nil
}
