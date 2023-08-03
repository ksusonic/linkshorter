package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	storage Storage
}

type Storage interface {
	GetLink(ctx context.Context, id string) (string, error)
	Insert(ctx context.Context, id string, url string) error
}

func NewUrlController(storage Storage) *Controller {
	return &Controller{storage: storage}
}

func (c *Controller) Shorten(ctx *gin.Context) {
	var body Shorten
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	url, err := checkNormalizeUrl(body.Link)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("incorrect url: %s - %v", body.Link, err)})
		return
	}
	hash, err := hashUrl(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("url hashing error:  %v", err)})
		return
	}
	err = c.storage.Insert(ctx.Request.Context(), hash, url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("creating url error:  %v", err)})
		log.Printf("got error: %v", err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (c *Controller) Redirect(ctx *gin.Context) {
	var redirect Redirect
	if err := ctx.BindUri(&redirect); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	link, err := c.storage.GetLink(ctx.Request.Context(), redirect.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("redirect error:  %v", err)})
		log.Printf("got error: %v", err)
		return
	}
	if link == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, link)
}
