package node

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"legato_server/internal/legato/api/rest/auth"
	"legato_server/internal/legato/api/rest/server"
	"legato_server/internal/legato/database"
	"legato_server/internal/legato/database/models"
	"legato_server/pkg/logger"
	"net/http"
	"strconv"
)

var log, _ = logger.NewLogger(logger.Config{})

type Node struct {
	db database.Database
}

func (n *Node) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/users/:username/scenarios/:scenario_id/nodes", n.AddNode)
	group.PUT("/users/:username/scenarios/:scenario_id/nodes/:node_id", n.UpdateNode)
	group.DELETE("/users/:username/scenarios/:scenario_id/nodes/:node_id", n.DeleteNode)
	group.GET("/users/:username/scenarios/:scenario_id/nodes/:node_id", n.GetNode)
	group.GET("/users/:username/scenarios/:scenario_id/nodes", n.GetScenarioNodes)
	group.GET("/users/:username/scenarios/:scenario_id/nodes/:node_id/children", n.GetNodesChildren)
}

func (n *Node) GetNodesChildren(c *gin.Context) {
	// Params
	username := c.Param("username")
	scenarioIdParam, _ := strconv.Atoi(c.Param("scenario_id"))
	scenarioId := uint(scenarioIdParam)
	nodeIdParam, _ := strconv.Atoi(c.Param("node_id"))
	nodeId := uint(nodeIdParam)

	// Auth
	loggedInUser := auth.CheckAuth(c, []string{username})
	if loggedInUser == nil {
		return
	}

	scenario, err := n.db.GetUserScenarioById(loggedInUser, scenarioId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not find this scenario",
			"error":   err.Error(),
		})
		return
	}

	serv, err := n.db.GetScenarioServiceById(&scenario, nodeId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not find this node in the scenario",
			"error":   err.Error(),
		})
		return
	}

	// TODO: FROM HERE
	nodesChildren, err := n.db.GetServiceChildrenById(&serv)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not fetch the children",
			"error":   err.Error(),
		})
		return
	}

	// Convert to Response model
	var nodesResponse []ServiceNodeResponse
	for _, srv := range nodesChildren {
		nodesResponse = append(nodesResponse, ServiceNodeResponse{
			Id:       srv.ID,
			ParentId: srv.ParentID,
			Name:     srv.Name,
			Type:     srv.Type,
			SubType:  srv.SubType,
			Position: Position{
				X: srv.PosX,
				Y: srv.PosY,
			},
			Data: srv.Data,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"children": nodesResponse,
	})
}

func (n *Node) GetScenarioNodes(c *gin.Context) {
	// Params
	username := c.Param("username")
	scenarioIdParam, _ := strconv.Atoi(c.Param("scenario_id"))
	scenarioId := uint(scenarioIdParam)

	// Queries
	onlyRootNodesQuery := c.DefaultQuery("only_root", "false")
	onlyRootNodes, _ := strconv.ParseBool(onlyRootNodesQuery)

	// Auth
	loginUser := auth.CheckAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	scenario, err := n.db.GetUserScenarioById(loginUser, scenarioId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not fetch this scenario",
			"error":   err.Error(),
		})
		return
	}
	services := scenario.Services

	if onlyRootNodes {
		services, err = n.db.GetScenarioRootServices(scenario)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "can not fetch root services for this scenario",
				"error":   err.Error(),
			})
			return
		}
	}

	// Convert to Response model
	var nodesResponse []ServiceNodeResponse
	for _, srv := range services {
		nodesResponse = append(nodesResponse, ServiceNodeResponse{
			Id:       srv.ID,
			ParentId: srv.ParentID,
			Name:     srv.Name,
			Type:     srv.Type,
			SubType:  srv.SubType,
			Position: Position{
				X: srv.PosX,
				Y: srv.PosY,
			},
			Data: srv.Data,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodesResponse,
	})
}

func (n *Node) AddNode(c *gin.Context) {
	username := c.Param("username")
	scenarioIdParam, _ := strconv.Atoi(c.Param("scenario_id"))
	scenarioId := uint(scenarioIdParam)

	// Auth
	loggedInUser := auth.CheckAuth(c, []string{username})
	if loggedInUser == nil {
		return
	}

	// Validate JSON
	newNode := NewServiceNodeRequest{}
	if err := c.ShouldBindJSON(&newNode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not create this scenario",
			"error":   err.Error(),
		})
		return
	}

	// Service Switch
	// NOTE: handle other non-service state
	var err error
	var addedServ ServiceNodeResponse
	switch newNode.Type {
	//case "webhooks":
	//	addedServ, err = resolvers.WebhookUseCase.AddWebhookToScenario(loggedInUser, uint(scenarioId), newNode)
	//	break
	case "https":
		//addedServ, err = resolvers.HttpUserCase.AddToScenario(loggedInUser, uint(scenarioId), newNode)

		user, err := n.db.GetUserByUsername(loggedInUser.Username)
		if err != nil {
			break
		}

		scenario, err := n.db.GetUserScenarioById(&user, scenarioId)
		if err != nil {
			break
		}

		httpService := models.Service{
			Name:     newNode.Name,
			Type:     "https",
			ParentID: newNode.ParentId,
			PosX:     newNode.Position.X,
			PosY:     newNode.Position.Y,
		}

		addedHttpService, err := n.db.AddNodeToScenario(&scenario, httpService)
		if err != nil {
			return
		}

		addedServ = ServiceNodeResponse{
			Id:       addedHttpService.ID,
			ParentId: addedHttpService.ParentID,
			Name:     addedHttpService.Name,
			Type:     addedHttpService.Type,
			SubType:  addedHttpService.SubType,
			Position: Position{
				X: addedHttpService.PosX,
				Y: addedHttpService.PosY,
			},
			Data: addedHttpService.Data,
		}

		break
	//case "telegrams":
	//	addedServ, err = resolvers.TelegramUseCase.AddToScenario(loggedInUser, uint(scenarioId), newNode)
	//	break
	//case "spotifies":
	//	addedServ, err = resolvers.SpotifyUseCase.AddToScenario(loggedInUser, uint(scenarioId), newNode)
	//	break
	//case "sshes":
	//	addedServ, err = resolvers.SshUseCase.AddToScenario(loggedInUser, uint(scenarioId), newNode)
	//	break
	//case "gmails":
	//	addedServ, err = resolvers.GmailUseCase.AddToScenario(loggedInUser, uint(scenarioId), newNode)
	//	break
	//case "githubs":
	//	addedServ, err = resolvers.GithubUseCase.AddToScenario(loggedInUser, uint(scenarioId), newNode)
	//	break
	//case "discords":
	//	addedServ, err = resolvers.DiscordUseCase.AddToScenario(loggedInUser, uint(scenarioId), newNode)
	//	break
	//case "tool_boxes":
	//	addedServ, err = resolvers.ToolBoxUseCase.AddToScenario(loggedInUser, uint(scenarioId), newNode)
	//	break

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not create this node",
			"error":   fmt.Sprintf("there is not any service with name %s", newNode.Type),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not create this node",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "node is created successfully.",
		"node":    addedServ,
	})
}

func (n *Node) GetNode(c *gin.Context) {
	username := c.Param("username")
	scenarioIdParam, _ := strconv.Atoi(c.Param("scenario_id"))
	scenarioId := uint(scenarioIdParam)
	nodeIdParam, _ := strconv.Atoi(c.Param("node_id"))
	nodeId := uint(nodeIdParam)

	// Auth
	loggedInUser := auth.CheckAuth(c, []string{username})
	if loggedInUser == nil {
		return
	}

	scenario, err := n.db.GetUserScenarioById(loggedInUser, scenarioId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not find this scenario",
			"error":   err.Error(),
		})
		return
	}

	service, err := n.db.GetScenarioServiceById(&scenario, nodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not fetch this node: %s", err),
			"error":   err.Error(),
		})
		return
	}

	node := ServiceNodeResponse{
		Id:       service.ID,
		ParentId: service.ParentID,
		Name:     service.Name,
		Type:     service.Type,
		SubType:  service.SubType,
		Position: Position{
			X: service.PosX,
			Y: service.PosY,
		},
		Data: service.Data,
	}

	c.JSON(http.StatusOK, gin.H{
		"node": node,
	})
}

func (n *Node) UpdateNode(c *gin.Context) {
	username := c.Param("username")
	scenarioIdParam, _ := strconv.Atoi(c.Param("scenario_id"))
	scenarioId := uint(scenarioIdParam)
	nodeIdParam, _ := strconv.Atoi(c.Param("node_id"))
	nodeId := uint(nodeIdParam)

	// Auth
	loggedInUser := auth.CheckAuth(c, []string{username})
	if loggedInUser == nil {
		return
	}

	newNode := UpdateServiceNodeRequest{}
	if err := c.ShouldBindJSON(&newNode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not update this scenario",
			"error":   err.Error(),
		})
		return
	}

	// Get the existing service and get the type
	scenario, err := n.db.GetUserScenarioById(loggedInUser, scenarioId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not find this scenario",
			"error":   err.Error(),
		})
		return
	}

	serv, err := n.db.GetScenarioServiceById(&scenario, nodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not fetch this node: %s", err),
			"error":   err.Error(),
		})
		return
	}

	// Service Switch
	// NOTE: handle other non-service state
	switch newNode.Type {
	//case "webhooks":
	//	err = resolvers.WebhookUseCase.Update(loggedInUser, uint(scenarioId), uint(nodeId), newNode)
	//	break
	case "https":
		newHttp := models.Service{
			Name:     newNode.Name,
			Type:     newNode.Type,
			SubType:  newNode.SubType,
			ParentID: newNode.ParentId,
			PosX:     newNode.Position.X,
			PosY:     newNode.Position.Y,
			Data:     newNode.Data,
		}

		err = n.db.UpdateScenarioNode(&scenario, serv.ID, newHttp)
		break
	//case "telegrams":
	//	err = resolvers.TelegramUseCase.Update(loggedInUser, uint(scenarioId), uint(nodeId), newNode)
	//	break
	//case "spotifies":
	//	err = resolvers.SpotifyUseCase.Update(loggedInUser, uint(scenarioId), uint(nodeId), newNode)
	//	break
	//case "sshes":
	//	err = resolvers.SshUseCase.Update(loggedInUser, uint(scenarioId), uint(nodeId), newNode)
	//	break
	//case "gmails":
	//	err = resolvers.GmailUseCase.Update(loggedInUser, uint(scenarioId), uint(nodeId), newNode)
	//	break
	//case "githubs":
	//	err = resolvers.GithubUseCase.Update(loggedInUser, uint(scenarioId), uint(nodeId), newNode)
	//	break
	//case "discords":
	//	err = resolvers.DiscordUseCase.Update(loggedInUser, uint(scenarioId), uint(nodeId), newNode)
	//	break
	//case "tool_boxes":
	//	err = resolvers.ToolBoxUseCase.Update(loggedInUser, uint(scenarioId), uint(nodeId), newNode)
	//	break
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not update the node",
			"error":   fmt.Sprintf("there is not any service with name %s", newNode.Type),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not update this node",
			"error":   err.Error(),
		})
		return
	}

	updatedService, err := n.db.GetScenarioServiceById(&scenario, nodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not fetch updated node",
			"error":   err.Error(),
		})
		return
	}

	updatedNode := ServiceNodeResponse{
		Id:       updatedService.ID,
		ParentId: updatedService.ParentID,
		Name:     updatedService.Name,
		Type:     updatedService.Type,
		SubType:  updatedService.SubType,
		Position: Position{
			X: updatedService.PosX,
			Y: updatedService.PosY,
		},
		Data: updatedService.Data,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "node is updated successfully.",
		"node":    updatedNode,
	})
}

func (n *Node) DeleteNode(c *gin.Context) {
	// Params
	username := c.Param("username")
	scenarioIdParam, _ := strconv.Atoi(c.Param("scenario_id"))
	scenarioId := uint(scenarioIdParam)
	nodeIdParam, _ := strconv.Atoi(c.Param("node_id"))
	nodeId := uint(nodeIdParam)

	// Auth
	loggedInUser := auth.CheckAuth(c, []string{username})
	if loggedInUser == nil {
		return
	}

	scenario, err := n.db.GetUserScenarioById(loggedInUser, scenarioId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not find this scenario",
			"error":   err.Error(),
		})
		return
	}

	err = n.db.DeleteServiceById(&scenario, nodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not delete this node",
			"error":   err.Error(),
		})
		return
	}

	scenario, err = n.db.GetUserScenarioById(loggedInUser, scenarioId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not find this scenario",
			"error":   err.Error(),
		})
		return
	}

	// Convert to Response model
	var nodesResponse []ServiceNodeResponse
	for _, srv := range scenario.Services {
		nodesResponse = append(nodesResponse, ServiceNodeResponse{
			Id:       srv.ID,
			ParentId: srv.ParentID,
			Name:     srv.Name,
			Type:     srv.Type,
			SubType:  srv.SubType,
			Position: Position{
				X: srv.PosX,
				Y: srv.PosY,
			},
			Data: srv.Data,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "node is deleted successfully",
		"nodes":   nodesResponse,
	})
}

func NewNodeModule(db database.Database) (server.RestModule, error) {
	return &Node{
		db: db,
	}, nil
}
