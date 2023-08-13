package service

import (
	"common/config"
	model "common/model/practiceModel"
	"common/proto/practiceProto"
	"common/results"
	"common/utils"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"log"
	"practice/grpc/practiceGrpc"
	"practice/rabbitmq"
	"reflect"
	"sort"
)

func QuestionList(c *gin.Context) *results.JsonResult {
	//构建context
	reqCtx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("token", c.GetHeader("token")))
	//发起grpc请求
	res, err := practiceGrpc.PracticeServiceClient.QuestionList(reqCtx, &practiceProto.QuestionListRequest{PaperID: c.GetHeader("paperId")})
	if err != nil {
		return results.Fail("rpc请求失败", err)
	}
	return results.Success("请求成功", res.QuestionList)
}

func TableStruct(c *gin.Context) *results.JsonResult {
	//构建context
	reqCtx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("token", c.GetHeader("token")))
	//发起grpc请求
	res, err := practiceGrpc.PracticeServiceClient.TableStruct(reqCtx, &practiceProto.TableStructRequest{PaperID: c.GetHeader("paperId")})
	if err != nil {
		return results.Fail("rpc请求失败", err)
	}
	resultMap := make(map[string][]*practiceProto.ColumnComment)
	for k, v := range res.TableStructMap {
		resultMap[k] = v.ColumnCommentList
	}

	return results.Success("请求成功", resultMap)

	//控制携程
	//treeChannel := make(chan struct{}, 0)
	//studentChannel := make(chan struct{}, 0)

	//子携程
	//studentChannel <- struct{}{}
	//defer close(studentChannel)

	//<-treeChannel
	//defer close(treeChannel)

	//主携程阻塞
	//treeChannel <- struct{}{}
	//<-studentChannel

	//log.Println("查询表结构耗时", time.Now().UnixMilli()-milli)
}

func LevelList(c *gin.Context) *results.JsonResult {

	//构建context
	reqCtx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("token", c.GetHeader("token")))
	//发起grpc请求
	res, err := practiceGrpc.PracticeServiceClient.LevelList(reqCtx, &practiceProto.LevelListRequest{})
	if err != nil {
		return results.Fail("rpc请求失败", err)
	}
	return results.Success("查找成功", res.LevelList)

}

func PaperList(c *gin.Context) *results.JsonResult {

	//构建context
	reqCtx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("token", c.GetHeader("token")))
	//发起grpc请求
	res, err := practiceGrpc.PracticeServiceClient.PaperList(reqCtx, &practiceProto.PaperListRequest{})
	if err != nil {
		return results.Fail("rpc请求失败", err)
	}

	resultMap := make(map[int32][]*practiceProto.Paper)
	for k, v := range res.PaperMap {
		resultMap[k] = v.PaperList
	}

	return results.Success("查找成功", resultMap)
}

func CommitSQL(c *gin.Context) *results.JsonResult {
	data, _ := c.GetRawData()

	var body struct {
		PracticeId string `json:"practiceId"`
		SQLString  string `json:"sqlString"`
	}
	if err := sonic.Unmarshal(data, &body); err != nil {
		return results.Error("参数错误", err)
	}
	var resultMapSlice []map[string]any
	if result := config.DB.Raw(body.SQLString).Scan(&resultMapSlice); result.Error != nil {
		return results.Success("查找成功", resultMapSlice)
	}

	//答题sql查到数据
	answerMapSliceLen := len(resultMapSlice)
	if answerMapSliceLen > 0 {
		var practice model.Practice
		if result2 := config.DB.Model(&practice).Select("answer", "paper_id").Where("id = ?", body.PracticeId).Find(&practice); result2.Error != nil {
			return results.Success("查找成功", resultMapSlice)
		}

		//预定义的答案
		var answerMapSlice []map[string]any
		if result3 := config.DB.Raw(practice.Answer).Scan(&answerMapSlice); result3.Error != nil {
			return results.Success("查找成功", resultMapSlice)
		}

		if len(answerMapSlice) != answerMapSliceLen {
			//todo
			return results.Success("查找成功", resultMapSlice)
		}

		//答案map的values == 结果map的values ?
		var resultStringSlice []string
		var answerStringSlice []string

		for _, v := range resultMapSlice {
			for _, v2 := range v {
				resultStringSlice = append(resultStringSlice, fmt.Sprintf("%v", v2))
			}
		}

		for _, v := range answerMapSlice {
			for _, v2 := range v {
				answerStringSlice = append(answerStringSlice, fmt.Sprintf("%v", v2))
			}
		}

		sort.Strings(resultStringSlice)
		sort.Strings(answerStringSlice)
		areEqual := reflect.DeepEqual(resultStringSlice, answerStringSlice)

		//答案正确
		if areEqual {
			token := c.Request.Header.Get("Authorization")
			// 未传token
			if len(token) <= 0 {
				return results.Success("查找成功", resultMapSlice)
			}

			userId, err := utils.ParseToken(token)
			if err != nil {
				// token解析失败
				log.Println("token解析失败", err)
				return results.Success("查找成功", resultMapSlice)
			}

			if userId == 0 {
				return results.Success("查找成功", resultMapSlice)
			}

			//用户已登录
			//异步 记录答案 + sql执行耗时排行
			elapsed, err := elapsedSQL(body.SQLString)
			if err != nil {
				log.Println("sql执行耗时统计失败", err)
				return results.Success("查找成功", resultMapSlice)
			}
			userCommitAnswerLog := model.UserCommitAnswerLog{
				QuestionId: body.PracticeId,
				UserId:     fmt.Sprintf("%v", userId),
				SQLRunTime: elapsed,
				AnswerSql:  body.SQLString}

			marshal, err := sonic.Marshal(userCommitAnswerLog)
			if err != nil {
				log.Println("序列化失败", err)
				return results.Success("查找成功", resultMapSlice)
			}
			utils.PushMessage(marshal, rabbitmq.PracticeAnswerChannel, "practiceExchange", "practiceAnswerKey")
		}

	}

	return results.Success("查找成功", resultMapSlice)
}

// elapsedSQL 获取SQL执行耗时
func elapsedSQL(commitSql string) (float64, error) {

	//打开性能统计
	config.DB.Exec("set profiling=1;")

	//需要统计的sql
	var resultMapSlice []map[string]any
	config.DB.Raw(commitSql).Scan(&resultMapSlice)

	//展示统计信息
	var profiledModel model.ProfileModel
	result := config.DB.Raw("show profiles;").Scan(&profiledModel)

	if result.Error != nil {
		log.Println("查询耗时失败", result.Error)
		return 0, result.Error
	}
	return profiledModel.Duration, nil
}
