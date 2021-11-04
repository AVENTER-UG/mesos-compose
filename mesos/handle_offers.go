package mesos

import (
	"github.com/sirupsen/logrus"

	cfg "github.com/AVENTER-UG/mesos-compose/types"

	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
)

func defaultResources(cmd cfg.Command) []mesosproto.Resource {
	CPU := "cpus"
	MEM := "mem"
	cpu := cmd.CPU
	mem := cmd.Memory
	PORT := "ports"

	res := []mesosproto.Resource{
		{
			Name:   CPU,
			Type:   mesosproto.SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: cpu},
		},
		{
			Name:   MEM,
			Type:   mesosproto.SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: mem},
		},
	}

	var portBegin, portEnd uint64

	if cmd.DockerPortMappings != nil {
		portBegin = uint64(cmd.DockerPortMappings[0].HostPort)
		portEnd = portBegin + 2

		res = []mesosproto.Resource{
			{
				Name:   CPU,
				Type:   mesosproto.SCALAR.Enum(),
				Scalar: &mesosproto.Value_Scalar{Value: cpu},
			},
			{
				Name:   MEM,
				Type:   mesosproto.SCALAR.Enum(),
				Scalar: &mesosproto.Value_Scalar{Value: mem},
			},
			{
				Name: PORT,
				Type: mesosproto.RANGES.Enum(),
				Ranges: &mesosproto.Value_Ranges{
					Range: []mesosproto.Value_Range{
						{
							Begin: portBegin,
							End:   portEnd,
						},
					},
				},
			},
		}
	}

	return res
}

// getOffer get out the offer for the mesos task
func getOffer(offers *mesosproto.Event_Offers, cmd cfg.Command) (mesosproto.Offer, []mesosproto.OfferID) {
	offerIds := []mesosproto.OfferID{}

	count := 0
	for n, offer := range offers.Offers {
		logrus.Debug("Got Offer From:", offer.GetHostname())
		offerIds = append(offerIds, offer.ID)

		count = n
	}

	return offers.Offers[count], offerIds

}

// HandleOffers will handle the offers event of mesos
func HandleOffers(offers *mesosproto.Event_Offers) error {
	_, offerIds := getOffer(offers, cfg.Command{})

	select {
	case cmd := <-config.CommandChan:

		takeOffer, offerIds := getOffer(offers, cmd)
		logrus.Debug("Take Offer From:", takeOffer.GetHostname())

		var taskInfo []mesosproto.TaskInfo
		RefuseSeconds := 5.0

		logrus.Debug("Schedule Command: ", cmd.Command)

		taskInfo, _ = prepareTaskInfoExecuteContainer(takeOffer.AgentID, cmd)

		logrus.Debug("HandleOffers cmd: ", taskInfo)

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
		err := Call(accept)
		if err != nil {
			logrus.Error("Handle Offers: ", err)
		}

		// decline unneeded offer
		logrus.Info("Offer Decline: ", offerIds)
		decline := &mesosproto.Call{
			Type:    mesosproto.Call_DECLINE,
			Decline: &mesosproto.Call_Decline{OfferIDs: offerIds},
		}
		return Call(decline)
	default:
		// decline unneeded offer
		logrus.Info("Decline unneeded offer: ", offerIds)
		decline := &mesosproto.Call{
			Type:    mesosproto.Call_DECLINE,
			Decline: &mesosproto.Call_Decline{OfferIDs: offerIds},
		}
		return Call(decline)

	}
}
