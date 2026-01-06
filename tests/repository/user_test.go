package repository_test

import (
	"ai-chat/internal/model"
	"ai-chat/internal/repository"
	"ai-chat/tests/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewUserRepository(db)

	t.Run("成功创建用户", func(t *testing.T) {
		user := &model.User{
			Name:     "测试用户",
			Email:    "create@example.com",
			Password: "hashedPassword",
			Salt:     "salt",
		}

		err := repo.Create(user)

		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
	})

	t.Run("邮箱重复创建失败", func(t *testing.T) {
		user1 := &model.User{
			Name:     "用户1",
			Email:    "duplicate@example.com",
			Password: "hashedPassword",
			Salt:     "salt",
		}
		repo.Create(user1)

		user2 := &model.User{
			Name:     "用户2",
			Email:    "duplicate@example.com",
			Password: "hashedPassword",
			Salt:     "salt",
		}
		err := repo.Create(user2)

		assert.Error(t, err)
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewUserRepository(db)

	t.Run("成功查找用户", func(t *testing.T) {
		user := &model.User{
			Name:     "查找测试用户",
			Email:    "findbyemail@example.com",
			Password: "hashedPassword",
			Salt:     "salt",
		}
		repo.Create(user)

		result, err := repo.FindByEmail("findbyemail@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "findbyemail@example.com", result.Email)
		assert.Equal(t, "查找测试用户", result.Name)
	})

	t.Run("用户不存在", func(t *testing.T) {
		result, err := repo.FindByEmail("notfound@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewUserRepository(db)

	t.Run("成功查找用户", func(t *testing.T) {
		user := &model.User{
			Name:     "ID查找测试用户",
			Email:    "findbyid@example.com",
			Password: "hashedPassword",
			Salt:     "salt",
		}
		repo.Create(user)

		result, err := repo.FindByID(user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, "ID查找测试用户", result.Name)
	})

	t.Run("用户不存在", func(t *testing.T) {
		result, err := repo.FindByID(99999)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUserRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewUserRepository(db)

	t.Run("成功更新用户", func(t *testing.T) {
		user := &model.User{
			Name:     "原始名称",
			Email:    "update@example.com",
			Password: "hashedPassword",
			Salt:     "salt",
		}
		repo.Create(user)

		user.Name = "更新后的名称"
		err := repo.Update(user)

		assert.NoError(t, err)

		// 验证更新
		updated, _ := repo.FindByID(user.ID)
		assert.Equal(t, "更新后的名称", updated.Name)
	})
}

func TestUserRepository_DeleteByID(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewUserRepository(db)

	t.Run("成功删除用户", func(t *testing.T) {
		user := &model.User{
			Name:     "待删除用户",
			Email:    "delete@example.com",
			Password: "hashedPassword",
			Salt:     "salt",
		}
		repo.Create(user)

		err := repo.DeleteByID(user.ID)

		assert.NoError(t, err)

		// 验证已删除
		_, err = repo.FindByID(user.ID)
		assert.Error(t, err)
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewUserRepository(db)

	// 创建测试数据
	for i := 0; i < 5; i++ {
		user := &model.User{
			Name:     "用户" + string(rune('A'+i)),
			Email:    "user" + string(rune('a'+i)) + "@example.com",
			Password: "hashedPassword",
			Salt:     "salt",
		}
		repo.Create(user)
	}

	t.Run("分页查询", func(t *testing.T) {
		users, total, err := repo.FindAll(1, 3, "")

		assert.NoError(t, err)
		assert.Len(t, users, 3)
		assert.Equal(t, int64(5), total)
	})

	t.Run("搜索查询", func(t *testing.T) {
		users, total, err := repo.FindAll(1, 10, "用户A")

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, int64(1), total)
		for _, u := range users {
			assert.Contains(t, u.Name, "用户A")
		}
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewUserRepository(db)

	t.Run("成功更新密码", func(t *testing.T) {
		user := &model.User{
			Name:     "密码更新测试",
			Email:    "password@example.com",
			Password: "oldPassword",
			Salt:     "oldSalt",
		}
		repo.Create(user)

		err := repo.UpdatePassword(user.ID, "newHashedPassword", "newSalt")

		assert.NoError(t, err)

		// 验证更新
		updated, _ := repo.FindByID(user.ID)
		assert.Equal(t, "newHashedPassword", updated.Password)
		assert.Equal(t, "newSalt", updated.Salt)
	})
}
