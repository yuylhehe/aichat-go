package repository_test

import (
	"ai-chat/internal/common"
	"ai-chat/internal/model"
	"ai-chat/internal/repository"
	"ai-chat/tests/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewMessageRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "msg_create@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "测试会话", user.ID)

	t.Run("成功创建消息", func(t *testing.T) {
		msg := &model.Message{
			ConversationID: conv.ID,
			Content:        "测试消息内容",
			Type:           common.RoleUser,
			Sort:           1,
		}

		err := repo.Create(msg)

		assert.NoError(t, err)
		assert.NotZero(t, msg.ID)
	})
}

func TestMessageRepository_FindByID(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewMessageRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "msg_findbyid@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "测试会话", user.ID)
	msg := testutil.CreateTestMessage(db, conv.ID, "测试消息", common.RoleUser, 1)

	t.Run("成功查找消息", func(t *testing.T) {
		result, err := repo.FindByID(msg.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "测试消息", result.Content)
	})

	t.Run("消息不存在", func(t *testing.T) {
		result, err := repo.FindByID(99999)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestMessageRepository_FindByConversationID(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewMessageRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "msg_list@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "测试会话", user.ID)

	// 创建多条消息
	testutil.CreateTestMessage(db, conv.ID, "消息1", common.RoleUser, 1)
	testutil.CreateTestMessage(db, conv.ID, "消息2", common.RoleAssistant, 2)
	testutil.CreateTestMessage(db, conv.ID, "消息3", common.RoleUser, 3)

	t.Run("获取会话所有消息", func(t *testing.T) {
		result, err := repo.FindByConversationID(conv.ID)

		assert.NoError(t, err)
		assert.Len(t, result, 3)
	})

	t.Run("按排序值排序", func(t *testing.T) {
		result, err := repo.FindByConversationID(conv.ID)

		assert.NoError(t, err)
		assert.Equal(t, 1, result[0].Sort)
		assert.Equal(t, 2, result[1].Sort)
		assert.Equal(t, 3, result[2].Sort)
	})
}

func TestMessageRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewMessageRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "msg_update@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "测试会话", user.ID)
	msg := testutil.CreateTestMessage(db, conv.ID, "原始内容", common.RoleUser, 1)

	t.Run("成功更新消息", func(t *testing.T) {
		updates := map[string]interface{}{
			"content": "更新后的内容",
		}
		err := repo.Update(msg, updates)

		assert.NoError(t, err)

		// 验证更新
		updated, _ := repo.FindByID(msg.ID)
		assert.Equal(t, "更新后的内容", updated.Content)
	})
}

func TestMessageRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewMessageRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "msg_delete@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "测试会话", user.ID)
	msg := testutil.CreateTestMessage(db, conv.ID, "待删除消息", common.RoleUser, 1)

	t.Run("成功删除消息", func(t *testing.T) {
		err := repo.Delete(msg.ID)

		assert.NoError(t, err)

		// 验证已删除
		_, err = repo.FindByID(msg.ID)
		assert.Error(t, err)
	})
}

func TestMessageRepository_DeleteByConversationID(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewMessageRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "msg_delconv@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "测试会话", user.ID)

	testutil.CreateTestMessage(db, conv.ID, "消息1", common.RoleUser, 1)
	testutil.CreateTestMessage(db, conv.ID, "消息2", common.RoleAssistant, 2)

	t.Run("成功删除会话所有消息", func(t *testing.T) {
		err := repo.DeleteByConversationID(conv.ID)

		assert.NoError(t, err)

		// 验证消息已全部删除
		messages, _ := repo.FindByConversationID(conv.ID)
		assert.Len(t, messages, 0)
	})
}

func TestMessageRepository_NextSort(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewMessageRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "msg_nextsort@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "测试会话", user.ID)

	t.Run("空会话返回1", func(t *testing.T) {
		result, err := repo.NextSort(conv.ID)

		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("有消息返回最大值+1", func(t *testing.T) {
		testutil.CreateTestMessage(db, conv.ID, "消息1", common.RoleUser, 1)
		testutil.CreateTestMessage(db, conv.ID, "消息2", common.RoleAssistant, 5)

		result, err := repo.NextSort(conv.ID)

		assert.NoError(t, err)
		assert.Equal(t, 6, result)
	})
}

func TestMessageRepository_CountByConversationID(t *testing.T) {
	db := testutil.SetupTestDB()
	repo := repository.NewMessageRepository(db)

	user := testutil.CreateTestUser(db, "测试用户", "msg_count@example.com", "password", "salt")
	conv := testutil.CreateTestConversation(db, "测试会话", user.ID)

	testutil.CreateTestMessage(db, conv.ID, "消息1", common.RoleUser, 1)
	testutil.CreateTestMessage(db, conv.ID, "消息2", common.RoleAssistant, 2)
	testutil.CreateTestMessage(db, conv.ID, "消息3", common.RoleUser, 3)

	t.Run("正确统计消息数量", func(t *testing.T) {
		count, err := repo.CountByConversationID(conv.ID)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("空会话返回0", func(t *testing.T) {
		emptyConv := testutil.CreateTestConversation(db, "空会话", user.ID)

		count, err := repo.CountByConversationID(emptyConv.ID)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}
