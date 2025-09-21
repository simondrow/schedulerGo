package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name  string            `json:"name" bson:"name"`
	Tasks map[string][]string `json:"tasks" bson:"tasks"`
}

type UserResponse struct {
	Name  string            `json:"name"`
	Tasks map[string][]string `json:"任务"`
}

type UpdateTasksRequest struct {
	Name  string            `json:"name"`
	Tasks map[string][]string `json:"tasks,omitempty"`
}

type UpdateDayTasksRequest struct {
	Tasks []string `json:"tasks"`
}

var client *mongo.Client
var usersCollection *mongo.Collection

// 安全的 JSON 序列化函数
func safeStringify(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", data)
	}
	return string(jsonBytes)
}

// 日志中间件
func loggingMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		
		// 请求开始日志
		fmt.Printf("[REQ] %s %s from %s\n", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		
		// 保存原始的 JSON 方法
		originalJson := c.Writer.Header().Get("Content-Type")
		
		// 包装 JSON 响应
		writer := c.Writer
		c.Writer = &responseWriter{ResponseWriter: writer, originalJson: originalJson, context: c}
		
		c.Next()
		
		// 请求完成日志
		fmt.Printf("[RES] %s %s -> %d (%dms)\n", 
			c.Request.Method, 
			c.Request.URL.Path, 
			c.Writer.Status(), 
			time.Since(start).Milliseconds())
	})
}

// 自定义响应写入器
type responseWriter struct {
	gin.ResponseWriter
	originalJson string
	context      *gin.Context
}

func (w *responseWriter) Write(data []byte) (int, error) {
	// 如果是 GET 请求且响应是 JSON，记录响应内容
	if w.context.Request.Method == "GET" && w.Header().Get("Content-Type") == "application/json" {
		fmt.Printf("[GET] %s response:\n%s\n", w.context.Request.URL.Path, safeStringify(string(data)))
	}
	return w.ResponseWriter.Write(data)
}

// 初始化示例数据
func initializeData() error {
	// 检查是否已有数据
	count, err := usersCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	
	if count > 0 {
		return nil
	}

	sampleUsers := []User{
		{
			Name: "Waner",
			Tasks: map[string][]string{
				"1": {
					"7点半滴眼药水",
					"7点半练反转拍",
					"7点半朗读英语背单词",
					"9-11作文课",
					"11点户外",
					"14-15练琴",
					"15点写日记",
					"16-17写作业",
					"17-18户外",
					"19点读书60页",
					"21点半抹油上床",
				},
				"2": {
					"7点滴眼药水",
					"7点朗读英语背单词",
					"16-17钢琴课",
					"17点练反转拍",
					"17点写作业",
					"17点半户外",
					"18点半家教课",
					"21点半抹油上床",
				},
				"3": {
					"7点滴眼药水",
					"7点朗读语文",
					"16点户外",
					"17-19剑桥英语",
					"19点练反转拍",
					"19点练琴",
					"20-21写作业",
					"21点半抹油上床",
				},
				"4": {
					"7点滴眼药水",
					"7点朗读英语背单词",
					"16点户外",
					"17点练反转拍",
					"17-18写校内作业",
					"19点-20点写家庭作业",
					"20点练琴",
					"21点半抹油上床",
				},
				"5": {
					"7点滴眼药水",
					"7点朗读语文",
					"16-17钢琴课",
					"17点练反转拍",
					"17点写校内作业",
					"17点半户外",
					"19点-20点写家庭作业",
					"21点半抹油上床",
				},
				"6": {
					"7点滴眼药水",
					"7点朗读英语背单词",
					"16点户外",
					"17点练反转拍",
					"17点写作业",
					"18点家教课",
					"20点练琴",
					"21点半抹油上床",
				},
				"7": {
					"7点半滴眼药水",
					"7点半练反转拍",
					"7点半朗读语文",
					"9-10读书60页",
					"10-11户外",
					"11-12写一页字",
					"14-15练琴",
					"15点写日记",
					"16-17写作业",
					"17-18户外",
					"21点半抹油上床",
				},
			},
		},
		{
			Name: "John",
			Tasks: map[string][]string{
				"1": {
					"7点半滴眼药水",
					"7点半练反转拍",
					"7点半朗读英语背单词",
					"9-11奥数课",
					"11点户外",
					"14-15练琴",
					"15点写日记",
					"16-17写作业",
					"17-18户外",
					"19点读书60页",
					"21点半抹油上床",
				},
				"2": {
					"7点滴眼药水",
					"7点朗读英语背单词",
					"16点练反转拍",
					"16点写校内作业",
					"17点钢琴课",
					"17点半户外",
					"19点写家庭作业",
					"20点家教课",
					"21点半抹油上床",
				},
				"3": {
					"7点滴眼药水",
					"7点朗读语文",
					"16点户外",
					"17点练反转拍",
					"17-18写校内作业",
					"19点-20点写家庭作业",
					"20点练琴",
					"21点半抹油上床",
				},
				"4": {
					"7点滴眼药水",
					"7点朗读英语背单词",
					"16点户外",
					"17-19剑桥英语",
					"19点练反转拍",
					"19点练琴",
					"20-21写作业",
					"21点半抹油上床",
				},
				"5": {
					"7点滴眼药水",
					"7点朗读语文",
					"16点练反转拍",
					"16点写校内作业",
					"17点钢琴课",
					"17点半户外",
					"19点-20点写家庭作业",
					"21点半抹油上床",
				},
				"6": {
					"7点滴眼药水",
					"7点朗读英语背单词",
					"16点户外",
					"17点练反转拍",
					"17点写作业",
					"18点半练琴",
					"19点半家教课",
					"21点半抹油上床",
				},
				"7": {
					"7点半滴眼药水",
					"7点半练反转拍",
					"7点半朗读语文",
					"9-10读书60页",
					"10-11户外",
					"11-12写一页字",
					"14-15练琴",
					"15点写日记",
					"16-17写作业",
					"17-18户外",
					"21点半抹油上床",
				},
			},
		},
	}

	// 插入示例数据
	var docs []interface{}
	for _, user := range sampleUsers {
		docs = append(docs, user)
	}
	
	_, err = usersCollection.InsertMany(context.Background(), docs)
	if err != nil {
		return err
	}
	
	fmt.Println("示例数据初始化完成")
	return nil
}

// GET /api/users - 获取所有用户数据
func getAllUsers(c *gin.Context) {
	cursor, err := usersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Printf("获取用户数据失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	defer cursor.Close(context.Background())

	var users []User
	if err = cursor.All(context.Background(), &users); err != nil {
		fmt.Printf("解析用户数据失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	result := make(map[string]interface{})
	for _, user := range users {
		result[user.Name] = map[string]interface{}{
			"任务": user.Tasks,
		}
	}

	c.JSON(http.StatusOK, result)
}

// POST /api/users - 根据Name字段修改任务内容
func updateUserTasks(c *gin.Context) {
	var req UpdateTasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name字段是必需的"})
		return
	}

	updateData := bson.M{}
	if req.Tasks != nil {
		updateData["tasks"] = req.Tasks
	}

	// 打印修改内容
	fmt.Printf("[POST] %s modifications:\n%s\n", c.Request.URL.Path, safeStringify(req))

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var updatedUser User
	err := usersCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"name": req.Name},
		bson.M{"$set": updateData},
		opts,
	).Decode(&updatedUser)

	if err != nil {
		fmt.Printf("更新用户数据失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户数据更新成功",
		"user": gin.H{
			"name":  updatedUser.Name,
			"tasks": updatedUser.Tasks,
		},
	})
}

// GET /api/users/:name - 获取特定用户的任务
func getUserTasks(c *gin.Context) {
	name := c.Param("name")
	
	var user User
	err := usersCollection.FindOne(context.Background(), bson.M{"name": name}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		fmt.Printf("获取用户任务失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":  user.Name,
		"tasks": user.Tasks,
	})
}

// GET /api/users/:name/day/:day - 获取特定用户特定日期的任务
func getUserDayTasks(c *gin.Context) {
	name := c.Param("name")
	dayStr := c.Param("day")
	
	dayNum, err := strconv.Atoi(dayStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误"})
		return
	}
	
	if dayNum < 1 || dayNum > 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期必须在1-7之间（1=周日，7=周六）"})
		return
	}
	
	var user User
	err = usersCollection.FindOne(context.Background(), bson.M{"name": name}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		fmt.Printf("获取用户日期任务失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	dayTasks := user.Tasks[dayStr]
	if dayTasks == nil {
		dayTasks = []string{}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"name":  user.Name,
		"day":   dayNum,
		"tasks": dayTasks,
	})
}

// POST /api/users/:name/day/:day - 更新特定用户特定日期的任务
func updateUserDayTasks(c *gin.Context) {
	name := c.Param("name")
	dayStr := c.Param("day")
	
	dayNum, err := strconv.Atoi(dayStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误"})
		return
	}
	
	if dayNum < 1 || dayNum > 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期必须在1-7之间（1=周日，7=周六）"})
		return
	}
	
	var req UpdateDayTasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}
	
	if req.Tasks == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tasks字段必须是数组"})
		return
	}
	
	// 打印修改内容
	fmt.Printf("[POST] %s modifications:\n%s\n", c.Request.URL.Path, safeStringify(gin.H{
		"name": name,
		"day": dayNum,
		"tasks": req.Tasks,
	}))
	
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var updatedUser User
	err = usersCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"name": name},
		bson.M{"$set": bson.M{fmt.Sprintf("tasks.%s", dayStr): req.Tasks}},
		opts,
	).Decode(&updatedUser)

	if err != nil {
		fmt.Printf("更新用户日期任务失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "任务更新成功",
		"user": gin.H{
			"name": updatedUser.Name,
			"day":  dayNum,
			"tasks": updatedUser.Tasks[dayStr],
		},
	})
}

func main() {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 路由器
	r := gin.Default()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(loggingMiddleware())
	
	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})

	// MongoDB 连接
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/scheduler"
	}

	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB连接失败:", err)
	}

	// 检查连接
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("MongoDB连接检查失败:", err)
	}

	fmt.Println("MongoDB连接成功")
	
	// 获取集合
	usersCollection = client.Database("scheduler").Collection("users")

	// 初始化数据
	if err := initializeData(); err != nil {
		log.Fatal("数据初始化失败:", err)
	}

	// API 路由
	api := r.Group("/api")
	{
		api.GET("/users", getAllUsers)
		api.POST("/users", updateUserTasks)
		api.GET("/users/:name", getUserTasks)
		api.GET("/users/:name/day/:day", getUserDayTasks)
		api.POST("/users/:name/day/:day", updateUserDayTasks)
	}

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 启动服务器
	fmt.Printf("服务器运行在端口 %s\n", port)
	r.Run("0.0.0.0:" + port)
}
