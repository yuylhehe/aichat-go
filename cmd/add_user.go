package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"ai-chat/config"
	"ai-chat/internal/repository"
	"ai-chat/internal/service"
)

func main() {
	// 1. 加载配置 (确保能读取到 .env)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("加载配置失败: ", err)
	}

	// 2. 连接数据库
	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatal("连接数据库失败: ", err)
	}

	reader := bufio.NewReader(os.Stdin)

	// 3. 获取用户输入
	fmt.Println("=== 手动添加用户账户 ===")

	fmt.Print("请输入用户名 (Name): ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("请输入邮箱 (Email): ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("请输入密码 (Password): ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if name == "" || email == "" || password == "" {
		log.Fatal("错误: 所有字段都不能为空")
	}

	// 4. 检查邮箱是否已存在
	var existingUser repository.User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		log.Fatalf("错误: 邮箱 '%s' 已被注册", email)
	}

	// 5. 创建用户 (使用由于 service 层提供的加密逻辑)
	salt := service.GenerateSalt()
	hashedPassword := service.HashPassword(password, salt)

	user := &repository.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Salt:     salt,
		IsActive: true,
	}

	if err := db.Create(user).Error; err != nil {
		log.Fatal("创建用户失败: ", err)
	}

	fmt.Println("\n✅ 用户创建成功！")
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("用户名: %s\n", user.Name)
	fmt.Printf("邮箱: %s\n", user.Email)
}
