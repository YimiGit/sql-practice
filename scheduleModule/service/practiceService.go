package service

import (
	"common/config"
	"common/model/adminModel"
	"common/utils"
	"context"
	"github.com/bytedance/sonic"
	"github.com/xxl-job/xxl-job-executor-go"
	"time"
)

func SqlPassTotal(c context.Context, param *xxl.RunReq) error {
	type sqlPassTotals struct {
		QuestionId string `json:"question_id"`
		Total      int    `json:"total"`
	}
	var sqlPassTotal []sqlPassTotals
	//查询sql作答正确的数据
	result := config.DB.Raw("select question_id,count(1) as total from user_question_answer where is_correct = 1 group by question_id").Scan(&sqlPassTotal)
	if result.Error != nil {
		return result.Error
	}

	marshal, err := sonic.Marshal(sqlPassTotal)
	if err != nil {
		return err
	}

	//放入redis
	redisResult := config.RedisClient.Set(c, "sqlPassTotal", marshal, time.Duration(utils.TodayRemainNanosecond()))
	if redisResult.Err() != nil {
		return redisResult.Err()
	}
	return nil
}

func CleanJogLog(c context.Context, param *xxl.RunReq) error {
	result := config.DB.Exec("delete from xxl_job_log where handle_code = 200 and trigger_time < DATE_SUB(NOW(),INTERVAL 3 DAY)")
	return result.Error
}

func UserLoginMonth(c context.Context, param *xxl.RunReq) error {
	var loginLogs []adminModel.UserLoginTotal
	//本月1号到现在的登录日志
	result := config.DB.Where("login_time > DATE_FORMAT(NOW(),'%Y-%m-01')").Find(&loginLogs)

	if result.Error != nil {
		return result.Error
	}

	marshal, err := sonic.Marshal(loginLogs)
	if err != nil {
		return err
	}

	//放入redis
	redisResult := config.RedisClient.Set(c, "userLoginMonth", marshal, time.Duration(utils.TodayRemainNanosecond()))
	if redisResult.Err() != nil {
		return redisResult.Err()
	}
	return nil
}
