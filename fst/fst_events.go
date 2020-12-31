package fst

const (
	EReady = "onReady"
	ERoute = "onRoute"
	EClose = "onClose"
)

type appEvents struct {
	eReadyHds FHandlers
	eRouteHds FHandlers
	eCloseHds FHandlers
}

func (ft *Faster) On(eType string, handles ...FHandler) {
	switch eType {
	case EReady:
		ft.eReadyHds = append(ft.eReadyHds, handles...)
	case ERoute:
		ft.eRouteHds = append(ft.eRouteHds, handles...)
	case EClose:
		ft.eCloseHds = append(ft.eCloseHds, handles...)
	default:
		panic("Server event type error, can't find this type.")
	}
}

func (ft *Faster) execHandlers(hds FHandlers) {
	for i, hLen := 0, len(hds); i < hLen; i++ {
		hds[i](ft)
	}
	//for _, next := range hds {
	//	next(ft)
	//}
}
