# 个人博客系统后端

基于 Go + Gin + GORM 开发的个人博客系统后端，支持用户认证、文章管理和评论功能。

## 功能特性

- ✅ 用户注册和登录（JWT 认证）
- ✅ 博客文章的 CRUD 操作
- ✅ 文章评论功能
- ✅ 权限控制（用户只能操作自己的文章）
- ✅ 错误处理和日志记录

## 技术栈

- **框架**: Gin
- **ORM**: GORM

### 运行步骤
``` bash
cd blog_system

go mod tidy

go run main.go