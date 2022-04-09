package mesos

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
	"github.com/AVENTER-UG/util"
)

// HandleOffers will handle the offers event of mesos
func HandleOffers(offers *mesosproto.Event_Offers) error {
	select {
	case cmd := <-framework.CommandChan:
		// if no taskid or taskname is given, it's a wrong task.
		if cmd.TaskID == "" || cmd.TaskName == "" {
			return nil
		}

		takeOffer, offerIds := getOffer(offers, cmd)
		if takeOffer.GetHostname() == "" {
			framework.CommandChan <- cmd
			return nil
		}
		logrus.Debug("Take Offer From:", takeOffer.GetHostname())

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
			},
		}

		if getLabelValue("biz.aventer.mesos_compose.executor", cmd) != "" {
			cmd.Executor.Resources = defaultResources(cmd)
			accept.Accept.Operations = []mesosproto.Offer_Operation{{
				Type: mesosproto.Offer_Operation_LAUNCH_GROUP,
				LaunchGroup: &mesosproto.Offer_Operation_LaunchGroup{
					Executor: cmd.Executor,
					TaskGroup: mesosproto.TaskGroupInfo{
						Tasks: taskInfo,
					},
				},
			}}
		} else {
			accept.Accept.Operations = []mesosproto.Offer_Operation{{
				Type: mesosproto.Offer_Operation_LAUNCH,
				Launch: &mesosproto.Offer_Operation_Launch{
					TaskInfos: taskInfo,
				},
			}}
		}

		d, _ := json.Marshal(&accept)
		logrus.Debug("HandleOffers msg: ", util.PrettyJSON(d))

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

// get the value of a label from the command
func getLabelValue(label string, cmd mesosutil.Command) string {
	for _, v := range cmd.Labels {
		if label == v.Key {
			return fmt.Sprint(v.GetValue())
		}
	}
	return ""
}

func getOffer(offers *mesosproto.Event_Offers, cmd mesosutil.Command) (mesosproto.Offer, []mesosproto.OfferID) {
	var offerIds []mesosproto.OfferID
	var offerret mesosproto.Offer

	for n, offer := range offers.Offers {
		logrus.Debug("Got Offer From:", offer.GetHostname())
		offerIds = append(offerIds, offer.ID)

		// if the ressources of this offer does not matched what the command need, the skip
		if !mesosutil.IsRessourceMatched(offer.Resources, cmd) {
			logrus.Debug("Could not found any matched ressources, get next offer")
			mesosutil.Call(mesosutil.DeclineOffer(offerIds))
			continue
		}

		// if contraint_hostname is set, only accept offer with the same hostname
		valHostname := getLabelValue("biz.aventer.mesos_compose.contraint_hostname", cmd)
		if valHostname != "" {
			if strings.ToLower(valHostname) == offer.GetHostname() {
				logrus.Debug("Set Server Constraint to:", offer.GetHostname())
				offerret = offers.Offers[n]
			}
		} else {
			offerret = offers.Offers[n]
		}
	}
	return offerret, offerIds
}
