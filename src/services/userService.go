package services

import (
	"fmt"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
)

type UserService struct {
}

func NewUserService(cfg *config.Config) *UserService {
	return &UserService{}
}

func (us *UserService) Init(req *dto.UserDTO) (err error) {
	fmt.Println(req)

	return nil
}
