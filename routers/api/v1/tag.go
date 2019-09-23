package v1

import (
	"MyBlog/models"
	"MyBlog/pkg/e"
	"MyBlog/pkg/setting"
	"MyBlog/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

//获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//新增文章标签
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	code := e.INVALID_PARAMS
	util.Valid(func(v *validation.Validation) {
		v.Required(name, "name").Message("名称不能为空")
		v.MaxSize(name, 100, "name").Message("名称最长为100字符")
		v.Required(createdBy, "created_by").Message("创建人不能为空")
		v.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
		v.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}, func() {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

//修改文章标签
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	code := e.INVALID_PARAMS
	var state = -1
	util.Valid(func(v *validation.Validation) {
		if arg := c.Query("state"); arg != "" {
			state = com.StrTo(arg).MustInt()
			v.Range(state, 0, 1, "state").Message("状态只允许0或1")
		}

		v.Required(id, "id").Message("ID不能为空")
		v.Required(modifiedBy, "modified_by").Message("修改人不能为空")
		v.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
		v.MaxSize(name, 100, "name").Message("名称最长为100字符")

	}, func() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}

			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

//删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	code := e.INVALID_PARAMS
	util.Valid(func(v *validation.Validation) {
		v.Min(id, 1, "id").Message("ID必须大于0")
	}, func() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
