package controller

import (
	"context"
	"fmt"
	"io"
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
	all, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	link := string(all)

	url, err := checkNormalizeUrl(link)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("incorrect url: %s - %v", link, err)})
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

	result := ctx.Request.URL
	result.Path = hash
	ctx.Status(http.StatusCreated)
	_, err = ctx.Writer.Write([]byte(result.String()))
	if err != nil {
		log.Printf("%v", err)
	}
}

func (c *Controller) Redirect(ctx *gin.Context) {
	var redirect struct {
		ID string `uri:"id"`
	}

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

	ctx.Header("Location", ctx.Request.RequestURI)
	ctx.Redirect(http.StatusTemporaryRedirect, link)
}
