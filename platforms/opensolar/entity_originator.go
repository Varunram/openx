package opensolar

import (
	"fmt"

	database "github.com/YaleOpenLab/openx/database"
	utils "github.com/YaleOpenLab/openx/utils"
)

// An originator is someone who approaches the recipient in real life and proposes
// that he can start a contract on the opensolar platform that will be open ot investors
// he needs to make clear that he is an originator and if interested, he can volunteer
// to be the contractor as well, in which case there will be no auction and we can go
// straight ahead to the auction phase with investors investing in the contract. A MOU
// must also be sigend between the originator and the recipient defining terms of agreement
// as per legal standards

// TODO: Consider any other information needed for originators that should be added

func NewOriginator(uname string, pwd string, seedpwd string, Name string,
	Address string, Description string) (Entity, error) {
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "originator")
}

// Originate creates and saves a new origin contract
func (contractor *Entity) Originate(panelSize string, totalValue float64, location string,
	years int, metadata string, recIndex int, auctionType string) (Project, error) {

	var pc Project
	var err error

	indexCheck, err := RetrieveAllProjects()
	if err != nil {
		return pc, fmt.Errorf("Projects could not be retrieved!")
	}
	pc.Index = len(indexCheck) + 1
	pc.PanelSize = panelSize
	pc.TotalValue = totalValue
	pc.Location = location
	pc.Years = years
	pc.Metadata = metadata
	pc.DateInitiated = utils.Timestamp()
	iRecipient, err := database.RetrieveRecipient(recIndex)
	if err != nil { // recipient does not exist
		return pc, err
	}
	pc.RecipientIndex = iRecipient.U.Index
	pc.Stage = 0 // 0 since we need to filter this out while retrieving the propsoed contracts
	pc.AuctionType = auctionType
	pc.Originator = *contractor
	pc.Reputation = totalValue // reputation is equal to the total value of the project
	// instead of storing in this proposedcontracts slice, store it as a project, but not a contract and retrieve by stage
	err = pc.Save()
	// don't insert the project since the contractor's projects are not final
	return pc, err
}

// RepOriginatedProject adds reputation to an originator on successful origination of a contract
func RepOriginatedProject(origIndex int, projIndex int) error {
	originator, err := RetrieveEntity(origIndex)
	if err != nil {
		return err
	}
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	return originator.U.IncreaseReputation(project.TotalValue * OriginatorWeight)
}