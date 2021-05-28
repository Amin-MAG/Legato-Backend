package usecase

import (
	"legato_server/api"
	legatoDb "legato_server/db"
	"legato_server/domain"
	"legato_server/helper/converter"
	"time"
)

type GitUseCase struct {
	db             *legatoDb.LegatoDB
	contextTimeout time.Duration
}

func NewHttpUseCase(db *legatoDb.LegatoDB, timeout time.Duration) domain.GitUseCase {
	return &GitUseCase{
		db:             db,
		contextTimeout: timeout,
	}
}
func (ss *GitUseCase) GetGitWithId(gid uint, username string) (api.GitInfo, error) {
	user, err := ss.db.GetUserByUsername(username)
	if err != nil {
		return api.GitInfo{}, err
	}

	git, err := ss.db.GetGitByID(gid, &user)
	if err != nil {
		return api.GitInfo{}, err
	}

	return converter.GitDbToGitInfo(&git), nil

}

func (git *GitUseCase) AddToScenario(u *api.UserInfo, scenarioId uint, ns api.NewServiceNode) (api.ServiceNode, error) {

	user, err := git.db.GetUserByUsername(u.Username)
	if err != nil {
		return api.ServiceNode{}, err
	}

	scenario, err := git.db.GetUserScenarioById(&user, scenarioId)
	if err != nil {
		return api.ServiceNode{}, err
	}

	var gg legatoDb.Github
	gg.Service = converter.NewServiceNodeToServiceDb(ns)
	g, err := git.db.CreateGitForScenario(&scenario, gg)
	if err != nil {
		return api.ServiceNode{}, err
	}
	return converter.ServiceDbToServiceNode(g.Service), nil
}

func (git *GitUseCase) Update(u *api.UserInfo, scenarioId uint, serviceId uint, ns api.NewServiceNode) error {
	user, err := git.db.GetUserByUsername(u.Username)
	if err != nil {
		return err
	}
	scenario, err := git.db.GetUserScenarioById(&user, scenarioId)
	if err != nil {
		return err
	}

	var g legatoDb.Github
	g.Service = converter.NewServiceNodeToServiceDb(ns)
	err = git.db.UpdateGit(&scenario, serviceId, g)
	if err != nil {
		return err
	}

	return nil
}
