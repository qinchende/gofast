// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

func CloseFiles() error {
	if myCnf.LogMedium == logMediumConsole {
		return nil
	}

	if debugLog != nil {
		if err := debugLog.Close(); err != nil {
			return err
		}
	}
	if infoLog != nil {
		if err := infoLog.Close(); err != nil {
			return err
		}
	}
	if warnLog != nil {
		if err := warnLog.Close(); err != nil {
			return err
		}
	}
	if errorLog != nil {
		if err := errorLog.Close(); err != nil {
			return err
		}
	}
	if stackLog != nil {
		if err := stackLog.Close(); err != nil {
			return err
		}
	}
	if statLog != nil {
		if err := statLog.Close(); err != nil {
			return err
		}
	}
	if slowLog != nil {
		if err := slowLog.Close(); err != nil {
			return err
		}
	}

	return nil
}

//func Disable() {
//	once.Do(func() {
//		atomic.StoreUint32(&initialized, 1)
//
//		//infoLog = iox.NopCloser(ioutil.Discard)
//		//errorLog = iox.NopCloser(ioutil.Discard)
//		//severeLog = iox.NopCloser(ioutil.Discard)
//		//slowLog = iox.NopCloser(ioutil.Discard)
//		//statLog = iox.NopCloser(ioutil.Discard)
//		//stackLog = ioutil.Discard
//	})
//}
