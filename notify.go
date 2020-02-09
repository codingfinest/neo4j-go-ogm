package gogm

import (
	"errors"
	"reflect"
)

func notifyPreSaveGraph(g graph, eventer eventer, registry *registry) error {

	if g.getValue().IsValid() {
		for _, eventListener := range eventer.eventListeners {
			eventListener.OnPreSave(event{g.getValue(), -1})
		}

		var (
			metadata metadata
			err      error
		)
		if metadata, err = registry.get(g.getValue().Type()); err != nil {
			return err
		}
		var label string
		if label, err = metadata.getLabel(*g.getValue()); err != nil {
			return err
		}
		g.setLabel(label)
		g.setProperties(metadata.getProperties(*g.getValue()))

		customIDName, customIDValue := metadata.getCustomID(*g.getValue())
		if customIDName != emptyString && customIDValue.Type().Kind() == reflect.Ptr && customIDValue.IsNil() {
			return errors.New("Custom ID cannot be nil in " + g.getValue().Type().String())
		}

	}

	return nil
}

func notifyPostSave(eventer eventer, g graph, lifeCycle lifeCycle) error {
	if g == nil {
		return nil
	}
	if g.getValue().IsValid() {
		for _, eventListener := range eventer.eventListeners {
			eventListener.OnPostSave(event{g.getValue(), lifeCycle})
		}
	}
	return nil
}

func notifyPostDelete(eventer eventer, g graph, lifeCycle lifeCycle) error {
	if g == nil {
		return nil
	}

	g.setID(initialGraphID)
	if g.getValue().IsValid() {
		id := initialGraphID
		unloadGraphID(g, &id)
	}
	if g.getValue().IsValid() {
		for _, eventListener := range eventer.eventListeners {
			eventListener.OnPostDelete(event{g.getValue(), lifeCycle})
		}
	}

	return nil
}
