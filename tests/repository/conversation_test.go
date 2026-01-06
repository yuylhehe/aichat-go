package repository_test

import (
	"ai-chat/internal/model"
	"ai-chat/internal/repository"
	"ai-chat/tests/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversationRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewConversationRepository(db)

	// 先创建用户
	user := testutil.CreateTestUser(db, "测试用户", "conv_create@example.com", "password", "salt")

	t.Run("成功创建会话", func(t *testing.T) {
		conv := &model.Conversation{
			Name:   "测试会话",
			UserID: user.ID,
		}

		err := repo.Create(conv)

		assert.NoError(t, err)
		assert.NotZero(t, conv.ID)
	})
}

func TestConversationRepository_FindByID(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewConversationRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "conv_findbyid@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "测试会话", user.ID)

	t.Run("成功查找会话", func(t *testing.T) {
		result, err := repo.FindByID(conv.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "测试会话", result.Name)
	})

	t.Run("会话不存在", func(t *testing.T) {
		result, err := repo.FindByID(99999)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "会话不存在")
	})
}

func TestConversationRepository_FindByIDAndUserID(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewConversationRepository(db)

	user1 := testutil.CreateTestUser(db, "用户1", "conv_user1@example.com", "password", "salt")
	user2 := testutil.CreateTestUser(db, "用户2", "conv_user2@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "用户1的会话", user1.ID)

	t.Run("成功查找归属用户的会话", func(t *testing.T) {
		result, err := repo.FindByIDAndUserID(conv.ID, user1.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "用户1的会话", result.Name)
	})

	t.Run("无权访问其他用户的会话", func(t *testing.T) {
		result, err := repo.FindByIDAndUserID(conv.ID, user2.ID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestConversationRepository_FindByUserID(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewConversationRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "conv_list@example.com", "password", "salt")
	testutil.CreateTestConversation(db, "会话1", user.ID)
	testutil.CreateTestConversation(db, "会话2", user.ID)
	testutil.CreateTestConversation(db, "搜索测试会话", user.ID)

	t.Run("获取用户所有会话", func(t *testing.T) {
		result, err := repo.FindByUserID(user.ID, "")

		assert.NoError(t, err)
		assert.Len(t, result, 3)
	})

	t.Run("搜索会话", func(t *testing.T) {
		// 注意: SQLite 不支持 ILIKE，这个测试在集成环境中会通过
		// 这里只测试查询为空字符串的情况
		t.Skip("SQLite 不支持 ILIKE 语法，跳过搜索测试")
	})
}

func TestConversationRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewConversationRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "conv_update@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "原始名称", user.ID)

	t.Run("成功更新会话", func(t *testing.T) {
		updates := map[string]interface{}{
			"name": "更新后的名称",
		}
		err := repo.Update(conv, updates)

		assert.NoError(t, err)

		// 验证更新
		updated, _ := repo.FindByID(conv.ID)
		assert.Equal(t, "更新后的名称", updated.Name)
	})
}

func TestConversationRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewConversationRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "conv_delete@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "待删除会话", user.ID)

	t.Run("成功删除会话", func(t *testing.T) {
		err := repo.Delete(conv.ID)

		assert.NoError(t, err)

		// 验证已删除
		_, err = repo.FindByID(conv.ID)
		assert.Error(t, err)
	})
}

func TestConversationRepository_CountByIDAndUserID(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewConversationRepository(db)

	user1 := testutil.CreateTestUser(db, "用户1", "conv_count1@example.com", "password", "salt")
	user2 := testutil.CreateTestUser(db, "用户2", "conv_count2@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "测试会话", user1.ID)

	t.Run("会话归属用户", func(t *testing.T) {
		count, err := repo.CountByIDAndUserID(conv.ID, user1.ID)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("会话不归属用户", func(t *testing.T) {
		count, err := repo.CountByIDAndUserID(conv.ID, user2.ID)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}
