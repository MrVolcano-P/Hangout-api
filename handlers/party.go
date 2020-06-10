package handlers

import (
	"fmt"
	"hangout-api/context"
	"hangout-api/models"
	Union "hangout-api/union"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Party struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Membership string    `json:"membership"`
	Date       time.Time `json:"date"`
	Members    []Member  `json:"members"`
	Leader     UserParty `json:"leader"`
}
type PartyReq struct {
	Name       string    `json:"name"`
	Membership string    `json:"membership"`
	Date       time.Time `json:"date"`
}

func (h *Handler) CreateParty(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		c.Status(401)
		return
	}
	pubId := c.Param("id")
	id, err := strconv.Atoi(pubId)
	if err != nil {
		Error(c, 400, err)
		return
	}
	req := new(PartyReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	party := new(models.Party)
	party.Name = req.Name
	party.Membership = req.Membership
	party.Date = req.Date
	party.PubID = uint(id)
	party.UserID = user.ID

	err = h.pts.Create(party)
	if err != nil {
		Error(c, 500, err)
		return
	}

	member := new(models.Member)
	member.UserID = user.ID
	member.PartyID = party.ID

	err = h.ms.Create(member)
	if err != nil {
		Error(c, 500, err)
		return
	}

	members := []Member{}
	members = append(members, Member{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	})

	res := new(Party)
	res.ID = party.ID
	res.Name = party.Name
	res.Date = party.Date
	res.Membership = party.Membership
	res.Leader = UserParty{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}
	res.Members = members

	c.JSON(201, res)
}

func (h *Handler) GetPartiesBypubID(c *gin.Context) {
	pubId := c.Param("id")
	id, err := strconv.Atoi(pubId)
	if err != nil {
		Error(c, 400, err)
		return
	}
	parties, err := h.pts.GetPartiesByPubId(uint(id))
	if err != nil {
		Error(c, 500, err)
		return
	}
	resParties := []Party{}
	for _, party := range parties {
		resMembers := []Member{}
		members, err := h.ms.GetByPartyID(party.ID)
		if err != nil {
			Error(c, 500, err)
			return
		}
		for _, member := range members {
			user, err := h.us.GetByID(member.UserID)
			if err != nil {
				Error(c, 500, err)
				return
			}
			resMembers = append(resMembers, Member{
				ID:       user.ID,
				Username: user.Username,
				Name:     user.Name,
			})
		}
		leader, err := h.us.GetByID(party.UserID)
		if err != nil {
			Error(c, 500, err)
			return
		}

		resParties = append(resParties, Party{
			ID:         party.ID,
			Name:       party.Name,
			Date:       party.Date,
			Membership: party.Membership,
			Leader: UserParty{
				ID:       leader.ID,
				Username: leader.Username,
				Name:     leader.Name,
			},
			Members: resMembers,
		})
	}
	c.JSON(200, resParties)
}

func (h *Handler) GetPartiesByuserID(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		c.Status(401)
		return
	}
	// from party
	parties, err := h.pts.GetPartiesByUserId(user.ID)
	if err != nil {
		Error(c, 500, err)
		return
	}
	partyId := []uint{}
	for _, party := range parties {
		partyId = append(partyId, party.ID)
	}
	// from member
	memberList, err := h.ms.GetByUserID(user.ID)
	if err != nil {
		Error(c, 500, err)
		return
	}
	partyIDFromMember := []uint{}
	for _, member := range memberList {
		partyIDFromMember = append(partyIDFromMember, member.PartyID)
	}

	unionParty := Union.Union(partyId, partyIDFromMember)
	sort.Slice(unionParty, func(i, j int) bool { return unionParty[i] < unionParty[j] })

	resParties := []Party{}
	for _, id := range unionParty {
		fmt.Println(id)
		party, err := h.pts.GetPartyById(id)
		if err != nil {
			Error(c, 500, err)
			return
		}

		resMembers := []Member{}
		members, err := h.ms.GetByPartyID(id)
		if err != nil {
			Error(c, 500, err)
			return
		}
		for _, member := range members {
			user, err := h.us.GetByID(member.UserID)
			if err != nil {
				Error(c, 500, err)
				return
			}
			resMembers = append(resMembers, Member{
				ID:       user.ID,
				Username: user.Username,
				Name:     user.Name,
			})
		}
		leader, err := h.us.GetByID(party.UserID)
		if err != nil {
			Error(c, 500, err)
			return
		}
		resParties = append(resParties, Party{
			ID:         party.ID,
			Name:       party.Name,
			Date:       party.Date,
			Membership: party.Membership,
			Leader: UserParty{
				ID:       leader.ID,
				Username: leader.Username,
				Name:     leader.Name,
			},
			Members: resMembers,
		})
	}
	
	c.JSON(200, resParties)
}

func (h *Handler) JoinParty(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		c.Status(401)
		return
	}
	partyId := c.Param("id")
	id, err := strconv.Atoi(partyId)
	if err != nil {
		Error(c, 400, err)
		return
	}
	member := new(models.Member)
	member.UserID = user.ID
	member.PartyID = uint(id)

	err = h.ms.Create(member)
	if err != nil {
		Error(c, 500, err)
		return
	}

	party, err := h.pts.GetPartyById(uint(id))
	if err != nil {
		Error(c, 500, err)
		return
	}

	resMembers := []Member{}
	members, err := h.ms.GetByPartyID(party.PubID)
	if err != nil {
		Error(c, 500, err)
		return
	}
	for _, member := range members {
		user, err := h.us.GetByID(member.UserID)
		if err != nil {
			Error(c, 500, err)
			return
		}
		resMembers = append(resMembers, Member{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
		})
	}
	leader, err := h.us.GetByID(party.UserID)
	if err != nil {
		Error(c, 500, err)
		return
	}
	res := new(Party)
	res.ID = party.ID
	res.Name = party.Name
	res.Date = party.Date
	res.Membership = party.Membership
	res.Leader = UserParty{
		ID:       leader.ID,
		Username: leader.Username,
		Name:     leader.Name,
	}
	res.Members = resMembers

	c.JSON(200, res)
}
