package main

import (
	"errors"
	"strings"
)

type ActionType string 

const (
	PUB ActionType = "PUB" 
	SUB ActionType = "SUB"
	EXIT ActionType = "EXIT"
	INVALID ActionType = ""
)

type ActionData []string

type pubData struct {
	topic Topic 
	message string 
}

type subData struct {
	topic Topic
}


func isValidAction(raw_action string) (ActionType, bool) {
	action := ActionType(raw_action)
	switch action {
	case PUB, SUB, EXIT:
		return action, true
	}
	return INVALID, false
}


func parseMessage(message string) (ActionType, ActionData, error) {
	words := strings.Fields(message)
	if len(words) == 0 {
		return INVALID, nil, errors.New("Empty message")
	}
	var (
		action ActionType
		valid bool
	)
	raw_action := words[0]
	if action, valid = isValidAction(raw_action); !valid {
		return action, ActionData{}, errors.New("Invalid action")
	}
	return action, words[1:], nil
}

func parsePubData(data ActionData) (pubData, error) {
	if len(data) < 2 {
		return pubData{}, errors.New("Invalid amount of data for PUB action")
	}
	return pubData{topic: Topic(data[0]), message: strings.Join(data[1:], " ")}, nil
}

func parseSubData(data ActionData) (subData, error) {
	if len(data) != 1 {
		return subData{}, errors.New("Invalid amount of data for SUB action")
	}
	return subData{topic: Topic(data[0])}, nil
}

