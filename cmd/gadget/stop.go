package main

import (
	//~ "fmt"
	"errors"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetStop(args []string, g *GadgetContext) error {
	
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info("  Stopping:")
	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	var stopFailed bool = false
	
	for _, container := range stagedContainers {
		log.Infof("    %s", container.Alias)
		
		stdout, stderr, err := RunRemoteCommand(client, "docker stop", container.Alias)
		
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"stop-stage": "stop",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"stop-stage": "stop",
		}).Debug(stderr)
		
		if err != nil {
			
			stopFailed = true
			
			log.WithFields(log.Fields{
				"function": "GadgetStop",
				"name": container.Alias,
				"stop-stage": "stop",
			}).Debug("This is likely due to specifying containers for a previous operation, but trying to stop all")


			log.Debug("Failed to stop container on Gadget,")
			log.Debug("it might have never been deployed,")
			log.Debug("Or stop otherwise failed")
			
		}

		stdout, stderr, err = RunRemoteCommand(client, "docker rm", container.Alias)
		
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"stop-stage": "rm",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"stop-stage": "rm",
		}).Debug(stderr)
		
		if err != nil {
			
			stopFailed = true
			
			log.WithFields(log.Fields{
				"function": "GadgetStop",
				"name": container.Alias,
				"stop-stage": "rm",
			}).Debug("This is likely due to specifying containers for a previous operation, but trying to stop all")


			log.Errorf("Failed to stop '%s' on Gadget", container.Name)
			log.Warn("Was it ever started?")
			
		} else {
			log.Info("    - stopped")
		}
	}
	
	if stopFailed {
		err = errors.New("A problem was encountered in GadgetStop")
	}
	
	return err
}
