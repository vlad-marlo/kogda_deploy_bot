package telebot

import "time"

func (ctrl *Controller) Poll() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctrl.stopChan:
			ctrl.log.Info("stopped controller polling")
			return
		case <-ticker.C:
			ctrl.HandleTick()
		}
	}
}
