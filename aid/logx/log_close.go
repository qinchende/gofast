// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

func CloseFiles() error {
	if myCnf.LogMedium == toConsole {
		return nil
	}

	if ioDebug != nil {
		if err := ioDebug.Close(); err != nil {
			return err
		}
	}
	if ioInfo != nil {
		if err := ioInfo.Close(); err != nil {
			return err
		}
	}
	if ioWarn != nil {
		if err := ioWarn.Close(); err != nil {
			return err
		}
	}
	if ioErr != nil {
		if err := ioErr.Close(); err != nil {
			return err
		}
	}
	if ioStack != nil {
		if err := ioStack.Close(); err != nil {
			return err
		}
	}
	if ioStat != nil {
		if err := ioStat.Close(); err != nil {
			return err
		}
	}
	if ioSlow != nil {
		if err := ioSlow.Close(); err != nil {
			return err
		}
	}
	if ioTimer != nil {
		if err := ioTimer.Close(); err != nil {
			return err
		}
	}

	return nil
}

func Disable() {
	initOnce.Do(func() {
		//atomic.StoreUint32(&initialized, 1)

		//ioInfo = iox.NopCloser(ioutil.Discard)
		//ioErr = iox.NopCloser(ioutil.Discard)
		//ioSlow = iox.NopCloser(ioutil.Discard)
		//ioStat = iox.NopCloser(ioutil.Discard)
		//ioStack = ioutil.Discard
	})
}
