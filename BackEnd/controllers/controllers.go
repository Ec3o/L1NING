package controllers

import (
	"BackEnd/config"
	"BackEnd/io"
	"BackEnd/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// AddTodo 增加Todo
func AddTodo(c *gin.Context) {
	var todos []models.Todo
	if err := c.ShouldBindJSON(&todos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claimsValue, _ := c.Get("claims")
	if claimsValue == nil {
		c.JSON(http.StatusUnauthorized, "未找到 Claims")
		c.Abort()
		return
	}

	claims, ok := claimsValue.(*models.Claims)
	if !ok {
		actualType := reflect.TypeOf(claimsValue)
		c.JSON(http.StatusUnauthorized, fmt.Sprintf("非法 Claims 类型: %s", actualType.String()))
		c.Abort()
		return
	}

	username := claims.Username
	fmt.Println("username:" + username)

	// 读取现有的TODOs
	existingTodos, err := io.LoadTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法加载现有的TODOs"})
		return
	}

	// 找到当前用户的最大ID
	maxID := uint(0)
	for _, todo := range existingTodos {
		if todo.User == username && todo.ID > maxID {
			maxID = todo.ID
		}
	}

	// 设置每个新TODO的用户名并分配唯一ID
	for i := range todos {
		todos[i].User = username
		maxID++
		todos[i].ID = maxID
	}

	// 将新TODOs添加到现有的TODOs列表中
	allTodos := append(existingTodos, todos...)

	// 保存所有TODOs
	err = io.SaveTodos(allTodos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存TODOs"})
		return
	}

	c.JSON(http.StatusOK, models.TodoSubmitSuccess)
}

// DeleteTodo 删除TODO
func DeleteTodo(c *gin.Context) {
	todoID := c.Param("id")
	println(todoID)
	id, err := strconv.ParseUint(todoID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的TODO ID"})
		return
	}

	claimsValue, _ := c.Get("claims")
	if claimsValue == nil {
		c.JSON(http.StatusUnauthorized, "未找到 Claims")
		c.Abort()
		return
	}

	claims, ok := claimsValue.(*models.Claims)
	if !ok {
		actualType := reflect.TypeOf(claimsValue)
		c.JSON(http.StatusUnauthorized, fmt.Sprintf("非法 Claims 类型: %s", actualType.String()))
		c.Abort()
		return
	}

	username := claims.Username

	// 读取现有的TODOs
	existingTodos, err := io.LoadTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法加载现有的TODOs"})
		return
	}

	// 查找并删除TODO
	found := false
	var newTodos []models.Todo
	for _, todo := range existingTodos {
		if todo.ID == uint(id) && todo.User == username {
			found = true
		} else {
			newTodos = append(newTodos, todo)
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "TODO未找到或无权限"})
		return
	}

	// 保存更新后的TODOs
	err = io.SaveTodos(newTodos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存TODOs"})
		return
	}

	c.JSON(http.StatusOK, models.TodoSubmitSuccess)
}

// GetTodo 查询TODO
func GetTodo(c *gin.Context) {
	claimsValue, _ := c.Get("claims")
	if claimsValue == nil {
		c.JSON(http.StatusUnauthorized, "未找到 Claims")
		c.Abort()
		return
	}

	claims, ok := claimsValue.(*models.Claims)
	if !ok {
		actualType := reflect.TypeOf(claimsValue)
		c.JSON(http.StatusUnauthorized, fmt.Sprintf("非法 Claims 类型: %s", actualType.String()))
		c.Abort()
		return
	}

	username := claims.Username
	var todos []models.Todo
	existingTodos, err := io.LoadTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, todo := range existingTodos {
		if todo.User == username {
			todos = append(todos, todo)
		}
	}
	if todos == nil {
		todos = []models.Todo{}
	}
	c.JSON(http.StatusOK, todos)
}

// UpdateTodo 修改TODO
func UpdateTodo(c *gin.Context) {
	var updatedTodo models.Todo
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claimsValue, _ := c.Get("claims")
	if claimsValue == nil {
		c.JSON(http.StatusUnauthorized, "未找到 Claims")
		c.Abort()
		return
	}

	claims, ok := claimsValue.(*models.Claims)
	if !ok {
		actualType := reflect.TypeOf(claimsValue)
		c.JSON(http.StatusUnauthorized, fmt.Sprintf("非法 Claims 类型: %s", actualType.String()))
		c.Abort()
		return
	}

	username := claims.Username

	// 读取现有的TODOs
	existingTodos, err := io.LoadTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法加载现有的TODOs"})
		return
	}

	// 查找并更新TODO
	found := false
	for i := range existingTodos {
		if existingTodos[i].ID == updatedTodo.ID && existingTodos[i].User == username {
			existingTodos[i] = updatedTodo
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "TODO未找到或无权限"})
		return
	}

	// 保存更新后的TODOs
	err = io.SaveTodos(existingTodos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存TODOs"})
		return
	}

	c.JSON(http.StatusOK, models.TodoSubmitSuccess)
}

// Register 注册 @Developed
func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, models.ErrInvalidUSERFormat)
		return
	}

	if len(user.Password) <= 6 { //密码长度过短提示重新设置
		c.JSON(400, models.ErrInvalidPassword)
		return
	}

	existingUsers, err := io.LoadUsers()
	if err != nil {
		c.JSON(500, models.ErrReadUserData)
		return
	}

	// 检查是否已经存在相同的用户名
	for _, existingUser := range existingUsers {
		if existingUser.Username == user.Username {
			c.JSON(400, models.ErrRegister)
			return
		}
	}

	existingUsers = append(existingUsers, user)
	err = io.SaveUsers(existingUsers)
	if err != nil {
		c.JSON(500, models.ErrSaveUserData)
		return
	}

	c.JSON(200, models.UserRegisterSuccess)
}

// Login 登录 @Developed
func Login(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, models.ErrInvalidUSERFormat)
		return
	}

	existingUsers, err := io.LoadUsers()
	if err != nil {
		c.JSON(500, models.ErrReadUserData)
		return
	}

	var foundUser models.User
	for _, existingUser := range existingUsers {
		if existingUser.Username == user.Username && existingUser.Password == user.Password {
			foundUser = existingUser
			break
		}
	}
	if foundUser.Username != "" {
		expirationTime := time.Now().Add(500 * time.Minute)
		claims := &models.Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(config.JwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "生成令牌失败")
			return
		}
		// 返回用户名和令牌
		c.JSON(200, gin.H{"status": "用户登录成功", "username": foundUser.Username, "token": tokenString})
	} else {
		c.JSON(404, models.ErrUserlogin)
	}
}
