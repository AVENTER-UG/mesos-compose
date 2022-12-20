package scheduler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/AVENTER-UG/util/util"
)

// HandleOffers will handle the offers event of mesos
func (e *Scheduler) HandleOffers(offers *mesosproto.Event_Offers) error {
	select {
	case cmd := <-e.Framework.CommandChan:
		// if no taskid or taskname is given, it's a wrong task.
		if cmd.TaskID == "" || cmd.TaskName == "" {
			return nil
		}

		takeOffer, offerIds := e.getOffer(offers, cmd)
		if takeOffer.GetHostname() == "" {
			cmd.State = ""
			e.Redis.SaveTaskRedis(cmd)
			logrus.WithField("func", "mesos.HandleOffers").Debug("No matched offer found.")
			logrus.WithField("func", "mesos.HandleOffers").Debug("Decline unneeded offer: ", offerIds)
			return e.Mesos.Call(e.Mesos.DeclineOffer(offerIds))
		}
		logrus.Debug("Take Offer From:", takeOffer.GetHostname())

		var taskInfo []mesosproto.TaskInfo
		RefuseSeconds := 5.0

		taskInfo, _ = e.PrepareTaskInfoExecuteContainer(takeOffer.AgentID, cmd)

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

		accept.Accept.Operations = []mesosproto.Offer_Operation{{
			Type: mesosproto.Offer_Operation_LAUNCH,
			Launch: &mesosproto.Offer_Operation_Launch{
				TaskInfos: taskInfo,
			},
		}}

		d, _ := json.Marshal(&accept)
		logrus.Debug("HandleOffers msg: ", util.PrettyJSON(d))

		logrus.Info("Offer Accept: ", takeOffer.GetID(), " On Node: ", takeOffer.GetHostname())
		err := e.Mesos.Call(accept)
		if err != nil {
			logrus.Error("Handle Offers: ", err)
			return err
		}

		// decline unneeded offer
		logrus.Info("Offer Decline: ", offerIds)
		return e.Mesos.Call(e.Mesos.DeclineOffer(offerIds))
	default:
		// decline unneeded offer
		_, offerIds := e.Mesos.GetOffer(offers, cfg.Command{})
		logrus.Info("Decline unneeded offer: ", offerIds)
		return e.Mesos.Call(e.Mesos.DeclineOffer(offerIds))
	}
}

// get the value of a label from the command
func (e *Scheduler) getLabelValue(label string, cmd cfg.Command) string {
	for _, v := range cmd.Labels {
		if label == v.Key {
			return fmt.Sprint(v.GetValue())
		}
	}
	return ""
}

func (e *Scheduler) getOffer(offers *mesosproto.Event_Offers, cmd cfg.Command) (mesosproto.Offer, []mesosproto.OfferID) {
	var offerIds []mesosproto.OfferID
	var offerret mesosproto.Offer

	for _, offer := range offers.Offers {
		logrus.Debug("Got Offer From:", offer.GetHostname())
		offerIds = append(offerIds, offer.ID)

		// if the ressources of this offer does not matched what the command need, the skip
		if !e.Mesos.IsRessourceMatched(offer.Resources, cmd) {
			logrus.Debug("Could not found any matched ressources, get next offer")
			e.Mesos.Call(e.Mesos.DeclineOffer(offerIds))
			continue
		}

		// if contraint_hostname is set, only accept offer with the same hostname
		valHostname := e.getLabelValue("__mc_placement_node_hostname", cmd)
		if valHostname != "" {
			if strings.ToLower(valHostname) == offer.GetHostname() {
				logrus.Debug("Set Server Hostname Constraint to:", offer.GetHostname())
			} else {
				logrus.Debug("Could not found hostname, get next offer")
				continue
			}
		}

		if !e.isAttributeMachted("__mc_placement_node_platform_os", "os", cmd, offer) {
			continue
		}

		if !e.isAttributeMachted("__mc_placement_node_platform_arch", "arch", cmd, offer) {
			continue
		}

		offerret = offer
	}
	// remove the offer we took
	offerIds = e.removeOffer(offerIds, offerret.ID.Value)
	return offerret, offerIds
}

// search matched mesos attributes
func (e *Scheduler) getAttributes(name string, offer mesosproto.Offer) string {
	for _, attribute := range offer.Attributes {
		if strings.EqualFold(attribute.Name, name) {
			return attribute.GetText().Value
		}
	}
	return ""
}

func (e *Scheduler) isAttributeMachted(label, attribute string, cmd cfg.Command, offer mesosproto.Offer) bool {
	valOS := e.getLabelValue(label, cmd)
	if valOS != "" {
		if strings.ToLower(valOS) == e.getAttributes(attribute, offer) {
			logrus.WithField("func", "mesos.isAttribute.Matched").Debugf("Set Server %s Constraint to: %s", attribute, offer.GetHostname())
			return true
		}
		logrus.WithField("func", "mesos.isAttribute.Matched").Debugf("Could not found %s, get next offer", attribute)
		return false
	}
	return true
}

// remove the offer we took from the list
func (e *Scheduler) removeOffer(offers []mesosproto.OfferID, clean string) []mesosproto.OfferID {
	var offerIds []mesosproto.OfferID
	for _, offer := range offers {
		if offer.Value != clean {
			offerIds = append(offerIds, offer)
		}
	}
	logrus.WithField("func", "mesos.removeOffer").Debug("Unused offers: ", offerIds)
	return offerIds
}
