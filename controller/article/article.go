package article

import (
	// "fmt"
	"strconv"
	"net/http"
	"github.com/xichengh/xcblog/model"
	"github.com/xichengh/xcblog/utils"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)
type PostArticle struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Status       uint          `json:"status"`
	Content      string        `json:"content"`
	Tags         []string      `json:"tags"`
	Logo       	 string        `json:"logo"`
}
type articleUpdate struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Status       uint          `json:"status"`
	Content      string        `json:"content"`
	Tags         []string      `json:"tags"`
	Logo       	 string        `json:"logo"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}
type articleRes struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	BrowseCount  uint          `json:"browseCount"`
	Status       uint          `json:"status"`
	Content      string        `json:"content"`
	Tags         []string      `json:"tags"`
	Author       author        `json:"author"`
	Logo       	 string        `json:"logo"`
}
type author struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name"`
	UserName  string        `json:"username"`
}
//获取文章详细数据
func GetArticleDetail(c *gin.Context) {
	var resData articleRes
	var authorData author
	id := c.Param("id")
	if findErr := model.DB.C("articles").Find(bson.M{
		"_id": bson.ObjectIdHex(id),
	}).One(&resData); findErr != nil {
		utils.SendBadResponse(c, "找不到该数据")
		return
	}
	if findErr := model.DB.C("users").Find(bson.M{
		"username": resData.Author.UserName,
	}).One(&authorData); findErr != nil {
		utils.SendBadResponse(c, "作者不存在")
		return
	}
	if err := model.DB.C("articles").Update(bson.M{
		"_id": bson.ObjectIdHex(id),
	}, bson.M{"$set": struct { BrowseCount  uint `json:"browseCount"` }{
		BrowseCount: resData.BrowseCount + 1,
	}}); err != nil {
		utils.SendBadResponse(c, "服务端错误")
		return
	}
	resData.Author = authorData
	c.JSON(http.StatusOK, gin.H{
		"data": resData,
		"msg": "success",
	})
}
//创建文章
func CreateArticle(c *gin.Context) {
	var postData PostArticle
	user, exist := c.Get("user")
	if(!exist) {
		utils.SendBadResponse(c, "未登录")
		return
	}
	if err := c.BindJSON(&postData); err != nil {
		utils.SendBadResponse(c, "服务端错误")
		return
	}
	articleData := model.Article{
		ID:          bson.NewObjectId(),
		CreatedAt:   time.Now(),
		UpdatedAt:	 time.Now(),
		Name:        postData.Name,
		Description: postData.Description,
		BrowseCount: 0,
		Status:      postData.Status,
		Content:     postData.Content,
		Tags:        postData.Tags,
		Author:      user.(model.Author),
		Logo:        postData.Logo,
	}
	if err := model.DB.C("articles").Insert(&articleData); err != nil {
		utils.SendBadResponse(c, "服务端错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": articleData.ID,
		"msg": "success",
	})
}
//删除文章
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	if findErr := model.DB.C("articles").Remove(bson.M{
		"_id": bson.ObjectIdHex(id),
	}); findErr != nil {
		utils.SendBadResponse(c, "找不到该数据")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": "",
		"msg": "success",
	})
}
//更新文章
func UpdateArticle(c *gin.Context) {
	var postData PostArticle
	id := c.Param("id")
	if err := c.BindJSON(&postData); err != nil {
		utils.SendBadResponse(c, "服务端错误")
		return
	}
	if err := model.DB.C("articles").Update(bson.M{
		"_id": bson.ObjectIdHex(id),
	}, bson.M{"$set": articleUpdate{
		UpdatedAt:	 time.Now(),
		Name:        postData.Name,
		Description: postData.Description,
		Status:      postData.Status,
		Content:     postData.Content,
		Tags:        postData.Tags,
		Logo:        postData.Logo,
	}}); err != nil {
		utils.SendBadResponse(c, "服务端错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": id,
		"msg": "success",
	})
}
//获取所有的文章列表数据
func GetAllArticleList(c *gin.Context) {
	articlelist := []articleRes{}
	allArticlelist := []articleRes{}
	var page, perPage int 
	var tag string
	pageQ:= c.Query("page")
	perPageQ := c.Query("perPage")
	tag = c.Query("tag")
	page, _ = strconv.Atoi(pageQ)
	perPage, _ = strconv.Atoi(perPageQ)
	if( page <= 0) {
		page = 1
	}
	if( perPage <= 0) {
		perPage = 10
	}
	findQuery := bson.M{}
	if( tag != "") {
		findQuery["tags"] = bson.M{
			"$elemMatch": bson.M{
				"$in": []string{ tag },
			},
		}
	}
	err := model.DB.C("articles").Find(findQuery).Sort("-createdat").Skip((page - 1) * perPage).Limit(perPage).All(&articlelist)
	allErr := model.DB.C("articles").Find(findQuery).All(&allArticlelist)
	
	if(err != nil || allErr != nil ) {
		utils.SendBadResponse(c, "服务端错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":articlelist,
			"count": len(allArticlelist),
		},
		"msg": "success",
	})
}
//获取公开的文章列表数据
func GetArticleList(c *gin.Context) {
	articlelist := []articleRes{}
	allArticlelist := []articleRes{}
	var page, perPage int 
	var tag string
	pageQ:= c.Query("page")
	perPageQ := c.Query("perPage")
	tag = c.Query("tag")
	page, _ = strconv.Atoi(pageQ)
	perPage, _ = strconv.Atoi(perPageQ)
	if( page <= 0) {
		page = 1
	}
	if( perPage <= 0) {
		perPage = 10
	}
	findQuery := bson.M{
		"status": 1,
	}
	if( tag != "") {
		findQuery["tags"] = bson.M{
			"$elemMatch": bson.M{
				"$in": []string{ tag },
			},
		}
	}
	err := model.DB.C("articles").Find(findQuery).Sort("-createdat").Skip((page - 1) * perPage).Limit(perPage).All(&articlelist)
	allErr := model.DB.C("articles").Find(findQuery).All(&allArticlelist)
	
	if(err != nil || allErr != nil ) {
		utils.SendBadResponse(c, "服务端错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":articlelist,
			"count": len(allArticlelist),
		},
		"msg": "success",
	})
}
