package mesos

import (
	"github.com/sirupsen/logrus"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
)

// HandleOffers will handle the offers event of mesos
func HandleOffers(offers *mesosproto.Event_Offers) error {
	select {
	case cmd := <-framework.CommandChan:

		takeOffer, offerIds := mesosutil.GetOffer(offers, cmd)
		if takeOffer.GetHostname() == "" {
			return nil
		}
		logrus.Debug("Take Offer From:", takeOffer.GetHostname())

		if cmd.TaskID == "" {
			return nil
		}

		var taskInfo []mesosproto.TaskInfo
		RefuseSeconds := 5.0

		taskInfo, _ = PrepareTaskInfoExecuteContainer(takeOffer.AgentID, cmd)

		accept := &mesosproto.Call{
			Type: mesosproto.Call_ACCEPT,
			Accept: &mesosproto.Call_Accept{
				OfferIDs: []mesosproto.OfferID{{
					Value: takeOffer.ID.Value,
				}},
				Filters: &mesosproto.Filters{
					RefuseSeconds: &RefuseSeconds,
				},
				Operations: []mesosproto.Offer_Operation{{
					Type: mesosproto.Offer_Operation_LAUNCH,
					Launch: &mesosproto.Offer_Operation_Launch{
						TaskInfos: taskInfo,
					}}}}}

		logrus.Info("Offer Accept: ", takeOffer.GetID(), " On Node: ", takeOffer.GetHostname())
		err := mesosutil.Call(accept)
		if err != nil {
			logrus.Error("Handle Offers: ", err)
			return err
		}

		// decline unneeded offer
		logrus.Info("Offer Decline: ", offerIds)
		return mesosutil.Call(mesosutil.DeclineOffer(offerIds))
	default:
		// decline unneeded offer
		_, offerIds := mesosutil.GetOffer(offers, mesosutil.Command{})
		logrus.Info("Decline unneeded offer: ", offerIds)
		return mesosutil.Call(mesosutil.DeclineOffer(offerIds))
	}
}
