package scheduler

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
)

// HandleOffers will handle the offers event of mesos
func (e *Scheduler) HandleOffers(offers *mesosproto.Event_Offers) error {
	var offerIds []*mesosproto.OfferID
	select {
	case cmd := <-e.Framework.CommandChan:
		// if no taskid or taskname is given, it's a wrong task.
		if cmd.TaskID == "" || cmd.TaskName == "" || cmd.Killed {
			return nil
		}

		// before schedule the task, check again if there are already
		// enough instances
		if e.Redis.CountRedisKey(cmd.TaskName+":*", "__KILL") > cmd.Instances {
			return nil
		}

		task := e.Redis.GetTaskFromTaskID(cmd.TaskID)

		var takeOffer *mesosproto.Offer
		takeOffer, offerIds = e.getOffer(offers, task)
		if takeOffer.GetHostname() == "" {
			task.State = ""
			e.Redis.SaveTaskRedis(task)
			logrus.WithField("func", "mesos.HandleOffers").Debug("No matched offer found.")
			return nil
		}
		logrus.WithField("func", "scheduler.HandleOffers").Info("Take Offer from " + takeOffer.GetHostname() + " for task " + task.TaskID + " (" + task.TaskName + ")")

		var taskInfo []*mesosproto.TaskInfo
		RefuseSeconds := 5.0

		taskInfo = e.PrepareTaskInfoExecuteContainer(takeOffer.GetAgentId(), task)

		// build mesos call object
		accept := &mesosproto.Call{
			Type: mesosproto.Call_ACCEPT.Enum(),
			Accept: &mesosproto.Call_Accept{
				OfferIds: []*mesosproto.OfferID{{
					Value: takeOffer.Id.Value,
				}},
				Filters: &mesosproto.Filters{
					RefuseSeconds: &RefuseSeconds,
				},
				Operations: []*mesosproto.Offer_Operation{{
					Type: mesosproto.Offer_Operation_LAUNCH.Enum(),
					Launch: &mesosproto.Offer_Operation_Launch{
						TaskInfos: taskInfo,
					}}}}}

		e.Redis.SaveTaskRedis(task)

		logrus.WithField("func", "scheduler.HandleOffers").Debug("Offer Accept: ", takeOffer.GetId(), " On Node: ", takeOffer.GetHostname())

		err := e.Mesos.Call(accept)
		if err != nil {
			logrus.WithField("func", "scheduler.HandleOffers").Error(err.Error())
			return err
		}

		// decline unneeded offer
		if len(offerIds) > 0 {
			logrus.WithField("func", "scheduler.HandleOffer").Debug("Offer Decline: ", offerIds)
			go e.Mesos.Call(e.Mesos.DeclineOffer(offerIds))
		}
	default:
		offerIds = e.getAllOfferIDs(offers)
	}

	e.Mesos.Call(e.Mesos.DeclineOffer(offerIds))
	return nil
}

// get the value of a label from the command
func (e *Scheduler) getLabelValue(label string, cmd *cfg.Command) string {
	for _, v := range cmd.Labels {
		if label == *v.Key {
			return fmt.Sprint(v.GetValue())
		}
	}
	return ""
}

func (e *Scheduler) getAllOfferIDs(offers *mesosproto.Event_Offers) []*mesosproto.OfferID {
	var offerIds []*mesosproto.OfferID
	for _, offer := range offers.Offers {
		offerIds = append(offerIds, offer.Id)
	}

	return offerIds
}

func (e *Scheduler) getOffer(offers *mesosproto.Event_Offers, cmd *cfg.Command) (*mesosproto.Offer, []*mesosproto.OfferID) {
	var offerret *mesosproto.Offer

	offerIds := e.getAllOfferIDs(offers)

	// check all offers
	for _, offer := range offers.Offers {
		logrus.WithField("func", "scheduler.getOffer").Debug("Got Offer From:", offer.GetHostname())

		// if contraint_hostname is set, only accept offer with the same hostname
		valHostname := e.getLabelValue("__mc_placement_node_hostname", cmd)
		if valHostname != "" {
			if strings.ToLower(valHostname) == offer.GetHostname() {
				logrus.WithField("func", "scheduler.getOffer").Debug("Set Server Hostname Constraint to:", offer.GetHostname())
			} else {
				logrus.WithField("func", "scheduler.getOffer").Debug("Could not find hostname, get next offer")
				continue
			}
		}

		if e.getLabelValue("__mc_placement", cmd) == "unique" {
			if e.alreadyRunningOnHostname(cmd) {
				logrus.WithField("func", "scheduler.getOffer").Debug("UNIQUE: Already running on node: ", offer.GetHostname())
				continue
			}
		}

		if !e.isAttributeMachted("__mc_placement_node_platform_os", "os", cmd, offer) {
			logrus.WithField("func", "scheduler.getOffer").Debug("OS: Does not match Attribute")
			continue
		}

		if !e.isAttributeMachted("__mc_placement_node_platform_arch", "arch", cmd, offer) {
			logrus.WithField("func", "scheduler.getOffer").Debug("ARCH: Does not match Attribute")
			continue
		}

		// if the ressources of this offer does not matched what the command need, the skip
		if !e.Mesos.IsRessourceMatched(offer.Resources, cmd) {
			logrus.WithField("func", "scheduler.getOffer").Debug("Could not found any matched ressources, get next offer")
			continue
		}

		offerret = offer
		break
	}
	// remove the offer we took
	if offerret.GetHostname() != "" {
		offerIds = e.removeOffer(offerIds, offerret.GetId().GetValue())
	}
	e.Mesos.Call(e.Mesos.DeclineOffer(offerIds))
	return offerret, offerIds
}

// search matched mesos attributes
func (e *Scheduler) getAttributes(name string, offer *mesosproto.Offer) string {
	for _, attribute := range offer.Attributes {
		if strings.EqualFold(*attribute.Name, name) {
			return attribute.GetText().GetValue()
		}
	}
	return ""
}

func (e *Scheduler) alreadyRunningOnHostname(cmd *cfg.Command) bool {
	keys := e.Redis.GetAllRedisKeys(cmd.TaskName + ":*")
	for keys.Next(e.Redis.CTX) {
		// continue if the key is not a mesos task
		if e.Redis.CheckIfNotTask(keys) {
			continue
		}
		// get the values of the current key
		key := e.Redis.GetRedisKey(keys.Val())

		task := e.Mesos.DecodeTask(key)

		// continue if it's a unvalid task
		if task.TaskID == "" {
			continue
		}

		if task.MesosAgent.Hostname == cmd.MesosAgent.Hostname && task.TaskID != cmd.TaskID {
			return true
		}
	}

	return false
}

func (e *Scheduler) isAttributeMachted(label, attribute string, cmd *cfg.Command, offer *mesosproto.Offer) bool {
	valOS := e.getLabelValue(label, cmd)
	if valOS != "" {
		if strings.ToLower(valOS) == e.getAttributes(attribute, offer) {
			logrus.WithField("func", "scheduler.isAttribute.Matched").Debugf("Set Server %s Constraint to: %s", attribute, offer.GetHostname())
			return true
		}
		logrus.WithField("func", "scheduler.isAttribute.Matched").Debugf("Could not found %s, get next offer", attribute)
		return false
	}
	return true
}

// remove the offer we took from the list
func (e *Scheduler) removeOffer(offers []*mesosproto.OfferID, clean string) []*mesosproto.OfferID {
	var offerIds []*mesosproto.OfferID
	for _, offer := range offers {
		if *offer.Value != clean {
			offerIds = append(offerIds, offer)
		}
	}
	logrus.WithField("func", "scheduler.removeOffer").Debug("Unused offers: ", offerIds)
	return offerIds
}
