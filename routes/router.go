package routes

import (
	"github.com/gin-gonic/gin"
	v1 "qiublog/api/v1"
	"qiublog/middleware"
	"qiublog/utils"
)

func InitRouter() {
	gin.SetMode(utils.Config.Server.AppMode)
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Use(gin.Recovery())
	if utils.Config.Server.AppMode == "release" {
		r.LoadHTMLGlob("web/index.html")
		r.Static("/assets", "web/assets")
		r.Static("favicons", "web/favicons")
		r.Static("img", "web/img")
		r.Static("effect", "web/effect")
		r.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", nil)
		})
	}

	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken(true))
	{
		//文章
		auth.POST("article/add", Handler()(v1.ReleaseArticle))  //发布文章
		auth.PUT("article/:id", Handler()(v1.ModifyArticle))    //更新文章
		auth.DELETE("article/:id", Handler()(v1.DeleteArticle)) //删除文章
		//菜单子项
		auth.POST("menuchild/add", Handler()(v1.AddMenuchild)) //添加菜单子项
		auth.PUT("menuchild/set", Handler()(v1.SetMenuchild))  //设置菜单子项
		//分类
		auth.POST("category/add", Handler()(v1.AddCategory))     //添加分类
		auth.PUT("category/list", Handler()(v1.ModifyCategorys)) //批量修改分类
		//上传
		auth.POST("upload/image", v1.Upload) //上}传文件
		//消息
		auth.PUT("message/updata", Handler()(v1.UpMessage))
		auth.DELETE("message/del", Handler()(v1.DelMessage))
		auth.PUT("message/reply", Handler()(v1.ReplyQuestio))
	}

	router := r.Group("api/v1")
	{
		//文章
		router.GET("article/list", Handler()(v1.GetsArticle))      //获取文章列表
		router.GET("article/:id", Handler()(v1.GetArticle))        //获取单个文章
		router.GET("article/tag/:id", Handler()(v1.TagGetArticle)) //根据标签获取所有文章
		//用户
		router.POST("user/register", Handler()(v1.Register)) //注册用户
		router.POST("user/login", Handler()(v1.Login))       //登陆用户
		//菜单子项
		router.GET("menuchilds", Handler()(v1.GetMenuchild))     //获取菜单子项
		router.GET("menuchild", Handler()(v1.GetSingleMenuItem)) //获取单个菜单项
		//分类
		router.GET("category", Handler()(v1.GetCategory)) //获取分类
		//标签
		router.GET("tags", Handler()(v1.GetTags)) //获取全部标签
		//消息
		router.POST("message", Handler()(v1.AddMessage))                              //进行留言
		router.POST("question", Handler()(v1.AddQuestion))                            //进行提问
		router.GET("message", middleware.JwtToken(false), Handler()(v1.GetMessage))   //获取留言
		router.GET("question", middleware.JwtToken(false), Handler()(v1.GetQuestion)) //获取提问
	}
	err := r.Run(utils.Config.Server.HttpPort)
	if err != nil {
		return
	}
}
