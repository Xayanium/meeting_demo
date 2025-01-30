package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"meeting_demo/internal/models"
	"meeting_demo/utils"
	"net/http"
	"time"
)

type MeetingListRequest struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Keyword string `json:"keyword"`
}

type MeetingListResponse struct {
	Identity string    `json:"identity"`
	Name     string    `json:"name"`
	BeginAt  time.Time `json:"begin_at"`
	EndAt    time.Time `json:"end_at"`
}

type MeetingCreateRequest struct {
	Name    string `json:"name"`
	BeginAt int64  `json:"begin_at"`
	EndAt   int64  `json:"end_at"`
}

type MeetingEditRequest struct {
	*MeetingCreateRequest
	Identity string `json:"identity"`
}

func MeetingList(c *gin.Context) {
	userAuth := c.MustGet("user_token").(*utils.UserAuth)

	meetingInfo := new(MeetingListRequest)
	err := c.ShouldBindQuery(&meetingInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Invalid Params",
		})
		log.Println("bind list meeting params error: " + err.Error())
		return
	}

	// 对列出的会议进行分页操作
	var count int64
	var rooms []*MeetingListResponse // 返回前端的最终结果
	var roomInfos []*models.RoomInfo // 查询结果暂存
	tx := models.DB.Model(&models.RoomInfo{})
	if meetingInfo.Keyword != "" {
		err = tx.Where("user_id = ? AND name LIKE ?", userAuth.Id, "%"+meetingInfo.Keyword+"%").
			Count(&count).Offset((meetingInfo.Page - 1) * meetingInfo.Size).Limit(meetingInfo.Size).Find(&rooms).Error
	} else {
		err = tx.Where("user_id = ?", userAuth.Id).
			Count(&count).Offset((meetingInfo.Page - 1) * meetingInfo.Size).Limit(meetingInfo.Size).Find(&rooms).Error
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "System Server Error",
		})
		log.Println("list room info error: " + err.Error())
		return
	}

	// 将临时结果部分字段映射到最终结果
	for _, roomInfo := range roomInfos {
		rooms = append(rooms, &MeetingListResponse{
			Identity: roomInfo.Identity,
			Name:     roomInfo.Name,
			BeginAt:  roomInfo.BeginAt,
			EndAt:    roomInfo.EndAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"rooms": rooms,
			"total": count,
		},
	})
}

func MeetingCreate(c *gin.Context) {
	userAuth := c.MustGet("user_token").(*utils.UserAuth)

	meetingInfo := new(MeetingCreateRequest)
	err := c.ShouldBindJSON(meetingInfo) // 将请求体中的JSON数据解析，绑定到Go的结构体中
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Invalid Params",
		})
		log.Println("bind create meeting params error: " + err.Error())
		return
	}

	err = models.DB.Create(&models.RoomInfo{
		Identity: utils.GetUUID(),
		Name:     meetingInfo.Name,
		BeginAt:  time.Unix(meetingInfo.BeginAt, 0), // 将前端传的时间戳(s)转为time.Time
		EndAt:    time.Unix(meetingInfo.EndAt, 0),
		UserId:   userAuth.Id, // 从中间件中获取用户信息
	}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "System Server Error",
		})
		log.Println("insert room info error: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

func MeetingEdit(c *gin.Context) {
	userAuth := c.MustGet("user_token").(*utils.UserAuth)

	meetingInfo := new(MeetingEditRequest)
	err := c.ShouldBindJSON(meetingInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Invalid Params",
		})
		log.Println("bind edit meeting params error: " + err.Error())
		return
	}

	// 只更新非空字段
	updates := make(map[string]interface{})
	if meetingInfo.Name != "" {
		updates["name"] = meetingInfo.Name
	}
	if meetingInfo.BeginAt != 0 {
		updates["begin_at"] = time.Unix(meetingInfo.BeginAt, 0)
	}
	if meetingInfo.EndAt != 0 {
		updates["end_at"] = time.Unix(meetingInfo.EndAt, 0)
	}

	// 修改时进行鉴权，验证会议标识及创建会议的用户
	if len(updates) != 0 {
		err = models.DB.Model(&models.RoomInfo{}).
			Where("identity = ? And user_id = ?", meetingInfo.Identity, userAuth.Id).
			Updates(updates).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "System Server Error",
			})
			log.Println("update room info error: " + err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

func MeetingDelete(c *gin.Context) {
	identity := c.Query("identity")
	userAuth := c.MustGet("user_token").(*utils.UserAuth)

	err := models.DB.Model(&models.RoomInfo{}).
		Where("identity = ? AND user_id = ?", identity, userAuth.Id).Delete(nil).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "System Server Error",
		})
		log.Println("delete room info error: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
