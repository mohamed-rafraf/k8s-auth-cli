package cmd

import (
	"gorm.io/gorm"
)

/*type Cluster struct {
	Name string `json:"name"`
	Status bool `json:"status"`
	Token string `json:"token"`
	API string `json:"apiserver"`
}


type ClusterUser struct {
	Name string `json:"name"`
	FullName string `json:"fullname"`
	Email string  `json:"mail"`
	Clusters []string `json:"clusters"`
}


}*/

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserCluster struct {
	gorm.Model
	ClusterID uint
	UserID    uint
}

type ClusterUser struct {
	gorm.Model
	Name     string    `gorm:"not null"`
	FullName string    `gorm:"not null"`
	Email    string    `gorm:"unique;not null"`
	Clusters []Cluster `gorm:"many2many:user_clusters;"`
}

type Admin struct {
	gorm.Model
	Name     string
	Mail     string `gorm:"unique"`
	FullName string
}

type Cluster struct {
	gorm.Model
	Name   string `gorm:"unique;not null"`
	Status bool
	Token  string
	API    string
}

type ClusterResponse struct {
	Status   string    `json:"status"`
	Message  string    `json:"message"`
	Clusters []Cluster `json:"clusters"`
}
type UserResponse struct {
	Status       string        `json:"status"`
	Message      string        `json:"message"`
	ClusterUsers []ClusterUser `json:"clusterusers"`
}
