package services

import (
	"legato_server/internal/legato/database"
	"legato_server/internal/legato/database/models"
)

type Pipeline struct {
	db           database.Database
	RootServices []Service
}

// Start executes a Pipeline
func (p *Pipeline) Start() {
	//err = legatoDb.CreateHistory(s.ID)
	//if err != nil {
	//	return err
	//}

	log.Debugln("Executing root services of this scenario")
	for _, serv := range p.RootServices {
		go func(s Service) {
			s.Execute()
		}(serv)
	}
	log.Debugln("Executing finished")
}

func NewPipeline(db database.Database, rootServiceModels []models.Service) (*Pipeline, error) {
	var rootServices []Service
	for _, rsm := range rootServiceModels {
		service, err := NewService(db, rsm)
		if err != nil {
			return nil, err
		}

		rootServices = append(rootServices, service)
	}

	return &Pipeline{
		db:           db,
		RootServices: rootServices,
	}, nil
}
